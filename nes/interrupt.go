package nes

type Interrupt struct {
	nmi bool
	irq bool

	noCopy noCopy
}

func NewInterrupt() *Interrupt {
	return &Interrupt{}
}

func (i *Interrupt) IsNMI() bool {
	return i.nmi
}

func (i *Interrupt) IsIRQ() bool {
	return i.irq
}

func (i *Interrupt) AssertNMI() {
	i.nmi = true
}

func (i *Interrupt) DeassertNMI() {
	i.nmi = false
}

func (i *Interrupt) AssertIRQ() {
	i.irq = true
}

func (i *Interrupt) DeassertIRQ() {
	i.irq = false
}
