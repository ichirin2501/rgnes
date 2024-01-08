package nes

// https://www.nesdev.org/wiki/APU_DMC
// > These periods are all even numbers because there are 2 CPU cycles in an APU cycle.
// > A rate of 428 means the output level changes every 214 APU cycles.
// NTSC
var dmcPeriodTable = []uint16{
	428, 380, 340, 320, 286, 254, 226, 214, 190, 160, 142, 128, 106, 84, 72, 54,
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
