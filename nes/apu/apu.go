package apu

// https://www.nesdev.org/wiki/APU_Pulse
var dutyTable = [][]byte{
	{0, 1, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 0, 0, 0},
	{1, 0, 0, 1, 1, 1, 1, 1},
}

// https://www.nesdev.org/wiki/APU_Triangle
var triangleTable = []byte{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
}

// https://www.nesdev.org/wiki/APU_Noise
// NTSC
var noisePeriodTable = []uint16{
	4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068,
}

// https://www.nesdev.org/wiki/APU_Length_Counter
var lengthTable = []byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

// https://www.nesdev.org/wiki/APU_DMC
// > These periods are all even numbers because there are 2 CPU cycles in an APU cycle.
// > A rate of 428 means the output level changes every 214 APU cycles.
// NTSC
var dmcPeriodTable = []uint16{
	428, 380, 340, 320, 286, 254, 226, 214, 190, 160, 142, 128, 106, 84, 72, 54,
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

type APU struct {
	pulse1 pulse
	pulse2 pulse
	tnd    triangle
	noise  noise
	dmc    dmc
}

func New() *APU {
	return &APU{}
}

// DDLC VVVV	Duty (D), envelope loop / length counter halt (L), constant volume (C), volume/envelope (V)
func writePulseController(p *pulse, val byte) {
	p.duty = (val >> 6) & 0b11

	// envelope
	p.constantVolume = (val & 0x10) == 0x10
	p.volume = val & 0x0F

	// length counter
	p.lengthCounterHalt = (val & 0x20) == 0x20
}

// https://www.nesdev.org/wiki/APU_Sweep
// EPPP NSSS	Sweep unit: enabled (E), period (P), negate (N), shift (S)
func writePulseSweep(p *pulse, val byte) {
	p.sweepEnabled = (val & 0x80) == 0x80
	// > The divider's period is P + 1 half-frames
	p.sweepPeriod = ((val >> 4) & 0b111) + 1
	p.sweepNegate = (val & 0x04) == 0x04
	p.sweepShiftCount = val & 0b111
	// > Side effects	Sets the reload flag
	p.sweepReload = true
}

// TTTT TTTT	Timer low (T)
func writePulseTimerLow(p *pulse, val byte) {
	p.timerPeriod = (p.timerPeriod & 0xFF00) | uint16(val)
}

// LLLL LTTT	Length counter load (L), timer high (T)
func writePulseLengthAndTimerHigh(p *pulse, val byte) {
	// todo: > If the enabled flag is set, the length counter is loaded with entry L of the length table
	p.lengthCounter = lengthTable[val>>3]
	p.timerPeriod = (p.timerPeriod & 0x00FF) | (uint16(val&0b111) << 8)

	// > The sequencer is immediately restarted at the first value of the current sequence.
	p.dutyPos = 0
}

// $4000
func (apu *APU) WritePulse1Controller(val byte) {
	writePulseController(&apu.pulse1, val)
}

// $4001
func (apu *APU) WritePulse1Sweep(val byte) {
	writePulseSweep(&apu.pulse1, val)
}

// $4002
func (apu *APU) WritePulse1TimerLow(val byte) {
	writePulseTimerLow(&apu.pulse1, val)
}

// $4003
func (apu *APU) WritePulse1LengthAndTimerHigh(val byte) {
	writePulseLengthAndTimerHigh(&apu.pulse1, val)
}

// $4004
func (apu *APU) WritePulse2Controller(val byte) {
	writePulseController(&apu.pulse2, val)
}

// $4005
func (apu *APU) WritePulse2Sweep(val byte) {
	writePulseSweep(&apu.pulse2, val)
}

// $4006
func (apu *APU) WritePulse2TimerLow(val byte) {
	writePulseTimerLow(&apu.pulse2, val)
}

// $4007
func (apu *APU) WritePulse2LengthAndTimerHigh(val byte) {
	writePulseLengthAndTimerHigh(&apu.pulse2, val)
}

// $4008
// CRRR RRRR	Length counter halt / linear counter control (C), linear counter load (R)
func (apu *APU) WriteTriangleController(val byte) {
	apu.tnd.lengthCounterHalt = (val & 0x80) == 0x80
	apu.tnd.linearCounter = val & 0x7F
}

// $400A
// TTTT TTTT	Timer low (T)
func (apu *APU) WriteTriangleTimerLow(val byte) {
	apu.tnd.timerPeriod = (apu.tnd.timerPeriod & 0xFF00) | uint16(val)
}

// $400B
// LLLL LTTT	Length counter load (L), timer high (T)
func (apu *APU) WriteTriangleLengthAndTimerHigh(val byte) {
	apu.tnd.lengthCounter = lengthTable[val>>3]
	apu.tnd.timerPeriod = (apu.tnd.timerPeriod & 0x00FF) | (uint16(val&0b111) << 8)
	apu.tnd.linearCounterReload = true
}

// $400C
// --LC VVVV	Envelope loop / length counter halt (L), constant volume (C), volume/envelope (V)
func (apu *APU) WriteNoiseController(val byte) {
	apu.noise.lengthCounterHalt = (val & 0x20) == 0x20
	apu.noise.constantVolume = (val & 0x10) == 0x10
	apu.noise.volume = val & 0x0F
}

// $400E
// L--- PPPP	Loop noise (L), noise period (P)
func (apu *APU) WriteNoiseLoopAndPeriod(val byte) {
	apu.noise.loop = (val & 0x80) == 0x80
	apu.noise.period = val & 0x0F
}

// $400F
// LLLL L---	Length counter load (L)
func (apu *APU) WriteNoiseLength(val byte) {
	apu.noise.lengthCounter = lengthTable[val>>3]
}

// $4010
// IL-- RRRR	IRQ enable (I), loop (L), frequency (R)
func (apu *APU) WriteDMCController(val byte) {
	apu.dmc.irqEnabled = (val & 0x80) == 0x80
	apu.dmc.loop = (val & 0x40) == 0x40
	apu.dmc.freq = val & 0x0F
}

// $4011
// -DDD DDDD	Load counter (D)
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
	if apu.pulse1.lengthCounter > 0 {
		res |= 0x01
	}
	if apu.pulse2.lengthCounter > 0 {
		res |= 0x02
	}
	if apu.tnd.lengthCounter > 0 {
		res |= 0x04
	}
	if apu.noise.lengthCounter > 0 {
		res |= 0x08
	}
	// todo: dmc, I, F

	return res
}

