package apu

// https://www.nesdev.org/wiki/APU_Pulse
var dutyTable = [][]byte{
	{0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 1, 1},
	{0, 0, 0, 0, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 0},
}

type pulse struct {
	channel byte
	lc      lengthCounter
	el      envelope

	targetPeriod uint16
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

func newPulse(channel byte) *pulse {
	return &pulse{
		channel: channel,
		timer:   newTimer(2),
	}
}

func (p *pulse) output() byte {
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

// > Whenever the current period changes for any reason, whether by $400x writes or by sweep, the target period also changes.
func (p *pulse) updateTargetPeriod() {
	// https://www.nesdev.org/wiki/APU_Sweep
	// > 1. A barrel shifter shifts the channel's 11-bit raw timer period right by the shift count, producing the change amount.
	// wiki的にはchange amountって言ってるから差分だと思ったんだけど、他実装エミュ見てると、ただのshift結果のコードになってる...
	delta := p.timer.period >> uint16(p.sweepShiftCount)

	// > 2. If the negate flag is true, the change amount is made negative.
	// > 3. The target period is the sum of the current period and the change amount.
	if p.sweepNegate {
		// > The two pulse channels have their adders' carry inputs wired differently, which produces different results when each channel's change amount is made negative:
		// > Pulse 1 adds the ones' complement (−c − 1). Making 20 negative produces a change amount of −21.
		// > Pulse 2 adds the two's complement (−c). Making 20 negative produces a change amount of −20.
		if p.channel == 1 {
			p.targetPeriod = p.timer.period - delta - 1
		} else {
			p.targetPeriod = p.timer.period - delta
		}
	} else {
		p.targetPeriod = p.timer.period + delta
	}
}

func (p *pulse) isMuteSweep() bool {
	// > Two conditions cause the sweep unit to mute the channel:
	// > 1. If the current period is less than 8, the sweep unit mutes the channel.
	// > 2. If at any time the target period is greater than $7FF, the sweep unit mutes the channel.
	return p.timer.period < 8 || p.targetPeriod > 0x7FF
}

func (p *pulse) tickSweep() {
	// > 1. If the divider's counter is zero, the sweep is enabled, and the sweep unit is not muting the channel: The pulse's period is adjusted.
	// > 2. If the divider's counter is zero or the reload flag is true: The counter is set to P and the reload flag is cleared. Otherwise, the counter is decremented.
	// > If the shift count is zero, the pulse channel's period is never updated, but muting logic still applies.
	tick := p.sweepDivider.tick()
	if tick && p.sweepEnabled && !p.isMuteSweep() && p.sweepShiftCount > 0 {
		p.timer.period = p.targetPeriod
		p.updateTargetPeriod()
	}
	if tick || p.sweepReload {
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

func (p *pulse) tickLengthCounter() {
	p.lc.tick()
}
