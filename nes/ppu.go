package nes

type PPU struct {
	r      *ppuRegister
	noCopy noCopy
}

func NewPPU() *PPU {
	return &PPU{
		r: newPPURegister(),
	}
}

// TODO
func (p *PPU) Read(addr uint16) byte {
	switch addr {
	case 0x0002:
	case 0x0004:
	case 0x0007:
	}
	return 0
}

// TODO
func (p *PPU) Write(addr uint16, val byte) {
	switch addr {
	case 0x0000:
	case 0x0001:
	case 0x0003:
	case 0x0004:
	case 0x0005:
	case 0x0006:
	case 0x0007:
	default:
	}
}
