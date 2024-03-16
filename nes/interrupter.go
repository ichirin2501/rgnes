package nes

type interruptLine int

const (
	interruptLineHigh interruptLine = iota
	interruptLineLow

	defaultInterruptLineState interruptLine = interruptLineHigh
)

func (i *interruptLine) SetHigh() {
	*i = interruptLineHigh
}

func (i *interruptLine) SetLow() {
	*i = interruptLineLow
}

func (i *interruptLine) IsHigh() bool {
	return *i == interruptLineHigh
}

func (i *interruptLine) IsLow() bool {
	return *i == interruptLineLow
}
