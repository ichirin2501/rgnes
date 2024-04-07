package nes

// ref: https://www.nesdev.org/wiki/APU_Frame_Counter
// > The sequencer is clocked on every other CPU cycle, so 2 CPU cycles = 1 APU cycle.
// > ... (with an additional delay of one CPU cycle for the quarter and half frame signals).
// e.g. APU Cycle 3728.5 => Envelopes & triangle's linear counter
//
// ref: https://www.nesdev.org/wiki/APU#Glossary
// > The triangle channel's timer is clocked on every CPU cycle
//
// For the above two reasons, I will design the system to move the APU one step when the CPU moves one step.
// Furthermore, the frame sequence value written on the wiki(APU_Frame_Counter) is doubled and stored in the following frame table.

// NTSC
// 4-step seq: 0,1,2,.....29829,29830(0),1,2,...
var frameTable = [][]int{
	{7457, 14913, 22371, 29828, 29829, 29830}, // 4-step seq
	{7457, 14913, 22371, 29829, 37281, 37282}, // 5-step seq
}

// https://www.nesdev.org/wiki/APU_Mixer#Lookup_Table
var pulseTable [31]float32
var tndTable [203]float32

func init() {
	for i := 0; i < 31; i++ {
		pulseTable[i] = 95.52 / (8128.0/float32(i) + 100)
	}
	for i := 0; i < 203; i++ {
		tndTable[i] = 163.67 / (24329.0/float32(i) + 100)
	}
}

type Player interface {
	Sample(float32)
	SampleRate() float64
}

type APU struct {
	irqLine      *interruptLine
	player       Player
	sampleRate   float64
	sampleTiming int
	clock        int
	pulse1       *pulse
	pulse2       *pulse
	tnd          *triangle
	noise        *noise
	dmc          *dmc

	// frame counter
	frameMode              byte // Sequencer mode: 0 selects 4-step sequence, 1 selects 5-step sequence
	frameInterruptInhibit  bool
	frameInterruptFlag     bool
	frameStep              int
	newFrameCounterVal     int // for delay
	writeDelayFrameCounter byte
}

func NewAPU(irqLine *interruptLine, p Player, dma *DMA) *APU {
	apu := &APU{
		player:     p,
		irqLine:    irqLine,
		sampleRate: p.SampleRate(),

		clock:  -1,
		pulse1: newPulse(1),
		pulse2: newPulse(2),
		tnd:    newTriangle(),
		noise:  newNoise(),
		dmc:    newDMC(dma, irqLine),

		frameStep:          -1,
		newFrameCounterVal: -1,
	}
	apu.sampleTiming = int(CPUClockFrequency / apu.sampleRate)
	return apu
}

func (apu *APU) PowerUp() {
	apu.writeStatus(0)
	apu.noise.shiftRegister = 1
	apu.writeDMCController(0)
	apu.writeDMCLoadCounter(0)
	apu.writeDMCSampleAddr(0)
	apu.writeDMCSampleLength(0)
}

func (apu *APU) Reset() {
	apu.writeStatus(0)
	apu.tnd.seqPos = 0
}

func (apu *APU) Step() {
	apu.clock++

	// e.g. https://www.nesdev.org/wiki/APU_Pulse
	// Looking at the diagram output to mixer via Gate, there are dependencies of components.
	// For example, in Pulse, timer depends on sweep, so my understanding is that the value of sweep must be determined before timer.
	// So, call tickFrameCounter() before tickTimers()
	apu.tickFrameCounter()
	apu.tickTimers()

	if apu.clock%apu.sampleTiming == 0 {
		out := apu.output()
		apu.player.Sample(out)
	}
}

// DDLC VVVV	Duty (D), envelope loop / length counter halt (L), constant volume (C), volume/envelope (V)
func writePulseController(p *pulse, val byte) {
	p.duty = (val >> 6) & 0b11
	p.lc.halt = (val & 0x20) == 0x20
	p.el.loop = (val & 0x20) == 0x20
	p.el.constantVolume = (val & 0x10) == 0x10
	p.el.divider.period = uint16(val & 0x0F)
}

