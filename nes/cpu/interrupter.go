package cpu

type InterruptType byte

const (
	InterruptNone InterruptType = iota
	InterruptNMI
	InterruptIRQ
)

type Interrupter struct {
	delayNMI  bool
	interrupt InterruptType
}

func (i *Interrupter) SetNMI(v bool) {
	if v {
		i.interrupt = InterruptNMI
	} else {
		i.interrupt = InterruptNone
	}
}
func (i *Interrupter) SetDelayNMI() {
	i.delayNMI = true
}
func (i *Interrupter) SetIRQ(v bool) {
	if v {
		i.interrupt = InterruptIRQ
	} else {
		i.interrupt = InterruptNone
	}
}
