package apu

// https://www.nesdev.org/wiki/APU_Frame_Counter
// > The sequencer is clocked on every other CPU cycle, so 2 CPU cycles = 1 APU cycle.
// > The sequencer keeps track of how many APU cycles have elapsed in total,
// > and each step of the sequence will occur once that total has reached the indicated amount (with an additional delay of one CPU cycle for the quarter and half frame signals).
// > Once the last step has executed, the count resets to 0 on the next APU cycle.
// apu cyclesを保持するとあるけど、timerはcpu cycleで動くので、ここでもcpu cyclesで統一してテーブルを用意する
// NTSC
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

type CPU interface {
	SetIRQ(val bool)
}

type Player interface {
	Sample(float32)
}

type Option func(*APU)

type APU struct {
	cpu          CPU
	player       Player
	sampleRate   int
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
	frameSequenceStep      int
	newFrameCounterVal     int // for delay
	writeDelayFrameCounter byte
}

func New(cpu CPU, p Player, opts ...Option) *APU {
	apu := &APU{
		player:     p,
		cpu:        cpu,
		sampleRate: 44100,

		clock:  -1,
		pulse1: newPulse(1),
		pulse2: newPulse(2),
		tnd:    newTriangle(),
		noise:  newNoise(),
		dmc:    newDMC(),

		frameStep:          -1,
		newFrameCounterVal: -1,
	}
	for _, opt := range opts {
		opt(apu)
	}
	apu.sampleTiming = 1789773 / apu.sampleRate
	return apu
}

func WithSampleRate(sampleRate int) Option {
	return func(apu *APU) {
		apu.sampleRate = sampleRate
	}
}

func (apu *APU) PowerUp() {
	// todo
	apu.WriteStatus(0)
	apu.noise.shiftRegister = 1
}

func (apu *APU) Reset() {
	// todo
	apu.WriteStatus(0)
}

