package cpu

type Interrupter struct {
	delayNMI bool
	nmi      bool
	irq      bool
}

func (i *Interrupter) SetNMI(v bool) {
	i.nmi = v
}
func (i *Interrupter) SetDelayNMI() {
	i.delayNMI = true
}
func (i *Interrupter) SetIRQ(v bool) {
	i.irq = v
}
