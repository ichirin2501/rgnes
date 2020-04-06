package ppu

import "github.com/ichirin2501/rgnes/nes/memory"

type PPU struct {
	r        *ppuRegister
	m        memory.Memory
	ScanLine int
	Cycle    int
}

func NewPPU(m memory.Memory) *PPU {
	return &PPU{
		r: newPPURegister(),
		m: m,
	}
}

// TODO
func (ppu *PPU) Read(addr uint16) byte {
	switch addr {
	case 0x0002:
		return ppu.readStatus()
	case 0x0004:
	case 0x0007:
	}
	return 0
}

// TODO
func (ppu *PPU) Write(addr uint16, val byte) {
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

func (ppu *PPU) readStatus() byte {
	v := ppu.r.Status & (spriteOverflowMask | sprite0HitMask | vBlankStartedMask)
	ppu.r.SetVBlankStarted(false)
	ppu.r.Latch = false
	return v
}

// func (p *PPU) writeScroll(val byte) {
// 	if p.r.Latch {

// 		p.r.Latch = false
// 	} else {

// 		p.r.Latch = true
// 	}
// }

func (ppu *PPU) rendering() bool {
	return ppu.r.ShowBackground() || ppu.r.ShowSprites()
}

func (ppu *PPU) visibleFrame() bool {
	return ppu.ScanLine == 261 || 0 <= ppu.ScanLine && ppu.ScanLine < 240
}

// ref: http://wiki.nesdev.com/w/images/4/4f/Ppu.svg
func (ppu *PPU) Step() {

	if ppu.rendering() {
		// vram fetch
		if ppu.visibleFrame() && (1 <= ppu.Cycle && ppu.Cycle < 257 || 321 <= ppu.Cycle && ppu.Cycle < 337) {
			if ppu.Cycle > 0 {
				switch ppu.Cycle % 8 {
				case 0: // High BG tile byte
				case 2: // NT byte
				case 4: // AT byte
				case 6: // Low BG tile byte
				}
			}
		}

		if ppu.visibleFrame() && (280 <= ppu.Cycle && ppu.Cycle < 305) {
			// vert(v) = vert(t)
		}

		if ppu.visibleFrame() && (0 < ppu.Cycle && ppu.Cycle < 256 || 321 <= ppu.Cycle && ppu.Cycle < 337) {
			// incr hori(v)
		}

		if ppu.visibleFrame() && ppu.Cycle == 256 {
			// incr vert(v)
		}

		if ppu.visibleFrame() && ppu.Cycle == 257 {
			// hori(v) = hori(t)
		}
	}

	if ppu.ScanLine == 241 && ppu.Cycle == 1 {
		ppu.r.SetVBlankStarted(true)
	}

	// Pre-render line
	if ppu.ScanLine == 261 && ppu.Cycle == 1 {
		ppu.r.SetVBlankStarted(false)
		ppu.r.SetSprite0HitFlag(false)
		ppu.r.SetSpriteOverflow(false)
	}

	ppu.Cycle++
	if ppu.Cycle == 341 {
		ppu.Cycle = 0
		ppu.ScanLine++
		if ppu.ScanLine == 262 {
			ppu.ScanLine = 0
			//ppu.Flame ^= 1
		}
	}
}
