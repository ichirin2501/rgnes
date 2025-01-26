package nes

type (
	irqInterruptLine int
	irqSource        int
)

const (
	irqSourceFrameCounter irqSource = (1 << iota)
	irqSourceDMC
)

func (i *irqInterruptLine) setHigh(src irqSource) {
	// For interrupt line, define high = 0
	// High = 0
	*i &= irqInterruptLine(^src)
}

func (i *irqInterruptLine) setLow(src irqSource) {
	*i |= irqInterruptLine(src)
}

func (i *irqInterruptLine) isHigh() bool {
	return *i == 0
}

func (i *irqInterruptLine) isLow() bool {
	return *i != 0
}

type nmiInterruptLine int

func (i *nmiInterruptLine) setHigh() {
	*i = 0
}

func (i *nmiInterruptLine) setLow() {
	*i = 1
}

func (i *nmiInterruptLine) isHigh() bool {
	return *i == 0
}

func (i *nmiInterruptLine) isLow() bool {
	return *i == 1
}