// https://www.nesdev.org/wiki/APU_Sweep
// EPPP NSSS	Sweep unit: enabled (E), period (P), negate (N), shift (S)
func writePulseSweep(p *pulse, val byte) {
	p.sweepEnabled = (val & 0x80) == 0x80
	// > The divider's period is P + 1 half-frames
	// The divider in this implementation works at a cycle of P+1 by default, so plus 1 is not necessary
	p.sweepDivider.period = (uint16((val >> 4) & 0b111))
	p.sweepNegate = (val & 0x08) == 0x08
	p.sweepShiftCount = val & 0b111
	// > Side effects	Sets the reload flag
	p.sweepReload = true
	p.updateTargetPeriod()
}

// TTTT TTTT	Timer low (T)
func writePulseTimerLow(p *pulse, val byte) {
	// timer is 11bit
	p.timer.period = (p.timer.period & 0x0700) | uint16(val)
	p.updateTargetPeriod()
}

// LLLL LTTT	Length counter load (L), timer high (T)
func writePulseLengthAndTimerHigh(p *pulse, val byte) {
	p.lc.load(val >> 3)
	p.timer.period = (p.timer.period & 0x00FF) | (uint16(val&0b111) << 8)
	p.updateTargetPeriod()

	// > The sequencer is immediately restarted at the first value of the current sequence.
	// > The envelope is also restarted.
	p.dutyPos = 0
	p.el.start = true
}

// $4000
func (apu *APU) writePulse1Controller(val byte) {
	writePulseController(apu.pulse1, val)
}

// $4001
func (apu *APU) writePulse1Sweep(val byte) {
	writePulseSweep(apu.pulse1, val)
}

// $4002
func (apu *APU) writePulse1TimerLow(val byte) {
	writePulseTimerLow(apu.pulse1, val)
}

// $4003
func (apu *APU) writePulse1LengthAndTimerHigh(val byte) {
	writePulseLengthAndTimerHigh(apu.pulse1, val)
}

// $4004
func (apu *APU) writePulse2Controller(val byte) {
	writePulseController(apu.pulse2, val)
}

// $4005
func (apu *APU) writePulse2Sweep(val byte) {
	writePulseSweep(apu.pulse2, val)
}

// $4006
func (apu *APU) writePulse2TimerLow(val byte) {
	writePulseTimerLow(apu.pulse2, val)
}

// $4007
func (apu *APU) writePulse2LengthAndTimerHigh(val byte) {
	writePulseLengthAndTimerHigh(apu.pulse2, val)
}

// $4008
// CRRR RRRR	Length counter halt / linear counter control (C), linear counter load (R)
func (apu *APU) writeTriangleController(val byte) {
	apu.tnd.lc.halt = (val & 0x80) == 0x80
	apu.tnd.linearCounterCtrl = (val & 0x80) == 0x80
	apu.tnd.linearCounterPeriod = val & 0x7F
}

// $400A
// TTTT TTTT	Timer low (T)
func (apu *APU) writeTriangleTimerLow(val byte) {
	// 11bit
	apu.tnd.timer.period = (apu.tnd.timer.period & 0x0700) | uint16(val)
}

// $400B
// LLLL LTTT	Length counter load (L), timer high (T)
func (apu *APU) writeTriangleLengthAndTimerHigh(val byte) {
	apu.tnd.lc.load(val >> 3)
	apu.tnd.timer.period = (apu.tnd.timer.period & 0x00FF) | (uint16(val&0b111) << 8)
	apu.tnd.linearCounterReload = true
}

