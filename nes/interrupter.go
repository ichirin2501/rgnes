package nes

type (
	IRQInterruptLine int
	IRQSource        int
)

const (
	IRQSourceFrameCounter IRQSource = (1 << iota)
	IRQSourceDMC
)

func (i *IRQInterruptLine) SetHigh(src IRQSource) {
	// For interrupt line, define high = 0
	// High = 0
	*i &= IRQInterruptLine(^src)
}

func (i *IRQInterruptLine) SetLow(src IRQSource) {
	*i |= IRQInterruptLine(src)
}

func (i *IRQInterruptLine) IsHigh() bool {
	return *i == 0
}

func (i *IRQInterruptLine) IsLow() bool {
	return *i != 0
}

type NMIInterruptLine int

func (i *NMIInterruptLine) SetHigh() {
	*i = 0
}

func (i *NMIInterruptLine) SetLow() {
	*i = 1
}

func (i *NMIInterruptLine) IsHigh() bool {
	return *i == 0
}

func (i *NMIInterruptLine) IsLow() bool {
	return *i == 1
}
