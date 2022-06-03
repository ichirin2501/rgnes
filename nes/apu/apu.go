package apu

// https://www.nesdev.org/wiki/APU_Pulse
var dutyTable = [][]byte{
	{0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 1, 1},
	{0, 0, 0, 0, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 0},
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
	pulse1 *pulse
	pulse2 *pulse
	tnd    *triangle
	noise  *noise
	dmc    *dmc
}

func New() *APU {
	return &APU{
		pulse1: newPulse(),
		pulse2: newPulse(),
		tnd:    newTriangle(),
		noise:  newNoise(),
		dmc:    newDMC(),
	}
}

func (apu *APU) Step() {
	// todo
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
}

// TTTT TTTT	Timer low (T)
func writePulseTimerLow(p *pulse, val byte) {
	// timer is 11bit
	p.timer.period = (p.timer.period & 0x0700) | uint16(val)
}

// LLLL LTTT	Length counter load (L), timer high (T)
func writePulseLengthAndTimerHigh(p *pulse, val byte) {
	// > If the enabled flag is set, the length counter is loaded with entry L of the length table
	if p.enabled {
		p.lc.load(val >> 3)
	}
	p.timer.period = (p.timer.period & 0x00FF) | (uint16(val&0b111) << 8)

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
// L--- PPPP	Loop noise (L), noise period (P)
func (apu *APU) WriteNoiseLoopAndPeriod(val byte) {
	apu.noise.loop = (val & 0x80) == 0x80
	apu.noise.period = val & 0x0F
}

// $400F
// LLLL L---	Length counter load (L)
func (apu *APU) WriteNoiseLength(val byte) {
	apu.noise.lc.load(val >> 3)
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
	apu.dmc.setEnabled((val & 0x10) == 0x10)
	apu.noise.setEnabled((val & 0x08) == 0x08)
	apu.tnd.setEnabled((val & 0x04) == 0x04)
	apu.pulse2.setEnabled((val & 0x02) == 0x02)
	apu.pulse1.setEnabled((val & 0x01) == 0x01)
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

	lc lengthCounter
	el envelope

	// https://www.nesdev.org/wiki/APU_Sweep
	sweepEnabled    bool
	sweepNegate     bool
	sweepShiftCount byte
	// > Each sweep unit contains the following: divider, reload flag.
	sweepDivider divider
	sweepReload  bool

	duty    byte
	dutyPos byte
	timer   timer
}

func newPulse() *pulse {
	return &pulse{
		timer: newTimer(2),
	}
}

func (p *pulse) setEnabled(v bool) {
	if !v {
		p.lc.value = 0
	}
	p.enabled = v
}

func (p *pulse) output() byte {
	if !p.enabled {
		return 0
	}
	if p.lc.value == 0 {
		return 0
	}
	if dutyTable[p.duty][p.dutyPos] == 0 {
		return 0
	}
	if p.timer.period < 8 {
		return 0
	}
	if p.isMuteSweep() {
		return 0
	}
	return p.el.output()
}

func (p *pulse) targetPeriod() uint16 {
	// https://www.nesdev.org/wiki/APU_Sweep
	// > 1. A barrel shifter shifts the channel's 11-bit raw timer period right by the shift count, producing the change amount.
	// wiki的にはchange amountって言ってるから差分だと思ったんだけど、他実装エミュ見てると、ただのshift結果のコードになってる...
	delta := p.timer.period >> uint16(p.sweepShiftCount)

	// > 2. If the negate flag is true, the change amount is made negative.
	// > 3. The target period is the sum of the current period and the change amount.
	if p.sweepNegate {
		return p.timer.period - delta
	} else {
		return p.timer.period + delta
	}
}

func (p *pulse) isMuteSweep() bool {
	// > Two conditions cause the sweep unit to mute the channel:
	// > 1. If the current period is less than 8, the sweep unit mutes the channel.
	// > 2. If at any time the target period is greater than $7FF, the sweep unit mutes the channel.
	return p.timer.period < 8 || p.targetPeriod() > 0x7FF
}

func (p *pulse) tickSweep() {
	// > 1. If the divider's counter is zero, the sweep is enabled, and the sweep unit is not muting the channel: The pulse's period is adjusted.
	// > 2. If the divider's counter is zero or the reload flag is true: The counter is set to P and the reload flag is cleared. Otherwise, the counter is decremented.

	if p.sweepDivider.tick() && p.sweepEnabled && !p.isMuteSweep() {
		p.timer.period = p.targetPeriod()
	}

	if p.sweepReload {
		p.sweepReload = false
		p.sweepDivider.reload()
	}
}

func (p *pulse) tickTimer() {
	if p.timer.tick() {
		p.dutyPos = (p.dutyPos - 1) & 7
	}
}

func (p *pulse) tickEnvelope() {
	p.el.tick()
}

type triangle struct {
	enabled             bool
	seqPos              byte
	lc                  lengthCounter
	linearCounterCtrl   bool
	linearCounter       byte
	linearCounterPeriod byte
	linearCounterReload bool

	timer timer
}

func newTriangle() *triangle {
	return &triangle{
		timer: newTimer(1),
	}
}

func (t *triangle) setEnabled(v bool) {
	if !v {
		t.lc.value = 0
	}
	t.enabled = v
}

func (t *triangle) output() byte {
	if !t.enabled {
		return 0
	}
	if t.lc.value == 0 {
		return 0
	}
	// todo
	return 0
}

func (t *triangle) tickTimer() {
	if t.timer.tick() {
		if t.lc.value > 0 && t.linearCounter > 0 {
			t.seqPos = (t.seqPos + 1) % 32
		}
	}
}

func (t *triangle) tickLinearCounter() {
	// > 1. If the linear counter reload flag is set, the linear counter is reloaded with the counter reload value, otherwise if the linear counter is non-zero, it is decremented.
	// > 2. If the control flag is clear, the linear counter reload flag is cleared.
	if t.linearCounterReload {
		t.linearCounter = t.linearCounterPeriod
	} else if t.linearCounter > 0 {
		t.linearCounter--
	}
	if !t.linearCounterCtrl {
		t.linearCounterReload = false
	}
}

type noise struct {
	enabled bool
	lc      lengthCounter
	el      envelope
	loop    bool
	period  byte
	timer   timer
}

func newNoise() *noise {
	return &noise{
		timer: newTimer(2),
	}
}

func (n *noise) setEnabled(v bool) {
	if !v {
		n.lc.value = 0
	}
	n.enabled = v
}

func (n *noise) output() byte {
	if !n.enabled {
		return 0
	}
	if n.lc.value == 0 {
		return 0
	}
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

func newDMC() *dmc {
	return &dmc{}
}

func (d *dmc) setEnabled(v bool) {
	if !v {
		// todo
	}
	d.enabled = v
}

func (d *dmc) output() byte {
	// todo
	return 0
}

type envelope struct {
	constantVolume    bool
	loop              bool
	start             bool
	divider           divider
	decayLevelCounter byte
}

func (e *envelope) output() byte {
	// https://www.nesdev.org/wiki/APU_Envelope
	// > The envelope unit's volume output depends on the constant volume flag:
	// > if set, the envelope parameter directly sets the volume, otherwise the decay level is the current volume.
	if e.constantVolume {
		// > bits 3-0	---- VVVV	Used as the volume in constant volume (C set) mode.
		// > Also used as the reload value for the envelope's divider (the period becomes V + 1 quarter frames).
		return byte(e.divider.period)
	} else {
		return e.decayLevelCounter
	}
}

func (e *envelope) tick() {
	// https://www.nesdev.org/wiki/APU_Envelope
	// > When clocked by the frame counter, one of two actions occurs:
	// > if the start flag is clear, the divider is clocked,
	// > otherwise the start flag is cleared, the decay level counter is loaded with 15, and the divider's period is immediately reloaded.
	if e.start {
		e.start = false
		e.decayLevelCounter = 15
		e.divider.reload()
	} else if e.divider.tick() {
		// > When the divider is clocked while at 0, it is loaded with V and clocks the decay level counter.
		// > Then one of two actions occurs: If the counter is non-zero, it is decremented, otherwise if the loop flag is set, the decay level counter is loaded with 15.
		if e.decayLevelCounter > 0 {
			e.decayLevelCounter--
		} else if e.loop {
			e.decayLevelCounter = 15
		}
	}
}

type lengthCounter struct {
	halt  bool
	value byte
}

func (lc *lengthCounter) load(v byte) {
	lc.value = lengthTable[v]
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
	d.counter--
	if d.counter == 0 {
		d.reload()
		return true
	}
	return false
}

func (d *divider) reload() {
	d.counter = d.period * (d.factor + 1)
}
