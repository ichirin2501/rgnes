package nes

type interruptLineStatus int

const (
	interruptLineHigh interruptLineStatus = iota
	interruptLineLow
)

type interruptLines struct {
	nmiLine interruptLineStatus
	irqLine interruptLineStatus
}

func (i *interruptLines) setIRQLine(v interruptLineStatus) {
	i.irqLine = v
}

func (i *interruptLines) setNMILine(v interruptLineStatus) {
	i.nmiLine = v
}