// $400C
// --LC VVVV	el loop / length counter halt (L), constant volume (C), volume/envelope (V)
func (apu *APU) writeNoiseController(val byte) {
	apu.noise.lc.halt = (val & 0x20) == 0x20
	apu.noise.el.loop = (val & 0x20) == 0x20
	apu.noise.el.constantVolume = (val & 0x10) == 0x10
	apu.noise.el.divider.period = uint16(val & 0x0F)
}

// $400E
// M---.PPPP	Mode and period (write)
func (apu *APU) writeNoiseLoopAndPeriod(val byte) {
	apu.noise.mode = (val & 0x80) == 0x80
	apu.noise.loadPeriod(val & 0x0F)
}

// $400F
// LLLL L---	Length counter load (L)
func (apu *APU) writeNoiseLength(val byte) {
	apu.noise.lc.load(val >> 3)
	apu.noise.el.start = true
}

// $4010
// IL-- RRRR	IRQ enable (I), loop (L), frequency (R)
func (apu *APU) writeDMCController(val byte) {
	apu.dmc.irqEnabled = (val & 0x80) == 0x80
	if !apu.dmc.irqEnabled {
		apu.dmc.clearInterruptFlag()
	}
	apu.dmc.loop = (val & 0x40) == 0x40
	apu.dmc.loadRate(val & 0x0F)
}

// $4011
// -DDD DDDD	load counter (D)
func (apu *APU) writeDMCLoadCounter(val byte) {
	apu.dmc.level = val & 0x7F
}

// $4012
// AAAA AAAA	Sample address (A)
func (apu *APU) writeDMCSampleAddr(val byte) {
	apu.dmc.sampleAddr = 0xC000 + (uint16(val) * 64)
}

// $4013
// LLLL LLLL	Sample length (L)
func (apu *APU) writeDMCSampleLength(val byte) {
	apu.dmc.sampleLength = (uint16(val) * 16) + 1
}

// $4015 read
// IF-D NT21	DMC interrupt (I), frame interrupt (F), DMC active (D), length counter > 0 (N/T/2/1)
func (apu *APU) readStatus() byte {
	res := byte(0)
	if apu.pulse1.lc.value > 0 {
		res |= 0x01
	}
	if apu.pulse2.lc.value > 0 {
		res |= 0x02
	}
	if apu.tnd.lc.value > 0 {
		res |= 0x04
	}
	if apu.noise.lc.value > 0 {
		res |= 0x08
	}
	if apu.dmc.bytesRemaining > 0 {
		res |= 0x10
	}
	if apu.frameInterruptFlag {
		res |= 0x40
	}
	if apu.dmc.interruptFlag {
		res |= 0x80
	}
	apu.frameInterruptFlag = false
	apu.irqLine.SetHigh()

	return res
}

// PeekStatus is used for debugging
func (apu *APU) PeekStatus() byte {
	res := byte(0)
	if apu.pulse1.lc.value > 0 {
		res |= 0x01
	}
	if apu.pulse2.lc.value > 0 {
		res |= 0x02
	}
	if apu.tnd.lc.value > 0 {
		res |= 0x04
	}
	if apu.noise.lc.value > 0 {
		res |= 0x08
	}
	if apu.dmc.bytesRemaining > 0 {
		res |= 0x10
	}
	if apu.frameInterruptFlag {
		res |= 0x40
	}
	if apu.dmc.interruptFlag {
		res |= 0x80
	}
	return res
}

// $4015 write
// ---D NT21	Enable DMC (D), noise (N), triangle (T), and pulse channels (2/1)
func (apu *APU) writeStatus(val byte) {
	apu.dmc.setEnabled((val & 0x10) == 0x10)
	apu.dmc.clearInterruptFlag()
	apu.noise.lc.setEnabled((val & 0x08) == 0x08)
	apu.tnd.lc.setEnabled((val & 0x04) == 0x04)
	apu.pulse2.lc.setEnabled((val & 0x02) == 0x02)
	apu.pulse1.lc.setEnabled((val & 0x01) == 0x01)
}