// PeekStatus is used for debugging
func (apu *APU) PeekStatus() byte {
	return apu.ReadStatus()
}

// $4015 write
// ---D NT21	Enable DMC (D), noise (N), triangle (T), and pulse channels (2/1)
func (apu *APU) WriteStatus(val byte) {
	apu.dmc.enabled = (val & 0x10) == 0x10
	apu.noise.enabled = (val & 0x08) == 0x08
	apu.tnd.enabled = (val & 0x04) == 0x04
	apu.pulse2.enabled = (val & 0x02) == 0x02
	apu.pulse1.enabled = (val & 0x01) == 0x01
}

// $4017
func (apu *APU) WriteFrameCounter(val byte) {}

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

type pulse struct {
	enabled bool
	// length counter
	lengthCounterHalt bool
	lengthCounter     byte

	// envelope
	// todo: envelop unit自体もdividerを持つ、かつ、noiseでも出てくるので、構造体にしたほうが良さそう?
	constantVolume bool
	volume         byte // volume/envelope (V)

	// https://www.nesdev.org/wiki/APU_Sweep
	sweepEnabled    bool
	sweepPeriod     byte
	sweepNegate     bool
	sweepShiftCount byte
	// > Each sweep unit contains the following: divider, reload flag.
	sweepDivider byte
	sweepReload  bool

	duty        byte
	dutyPos     byte
	timerPeriod uint16
}

func (p *pulse) output() byte {
	// todo
	return 0
}

type triangle struct {
	enabled bool
	// length counter
	lengthCounterHalt bool
	lengthCounter     byte

	linearCounter       byte
	linearCounterReload bool

	timerPeriod uint16
}

func (t *triangle) output() byte {
	// todo
	return 0
}

type noise struct {
	enabled bool
	// length counter
	lengthCounterHalt bool
	lengthCounter     byte

	// envelope
	constantVolume bool
	volume         byte // volume/envelope (V)

	loop        bool
	period      byte
	timerPeriod uint16
}

func (t *noise) output() byte {
	// todo
	return 0
}

type dmc struct {
	enabled      bool
	irqEnabled   bool
	loop         bool
	freq         byte
	counter      byte
	sampleAddr   byte
	sampleLength byte
}

func (d *dmc) output() byte {
	// todo
	return 0
}
