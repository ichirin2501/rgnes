package ppu

type PPU struct {
	r *ppuRegister

	ScanLine int
	Cycle    int
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
		return p.readStatus()
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

func (p *PPU) readStatus() byte {
	v := p.r.Status & (spriteOverflowMask | sprite0HitMask | vBlankStartedMask)
	p.r.SetVBlankStarted(false)
	p.r.Latch = false
	return v
}

// func (p *PPU) writeScroll(val byte) {
// 	if p.r.Latch {

// 		p.r.Latch = false
// 	} else {

// 		p.r.Latch = true
// 	}
// }

func (p *PPU) Step() {
	if p.ScanLine == 241 && p.Cycle == 1 {
		p.r.SetVBlankStarted(true)
	}

	// Pre-render line
	if p.ScanLine == 261 && p.Cycle == 1 {
		p.r.SetVBlankStarted(false)
		p.r.SetSprite0HitFlag(false)
		p.r.SetSpriteOverflow(false)
	}
}
