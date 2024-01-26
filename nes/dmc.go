package nes

// https://www.nesdev.org/wiki/APU_DMC
// The dmc timer is clocked only on every second CPU cycle
var dmcPeriodTable = []uint16{
	214, 190, 170, 160, 143, 127, 113, 107, 95, 80, 71, 64, 53, 42, 36, 27,
}

type dmc struct {
	enabled       bool
	irqEnabled    bool
	interruptFlag bool
	loop          bool
	rateIndex     byte
	timer         timer

	// memory reader
	sampleAddr     uint16
	currentAddr    uint16
	sampleLength   uint16
	bytesRemaining uint16
	sampleBuffer   []byte
	dma            *DMA

	// output unit
	rightShiftRegister   byte
	bitsRemainingCounter byte
	silenceFlag          bool
	level                byte // 7 bit
}

func newDMC(dma *DMA) *dmc {
	return &dmc{
		sampleBuffer: make([]byte, 0, 1),
		dma:          dma,
	}
}

func (d *dmc) setEnabled(v bool) {
	if !v {
		d.enabled = false
		// > If the DMC bit is clear, the DMC bytes remaining will be set to 0 and the DMC will silence when it empties.
		d.bytesRemaining = 0
	} else {
		d.enabled = true
		// > If the DMC bit is set, the DMC sample will be restarted only if its bytes remaining is 0.
		// > If there are bits remaining in the 1-byte sample buffer, these will finish playing before the next sample is fetched.
		if d.bytesRemaining == 0 {
			d.restart()
		}
		d.SignalDMA(false)
	}
}

func (d *dmc) restart() {
	d.bytesRemaining = d.sampleLength
	d.currentAddr = d.sampleAddr
}

func (d *dmc) SignalDMA(reloadTiming bool) {
	if d.enabled && d.bytesRemaining > 0 && len(d.sampleBuffer) == 0 {
		if reloadTiming {
			d.dma.TriggerOnDCMReload(d.currentAddr)
		} else {
			d.dma.TriggerOnDCMLoad(d.currentAddr)
		}
	}
}

func (d *dmc) output() byte {
	return d.level
}

func (d *dmc) loadRate(rateIndex byte) {
	d.timer.period = dmcPeriodTable[rateIndex]
}

func (d *dmc) setSampleBuffer(val byte) {
	if d.bytesRemaining > 0 {
		d.sampleBuffer = append(d.sampleBuffer, val)

		// > The address is incremented; if it exceeds $FFFF, it is wrapped around to $8000.
		d.currentAddr++
		if d.currentAddr == 0 {
			d.currentAddr = 0x8000
		}

		// > The bytes remaining counter is decremented; if it becomes zero and the loop flag is set, the sample is restarted (see above);
		// > otherwise, if the bytes remaining counter becomes zero and the IRQ enabled flag is set, the interrupt flag is set.
		d.bytesRemaining--
		if d.bytesRemaining == 0 {
			if d.loop {
				d.restart()
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
		if !d.silenceFlag {
			if (d.rightShiftRegister & 0x01) == 0x01 {
				if d.level <= 125 {
					d.level += 2
				}
			} else {
				if d.level >= 2 {
					d.level -= 2
				}
			}
		}
		// > 2. The right shift register is clocked.
		d.rightShiftRegister >>= 1

		// > 3. As stated above, the bits-remaining counter is decremented. If it becomes zero, a new output cycle is started.

		// > The bits-remaining counter is updated whenever the timer outputs a clock, regardless of whether a sample is currently playing.
		d.bitsRemainingCounter--
		if d.bitsRemainingCounter == 0 {
			// > When this counter reaches zero, we say that the output cycle ends.
			// > When an output cycle ends, a new cycle is started as follows:
			// > - The bits-remaining counter is loaded with 8.
			// > - If the sample buffer is empty, then the silence flag is set; otherwise, the silence flag is cleared and the sample buffer is emptied into the shift register.
			d.bitsRemainingCounter = 8
			if len(d.sampleBuffer) == 0 {
				d.silenceFlag = true
			} else {
				d.silenceFlag = false
				d.rightShiftRegister = d.sampleBuffer[0]
				d.sampleBuffer = d.sampleBuffer[:0]
				d.SignalDMA(true)
			}
		}
	}
}