// $4017
func (apu *APU) writeFrameCounter(val byte) {
	// ref: https://www.nesdev.org/wiki/APU#Frame_Counter_($4017)
	// > Writing to $4017 resets the frame counter and the quarter/half frame triggers happen simultaneously,
	// > but only on "odd" cycles (and only after the first "even" cycle after the write occurs)
	// > - thus, it happens either 2 or 3 cycles after the write (i.e. on the 2nd or 3rd cycle of the next instruction).
	// ref: https://www.nesdev.org/wiki/APU_Frame_Counter
	// > * If the write occurs during an APU cycle, the effects occur 3 CPU cycles after the $4017 write cycle,
	// > and if the write occurs between APU cycles, the effects occurs 4 CPU cycles after the write cycle.

	// ref: https://forums.nesdev.org/viewtopic.php?t=454
	// > The APU's master clock is at 1.79 MHz, same as the CPU clock. This can be divided into even and odd clocks
	// > Quick summary: A write to $4017 changes the APU mode (and restart it). Depending on when the write occurs,
	// > the mode change might be delayed by a single CPU clock, as if you wrote to $4017 one clock later.

	// Depending on a page, it was written as "2 or 3 cycles" or "3 or 4 cycles", so I couldn't understand it.
	// I'm not sure if the table below is correct, but my understanding is as follows
	//       apu.clock:   0,   1,   2,   3,   4,   5,   6,   7, ...
	// real APU Cycles:   0,   0,   1,   1,   2,   2,   3,   3, ...
	//                  0.0, 0.5, 1.0, 1.5, 2.0, 2.0, 2.5, 3.0, ...
	//         trigger:   t,    ,   t,    ,   t,    ,   t,    , ...
	//        even/odd:   e,   o,   e,   o,   e,   o,   e,   o, ...
	// Determine even/odd cycles with APU's master clock(=apu.clock) and delay by a single CPU.
	// And, the frame counter(sequencer) is clocked on every other CPU cycle(2 CPU cycles = 1 APU cycle, ^ trigger row in the above table).
	// I decided to adjust it according to the above timing of the trigger.

	apu.newFrameCounterVal = int(val)
	if apu.clock%2 == 0 {
		apu.writeDelayFrameCounter = 2
	} else {
		apu.writeDelayFrameCounter = 3
	}
	if (val & 0x40) == 0x40 {
		apu.frameInterruptInhibit = true
		apu.frameInterruptFlag = false
		apu.irqLine.SetHigh()
	} else {
		apu.frameInterruptInhibit = false
	}
}

// https://www.nesdev.org/wiki/APU_Mixer#Lookup_Table
// > output = pulse_out + tnd_out
// > pulse_out = pulse_table [pulse1 + pulse2]
// > tnd_out = tnd_table [3 * triangle + 2 * noise + dmc]
// > The values for pulse1, pulse2, triangle, noise, and dmc are the output values for the corresponding channel.
// > The dmc value ranges from 0 to 127 and the others range from 0 to 15.
func (apu *APU) output() float32 {
	pout := pulseTable[apu.pulse1.output()+apu.pulse2.output()]
	tout := tndTable[3*apu.tnd.output()+2*apu.noise.output()+apu.dmc.output()]
	return pout + tout
}

func (apu *APU) tickTimers() {
	// > The triangle channel's timer is clocked on every CPU cycle,
	// > but the pulse, noise, and DMC timers are clocked only on every second CPU cycle and thus produce only even periods
	if apu.clock%2 == 0 {
		apu.pulse1.tickTimer()
		apu.pulse2.tickTimer()
		//apu.noise.tickTimer()
		//apu.dmc.tickTimer()
	}
	// Since the DMC/Noise tables are defined in units of CPU cycles, I will run the timers every time for now
	apu.noise.tickTimer()
	apu.dmc.tickTimer()
	apu.tnd.tickTimer()
}

