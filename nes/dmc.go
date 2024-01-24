package nes

import "fmt"

type dmcMemoryReader struct {
	sampleAddr     uint16
	currentAddr    uint16
	sampleLength   uint16
	bytesRemaining uint16
	sampleBuffer   []byte
	dma            *DMA
}

func (r *dmcMemoryReader) fetch() byte {
	//debug
	if len(r.sampleBuffer) == 0 {
		panic(fmt.Sprint("len(sampleBuffer) == 0"))
	}

	ret := r.sampleBuffer[0]
	r.sampleBuffer = r.sampleBuffer[:0]

	// この時点でbuffer fillを呼び出して良いはず
	return ret
}

func (r *dmcMemoryReader) restart() {
	r.bytesRemaining = r.sampleLength
	r.currentAddr = r.sampleAddr
}

// // output unit: an 8-bit right shift register, a bits-remaining counter, a 7-bit output level (the same one that can be loaded directly via $4011), and a silence flag.
type dmcOutputUnit struct {
	rightShiftRegister   byte
	bitsRemainingCounter byte
	silenceFlag          bool
	level                byte // 7 bit
}

// https://www.nesdev.org/wiki/APU_DMC
// > These periods are all even numbers because there are 2 CPU cycles in an APU cycle.
// > A rate of 428 means the output level changes every 214 APU cycles.
// NTSC
var dmcPeriodTable = []uint16{
	428, 380, 340, 320, 286, 254, 226, 214, 190, 160, 142, 128, 106, 84, 72, 54,
}

type dmc struct {
	enabled       bool
	irqEnabled    bool
	interruptFlag bool
	loop          bool
	rateIndex     byte
	//outputLevel   byte // 7 bit
	// silenceFlag          bool
	// rightShiftRegister   byte
	// sampleBuffer         []byte
	// sampleAddr           uint16
	// sampleLength         uint16
	// currentAddr          uint16
	// bytesRemaining       uint16
	//bitsRemainingCounter byte
	memoryReader *dmcMemoryReader
	outputUnit   *dmcOutputUnit
	timer        timer
}

func newDMC(dma *DMA) *dmc {
	return &dmc{
		memoryReader: &dmcMemoryReader{
			sampleBuffer: make([]byte, 0, 1),
			dma:          dma,
		},
		outputUnit: &dmcOutputUnit{},
	}
}

func (d *dmc) setEnabled(v bool) {
	if !v {
		// > If the DMC bit is clear, the DMC bytes remaining will be set to 0 and the DMC will silence when it empties.
		d.memoryReader.bytesRemaining = 0
	} else {
		// > If the DMC bit is set, the DMC sample will be restarted only if its bytes remaining is 0.
		// > If there are bits remaining in the 1-byte sample buffer, these will finish playing before the next sample is fetched.
		if d.memoryReader.bytesRemaining == 0 {
			d.memoryReader.restart()
		} else {
			//d.enabled = v
		}
		// TODO: これは違う気がするが、どうやってmemory readerを起動させるんだ？
		// 実は常時動いている？？？
		// if len(d.sampleBuffer) == 0 {
		// 	d.dma.SignalDCMDMA(d.currentAddr)
		// }
	}
}

func (d *dmc) output() byte {
	return d.outputUnit.level
}

func (d *dmc) loadRate(rateIndex byte) {
	d.timer.period = dmcPeriodTable[rateIndex]
}

func (d *dmc) setSampleBuffer(val byte) {
	//debug
	if d.memoryReader.bytesRemaining == 0 {
		panic(fmt.Sprint("d.bytesRemaining == 0"))
	}

	if d.memoryReader.bytesRemaining > 0 {
		d.memoryReader.sampleBuffer = append(d.memoryReader.sampleBuffer, val)

		// > The address is incremented; if it exceeds $FFFF, it is wrapped around to $8000.
		d.memoryReader.currentAddr++
		if d.memoryReader.currentAddr == 0 {
			d.memoryReader.currentAddr = 0x8000
		}

		// > The bytes remaining counter is decremented; if it becomes zero and the loop flag is set, the sample is restarted (see above);
		// > otherwise, if the bytes remaining counter becomes zero and the IRQ enabled flag is set, the interrupt flag is set.
		d.memoryReader.bytesRemaining--
		if d.memoryReader.bytesRemaining == 0 {
			if d.loop {
				d.memoryReader.restart()
			} else if d.irqEnabled {
				d.interruptFlag = true
			}
		}
	}
}

func (d *dmc) tickTimer() {
	if d.timer.tick() {
		// > When the timer outputs a clock, the following actions occur in order:

		// > 1. If the silence flag is clear, the output level changes based on bit 0 of the shift register.
		// >    If the bit is 1, add 2; otherwise, subtract 2. But if adding or subtracting 2 would cause the output level to leave the 0-127 range, leave the output level unchanged.
		// >    This means subtract 2 only if the current level is at least 2, or add 2 only if the current level is at most 125.
		if !d.outputUnit.silenceFlag {
			if (d.outputUnit.rightShiftRegister & 0x01) == 0x01 {
				if d.outputUnit.level <= 125 {
					d.outputUnit.level += 2
				}
			} else {
				if d.outputUnit.level >= 2 {
					d.outputUnit.level -= 2
				}
			}
		}
		// > 2. The right shift register is clocked.
		d.outputUnit.rightShiftRegister >>= 1

		// > 3. As stated above, the bits-remaining counter is decremented. If it becomes zero, a new output cycle is started.

		// > The bits-remaining counter is updated whenever the timer outputs a clock, regardless of whether a sample is currently playing.
		d.outputUnit.bitsRemainingCounter--
		if d.outputUnit.bitsRemainingCounter == 0 {
			// > When this counter reaches zero, we say that the output cycle ends.
			// > When an output cycle ends, a new cycle is started as follows:
			// > - The bits-remaining counter is loaded with 8.
			// > - If the sample buffer is empty, then the silence flag is set; otherwise, the silence flag is cleared and the sample buffer is emptied into the shift register.
			d.outputUnit.bitsRemainingCounter = 8
			if len(d.memoryReader.sampleBuffer) == 0 {
				d.outputUnit.silenceFlag = true
			} else {
				d.outputUnit.silenceFlag = false
				// うーん、そもそもbytesRemaining見てなくない？このtickTimer()のときに
				// output unit の中に bytesRemaining の概念がない => memory readerがbytesRemainingとcurrentAddressを保持してるため、のはず
				// つまり、memory reader側で(bytesRemaining/currentAddress/buffer)を見てfetchするかどうかの判断, ってことかな
				// memory reader ってstructを作った方が良い
				// > The DMC channel contains the following: memory reader, interrupt flag, sample buffer, timer, output unit, 7-bit output level with up and down counter.
				// memory reader: bytesRemaining, currentAddr, Sample address, Sample Length, sample buffer(?)
				// output unit: an 8-bit right shift register, a bits-remaining counter, a 7-bit output level (the same one that can be loaded directly via $4011), and a silence flag.
				d.outputUnit.rightShiftRegister = d.memoryReader.fetch()
				//d.dma.SignalDCMDMA(d.currentAddr)
			}
		}
	}
}