func (apu *APU) Step() {
	apu.clock++
	apu.tickTimers()
	apu.tickFrameCounter()

	// test
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
	p.sweepDivider.period = (uint16((val >> 4) & 0b111)) + 1
	p.sweepNegate = (val & 0x04) == 0x04
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
func (apu *APU) WritePulse1Controller(val byte) {
	writePulseController(apu.pulse1, val)
}

// $4001
func (apu *APU) WritePulse1Sweep(val byte) {
	writePulseSweep(apu.pulse1, val)
}

// $4002
func (apu *APU) WritePulse1TimerLow(val byte) {
	writePulseTimerLow(apu.pulse1, val)
}

// $4003
func (apu *APU) WritePulse1LengthAndTimerHigh(val byte) {
	writePulseLengthAndTimerHigh(apu.pulse1, val)
}

// $4004
func (apu *APU) WritePulse2Controller(val byte) {
	writePulseController(apu.pulse2, val)
}

// $4005
func (apu *APU) WritePulse2Sweep(val byte) {
	writePulseSweep(apu.pulse2, val)
}

// $4006
func (apu *APU) WritePulse2TimerLow(val byte) {
	writePulseTimerLow(apu.pulse2, val)
}

// $4007
func (apu *APU) WritePulse2LengthAndTimerHigh(val byte) {
	writePulseLengthAndTimerHigh(apu.pulse2, val)
}

// $4008
// CRRR RRRR	Length counter halt / linear counter control (C), linear counter load (R)
func (apu *APU) WriteTriangleController(val byte) {
	apu.tnd.lc.halt = (val & 0x80) == 0x80
	apu.tnd.linearCounterCtrl = (val & 0x80) == 0x80
	apu.tnd.linearCounterPeriod = val & 0x7F
}

// $400A
// TTTT TTTT	Timer low (T)
func (apu *APU) WriteTriangleTimerLow(val byte) {
	// 11bit
	apu.tnd.timer.period = (apu.tnd.timer.period & 0x0700) | uint16(val)
}

// $400B
// LLLL LTTT	Length counter load (L), timer high (T)
func (apu *APU) WriteTriangleLengthAndTimerHigh(val byte) {
	apu.tnd.lc.load(val >> 3)
	apu.tnd.timer.period = (apu.tnd.timer.period & 0x00FF) | (uint16(val&0b111) << 8)
	apu.tnd.linearCounterReload = true
}

// $400C
// --LC VVVV	el loop / length counter halt (L), constant volume (C), volume/envelope (V)
func (apu *APU) WriteNoiseController(val byte) {
	apu.noise.lc.halt = (val & 0x20) == 0x20
	apu.noise.el.loop = (val & 0x20) == 0x20
	apu.noise.el.constantVolume = (val & 0x10) == 0x10
	apu.noise.el.divider.period = uint16(val & 0x0F)
}

// $400E
// M---.PPPP	Mode and period (write)
func (apu *APU) WriteNoiseLoopAndPeriod(val byte) {
	apu.noise.mode = (val & 0x80) == 0x80
	apu.noise.loadPeriod(val & 0x0F)
}

// $400F
// LLLL L---	Length counter load (L)
func (apu *APU) WriteNoiseLength(val byte) {
	apu.noise.lc.load(val >> 3)
	apu.noise.el.start = true
}

// $4010
// IL-- RRRR	IRQ enable (I), loop (L), frequency (R)
func (apu *APU) WriteDMCController(val byte) {
	apu.dmc.irqEnabled = (val & 0x80) == 0x80
	apu.dmc.loop = (val & 0x40) == 0x40
	apu.dmc.freq = val & 0x0F
}

// $4011
// -DDD DDDD	load counter (D)
func (apu *APU) WriteDMCLoadCounter(val byte) {
	apu.dmc.counter = val & 0x7F
}

// $4012
// AAAA AAAA	Sample address (A)
func (apu *APU) WriteDMCSampleAddr(val byte) {
	apu.dmc.sampleAddr = val
}

// $4013
// LLLL LLLL	Sample length (L)
func (apu *APU) WriteDMCSampleLength(val byte) {
	apu.dmc.sampleLength = val
}

// $4015 read
// IF-D NT21	DMC interrupt (I), frame interrupt (F), DMC active (D), length counter > 0 (N/T/2/1)
func (apu *APU) ReadStatus() byte {
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
	if apu.frameInterruptFlag {
		res |= 0x40
	}
	apu.frameInterruptFlag = false
	apu.cpu.SetIRQ(false)

	// todo: dmc, F

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
	if apu.frameInterruptFlag {
		res |= 0x40
	}
	// todo
	return res
}

// $4015 write
// ---D NT21	Enable DMC (D), noise (N), triangle (T), and pulse channels (2/1)
func (apu *APU) WriteStatus(val byte) {
	apu.dmc.setEnabled((val & 0x10) == 0x10)
	apu.noise.lc.setEnabled((val & 0x08) == 0x08)
	apu.tnd.lc.setEnabled((val & 0x04) == 0x04)
	apu.pulse2.lc.setEnabled((val & 0x02) == 0x02)
	apu.pulse1.lc.setEnabled((val & 0x01) == 0x01)
}

// $4017
func (apu *APU) WriteFrameCounter(val byte) {
	// https://www.nesdev.org/wiki/APU#Frame_Counter_($4017)
	// > Writing to $4017 resets the frame counter and the quarter/half frame triggers happen simultaneously,
	// > but only on "odd" cycles (and only after the first "even" cycle after the write occurs)
	// > - thus, it happens either 2 or 3 cycles after the write (i.e. on the 2nd or 3rd cycle of the next instruction).
	// > After 2 or 3 clock cycles (depending on when the write is performed), the timer is reset.
	// 2 or 3 って書いてるけど、別のFrameCounterのページ見たら 3 or 4 って書いてたりしてよくわからん
	// テストROMで適当に試したら3 or 4っぽかったのでこれでfixとする
	apu.newFrameCounterVal = int(val)
	if apu.clock%2 == 0 {
		apu.writeDelayFrameCounter = 3
	} else {
		apu.writeDelayFrameCounter = 4
	}
	if (val & 0x40) == 0x40 {
		apu.frameInterruptInhibit = true
		apu.frameInterruptFlag = false
		apu.cpu.SetIRQ(false)
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
	apu.pulse1.tickTimer()
	apu.pulse2.tickTimer()
	apu.tnd.tickTimer()
	apu.noise.tickTimer()
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
func (apu *APU) FetchFrameSeqStep() int {
	return apu.frameSequenceStep
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
	apu.frameSequenceStep = 0
}

func (apu *APU) tickFrameCounter() {
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

	apu.frameStep++
	if apu.frameStep >= frameTable[apu.frameMode][apu.frameSequenceStep] {
		if apu.frameMode == 0 {
			// 4 step
			switch apu.frameSequenceStep {
			case 0:
				apu.tickQuarterFrameCounter()
			case 1:
				apu.tickQuarterFrameCounter()
				apu.tickHalfFrameCounter()
			case 2:
				apu.tickQuarterFrameCounter()
			case 3:
				if !apu.frameInterruptInhibit {
					apu.frameInterruptFlag = true
					apu.cpu.SetIRQ(true)
				}
			case 4:
				apu.tickQuarterFrameCounter()
				apu.tickHalfFrameCounter()
				if !apu.frameInterruptInhibit {
					apu.frameInterruptFlag = true
					apu.cpu.SetIRQ(true)
				}
			case 5:
				if !apu.frameInterruptInhibit {
					apu.frameInterruptFlag = true
					apu.cpu.SetIRQ(true)
				}
			}
		} else {
			// 5 step
			switch apu.frameSequenceStep {
			case 0:
				apu.tickQuarterFrameCounter()
			case 1:
				apu.tickQuarterFrameCounter()
				apu.tickHalfFrameCounter()
			case 2:
				apu.tickQuarterFrameCounter()
			case 4:
				apu.tickQuarterFrameCounter()
				apu.tickHalfFrameCounter()
			}
		}

		apu.frameSequenceStep++
		if apu.frameSequenceStep >= 6 {
			apu.resetFrameCounter()
		}
	}

}

type timer struct {
	divider
}

func newTimer(factor uint16) timer {
	return timer{
		divider: divider{
			// factor is used internally with +1
			factor: factor - 1,
		},
	}
}

type divider struct {
	counter uint16
	period  uint16
	// > The triangle channel's timer is clocked on every CPU cycle,
	// > but the pulse, noise, and DMC timers are clocked only on every second CPU cycle and thus produce only even periods
	factor uint16
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
	d.counter = d.period * (d.factor + 1)
}
