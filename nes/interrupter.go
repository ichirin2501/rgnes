package nes

type interrupter struct {
	delayNMI bool
	nmi      bool
	irq      bool
}

func (i *interrupter) SetNMI(v bool) {
	i.nmi = v
}
func (i *interrupter) SetDelayNMI() {
	i.delayNMI = true
}
func (i *interrupter) SetIRQ(v bool) {
	i.irq = v
}
