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

type APU struct {
	pulse1 pulse
	pulse2 pulse
}

func New() *APU {
	return &APU{}
}

// TODO
func (p *APU) Read(addr uint16) byte {
	switch addr {
	case 0x0015:
	}
	return 0
}

// TODO
func (p *APU) Write(addr uint16, val byte) {
	switch addr {
	default:
	}
}

// $4000
func (apu *APU) WritePulse1Controller(val byte) {

}

// $4001
func (apu *APU) WritePulse1Sweep(val byte) {
	apu.pulse1.sweep = sweepUnit(val)
}

// $4002
func (apu *APU) WritePulse1TimerLow(val byte) {
	apu.pulse1.timerPeriod = (apu.pulse1.timerPeriod & 0xFF00) | uint16(val)
}

// $4003
func (apu *APU) WritePulse1LengthAndTimerHigh(val byte) {
	// todo len
	// apu.pulse1.timerPeriod = (apu.pulse1.timerPeriod & 0x00FF) | (uint16(val&0b111) << 8)
}

// $4004
func (apu *APU) WritePulse2Controller(val byte) {

}

// $4005
func (apu *APU) WritePulse2Sweep(val byte) {
	apu.pulse2.sweep = sweepUnit(val)
}

// $4006
func (apu *APU) WritePulse2TimerLow(val byte) {
	apu.pulse2.timerPeriod = (apu.pulse2.timerPeriod & 0xFF00) | uint16(val)
}

// $4007
func (apu *APU) WritePulse2LengthAndTimerHigh(val byte) {

}

// $4008
func (apu *APU) WriteTriangleController(val byte) {}

// $400A
func (apu *APU) WriteTriangleTimerLow(val byte) {}

// $400B
func (apu *APU) WriteTriangleLengthAndTimerHigh(val byte) {}

// $400C
func (apu *APU) WriteNoiseController(val byte) {}

// $400E
func (apu *APU) WriteNoiseLoopAndPeriod(val byte) {}

// $400F
func (apu *APU) WriteNoiseLength(val byte) {}

// $4010
func (apu *APU) WriteDMCController(val byte) {}

// $4011
func (apu *APU) WriteDMCLoadCounter(val byte) {}

// $4012
func (apu *APU) WriteDMCSampleAddr(val byte) {}

// $4013
func (apu *APU) WriteDMCSampleLength(val byte) {}

// $4015 read
func (apu *APU) ReadStatus() byte {
	return 0
}

// PeekStatus is used for debugging
func (apu *APU) PeekStatus() byte {
	return 0
}

// $4015 write
func (apu *APU) WriteStatus(val byte) {}

// $4017
func (apu *APU) WriteFrameCounter(val byte) {}

// https://www.nesdev.org/wiki/APU_Sweep
type sweepUnit byte

func (s *sweepUnit) Enabled() bool {
	return (byte(*s) & 0x80) == 0x80
}
func (s *sweepUnit) Period() byte {
	return ((byte(*s) >> 4) & 0b111) + 1
}
func (s *sweepUnit) Negate() bool {
	return (byte(*s) & 0x04) == 0x04
}
func (s *sweepUnit) ShiftCount() byte {
	return byte(*s) & 0b111
}

type pulseCtrl byte

type pulse struct {
	dutyCycle          byte // 2bit
	lengthCounterHalt  bool
	constantVolumeFlag bool
	sweep              sweepUnit
	timerPeriod        uint16
	lengthCounterLoad  byte
}
