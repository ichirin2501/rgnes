package apu

type APU struct {
}

func NewAPU() *APU {
	return &APU{}
}

// TODO
func (p *APU) Read(addr uint16) byte {
	switch addr {
	case 0x0015:
	}
	return 0
}

// TODO
func (p *APU) Write(addr uint16, val byte) {
	switch addr {
	default:
	}
}
