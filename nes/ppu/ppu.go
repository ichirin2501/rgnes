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
	case 0x0002: // $2002: PPUSTATUS
		return ppu.readStatus()
	case 0x0004:
	case 0x0007:
	}
	return 0
}

// TODO
func (ppu *PPU) Write(addr uint16, val byte) {
	switch addr {
	case 0x0000: // $2000: PPUCTRL
		ppu.writeController(val)
	case 0x0001: // $2001: PPUMASK
		ppu.writeMask(val)
	case 0x0003:
	case 0x0004:
	case 0x0005: // $2005: PPUSCROLL
		ppu.writeScroll(val)
	case 0x0006: // $2006: PPUADDR
		ppu.writeAddr(val)
	case 0x0007:
	default:
	}
}

// $2000: PPUCTRL
func (ppu *PPU) writeController(val byte) {
	ppu.r.Controller = val
	// t: ....BA.. ........ = d: ......BA
	ppu.r.t = (ppu.r.t & 0xF3FF) | (uint16(val)&0x03)<<10
}

// $2001: PPUMASK
func (ppu *PPU) writeMask(val byte) {
	ppu.r.Mask = val
}

// $2002: PPUSTATUS
func (ppu *PPU) readStatus() byte {
	v := ppu.r.Status & (spriteOverflowMask | sprite0HitMask | vBlankStartedMask)
	ppu.r.SetVBlankStarted(false)
	ppu.r.w = false
	return v
}

// $2005: PPUSCROLL
func (ppu *PPU) writeScroll(val byte) {
	if !ppu.r.w {
		// first write
		// d refers to the data written to the port, and A through H to individual bits of a value.
		// t: ........ ...HGFED = d: HGFED...
		// x:               CBA = d: .....CBA
		// w:                   = 1
		ppu.r.t = (ppu.r.t & 0xFFE0) | (uint16(val) >> 3) // HGFED
		ppu.r.x = val & 0x07                              // CBA
		ppu.r.w = true
	} else {
		// second write
		// t: .CBA..HG FED..... = d: HGFEDCBA
		// w:                   = 0
		t1 := (ppu.r.t & 0x8FFF) | ((uint16(val) & 0x07) << 12) // CBA
		t2 := (ppu.r.t & 0xFC1F) | ((uint16(val) & 0xF8) << 2)  // HGFED
		ppu.r.t = t1 | t2
		ppu.r.w = false
	}
}

// $2006: PPUADDR
func (ppu *PPU) writeAddr(val byte) {
	if !ppu.r.w {
		// first write
		// t: ..FEDCBA ........ = d: ..FEDCBA
		// t: .X...... ........ = 0
		// w:                   = 1
		ppu.r.t = (ppu.r.t & 0x80FF) | (uint16(val)&0x3F)<<8
		ppu.r.w = true
	} else {
		// second write
		// t: ........ HGFEDCBA = d: HGFEDCBA
		// v                    = t
		// w:                   = 0
		ppu.r.t = (ppu.r.t & 0xFF00) | uint16(val)
		ppu.r.Addr = ppu.r.t
		ppu.r.w = false
	}
}

func (ppu *PPU) incrementX() {
	v := ppu.r.Addr
	if (v & 0x001F) == 31 {
		v &= ^uint16(0x001F)
		v ^= 0x0400
	} else {
		v++
	}
	ppu.r.Addr = v
}

func (ppu *PPU) incrementY() {
	v := ppu.r.Addr
	if (v & 0x7000) != 0x7000 {
		v += 0x1000
	} else {
		v &= ^uint16(0x7000)
		y := (v & 0x03E0) >> 5
		if y == 29 {
			y = 0
			v ^= 0x0800
		} else if y == 31 {
			y = 0
		} else {
			y++
		}
		v = (v & ^uint16(0x03E0)) | (y << 5)
	}
	ppu.r.Addr = v
}

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