func (apu *APU) tickQuarterFrameCounter() {
	apu.pulse1.tickEnvelope()
	apu.pulse2.tickEnvelope()
	apu.noise.tickEnvelope()
	apu.tnd.tickLinearCounter()
}

func (apu *APU) tickHalfFrameCounter() {
	apu.pulse1.tickLengthCounter()
	apu.pulse2.tickLengthCounter()
	apu.tnd.tickLengthCounter()
	apu.noise.tickLengthCounter()
	apu.pulse1.tickSweep()
	apu.pulse2.tickSweep()
}

// debug
func (apu *APU) FetchFrameStep() int {
	return apu.frameStep
}
func (apu *APU) FetchFrameMode() int {
	return int(apu.frameMode)
}
func (apu *APU) FetchPulse1LC() int {
	return int(apu.pulse1.lc.value)
}
func (apu *APU) FetchFrameIRQFlag() bool {
	return apu.frameInterruptFlag
}
func (apu *APU) FetchNewFrameCounterVal() int {
	return apu.newFrameCounterVal
}
func (apu *APU) FetchWriteDelayFC() byte {
	return apu.writeDelayFrameCounter
}

func (apu *APU) resetFrameCounter() {
	apu.frameStep = 0
}

func (apu *APU) tickFrameCounter() {
	apu.frameStep++

	if apu.newFrameCounterVal >= 0 {
		if apu.writeDelayFrameCounter > 0 {
			apu.writeDelayFrameCounter--
		} else {
			apu.resetFrameCounter()
			if (apu.newFrameCounterVal & 0x80) == 0x80 {
				apu.frameMode = 1
				apu.tickHalfFrameCounter()
				apu.tickQuarterFrameCounter()
			} else {
				apu.frameMode = 0
			}
			apu.writeDelayFrameCounter = 0
			apu.newFrameCounterVal = -1
		}
	}

	if apu.frameMode == 0 {
		// 4 step
		switch apu.frameStep {
		case frameTable[apu.frameMode][0]:
			apu.tickQuarterFrameCounter()
		case frameTable[apu.frameMode][1]:
			apu.tickQuarterFrameCounter()
			apu.tickHalfFrameCounter()
		case frameTable[apu.frameMode][2]:
			apu.tickQuarterFrameCounter()
		case frameTable[apu.frameMode][3]:
			if !apu.frameInterruptInhibit {
				apu.frameInterruptFlag = true
				apu.irqLine.SetLow()
			}
		case frameTable[apu.frameMode][4]:
			apu.tickQuarterFrameCounter()
			apu.tickHalfFrameCounter()
			if !apu.frameInterruptInhibit {
				apu.frameInterruptFlag = true
				apu.irqLine.SetLow()
			}
		case frameTable[apu.frameMode][5]:
			if !apu.frameInterruptInhibit {
				apu.frameInterruptFlag = true
				apu.irqLine.SetLow()
			}
			apu.frameStep = 0
		}
	} else {
		// 5 step
		switch apu.frameStep {
		case frameTable[apu.frameMode][0]:
			apu.tickQuarterFrameCounter()
		case frameTable[apu.frameMode][1]:
			apu.tickQuarterFrameCounter()
			apu.tickHalfFrameCounter()
		case frameTable[apu.frameMode][2]:
			apu.tickQuarterFrameCounter()
		case frameTable[apu.frameMode][3]:
			// nothing
		case frameTable[apu.frameMode][4]:
			apu.tickQuarterFrameCounter()
			apu.tickHalfFrameCounter()
		case frameTable[apu.frameMode][5]:
			apu.frameStep = 0
		}
	}
}

type timer struct {
	divider
}

type divider struct {
	counter uint16
	period  uint16
}

func (d *divider) tick() bool {
	if d.counter == 0 {
		d.reload()
		return true
	} else {
		d.counter--
	}
	return false
}

func (d *divider) reload() {
	d.counter = d.period
}
