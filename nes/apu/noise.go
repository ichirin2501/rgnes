package apu

// https://www.nesdev.org/wiki/APU_Noise
// NTSC
var noisePeriodTable = []uint16{
	4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068,
}

type noise struct {
	lc            lengthCounter
	el            envelope
	mode          bool
	shiftRegister uint16
	timer         timer
}

func newNoise() *noise {
	return &noise{}
}

func (n *noise) output() byte {
	if n.lc.value == 0 {
		return 0
	}
	if (n.shiftRegister & 1) == 1 {
		return 0
	}
	return n.el.output()
}

func (n *noise) loadPeriod(p byte) {
	n.timer.period = noisePeriodTable[p]
}

func (n *noise) tickTimer() {
	if n.timer.tick() {
		fb := uint16(0)
		if n.mode {
			fb = (n.shiftRegister & 0x01) ^ ((n.shiftRegister >> 6) & 0x01)
		} else {
			fb = (n.shiftRegister & 0x01) ^ ((n.shiftRegister >> 1) & 0x01)
		}
		n.shiftRegister >>= 1
		n.shiftRegister |= (fb << 14)
	}
}

func (n *noise) tickEnvelope() {
	n.el.tick()
}

func (n *noise) tickLengthCounter() {
	n.lc.tick()
}
