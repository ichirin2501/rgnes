package nes

/*
ref: https://www.nesdev.org/wiki/PPU_registers#PPUCTRL

	7  bit  0
	---- ----
	VPHB SINN
	|||| ||||
	|||| ||++ - Base nametable address
	|||| ||     (0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00)
	|||| |+-- - VRAM address increment per CPU read/write of PPUDATA
	|||| |      (0: add 1, going across; 1: add 32, going down)
	|||| +--- - Sprite pattern table address for 8x8 sprites
	||||        (0: $0000; 1: $1000; ignored in 8x16 mode)
	|||+----- - Background pattern table address (0: $0000; 1: $1000)
	||+------ - Sprite size (0: 8x8 pixels; 1: 8x16 pixels)
	|+------- - PPU master/slave select
	|           (0: read backdrop from EXT pins; 1: output color on EXT pins)
	+-------- - Generate an NMI at the start of the
	            vertical blanking interval (0: off; 1: on)
*/
type ppuControlRegister byte

func (c *ppuControlRegister) baseNameTableAddr() uint16 {
	switch byte(*c) & 0x03 {
	case 0:
		return 0x2000
	case 1:
		return 0x2400
	case 2:
		return 0x2800
	case 3:
		return 0x2C00
	}
	panic("uwaaaaaaaaaaa")
}

func (c *ppuControlRegister) incrVRAMAddr() byte {
	if (byte(*c) & 0x04) == 0 {
		return 1
	} else {
		return 32
	}
}

func (c *ppuControlRegister) spritePatternAddr() uint16 {
	if (byte(*c) & 0x08) == 0 {
		return 0
	} else {
		return 0x1000
	}
}

func (c *ppuControlRegister) bgPatternAddr() uint16 {
	if (byte(*c) & 0x10) == 0 {
		return 0
	} else {
		return 0x1000
	}
}

func (c *ppuControlRegister) spriteSize() byte {
	if (byte(*c) & 0x20) == 0 {
		return 8
	} else {
		return 16
	}
}

func (c *ppuControlRegister) masterSlaveSelect() byte {
	if (byte(*c) & 0x40) == 0 {
		return 0
	} else {
		return 1
	}
}

func (c *ppuControlRegister) generateVBlankNMI() bool {
	return (byte(*c) & 0x80) == 0x80
}

/*
ref: https://www.nesdev.org/wiki/PPU_registers#PPUMASK

	7  bit  0
	---- ----
	BGRs bMmG
	|||| ||||
	|||| |||+ - Greyscale (0: normal color, 1: produce a greyscale display)
	|||| ||+- - 1: Show background in leftmost 8 pixels of screen, 0: Hide
	|||| |+-- - 1: Show sprites in leftmost 8 pixels of screen, 0: Hide
	|||| +--- - 1: Show background
	|||+----- - 1: Show sprites
	||+------ - Emphasize red (green on PAL/Dendy)
	|+------- - Emphasize green (red on PAL/Dendy)
	+-------- - Emphasize blue
*/
type ppuMaskRegister byte

func (m *ppuMaskRegister) isGreyscale() bool {
	return (byte(*m) & 0x01) == 0x01
}
func (m *ppuMaskRegister) showBGLeftMost8pxlScreen() bool {
	return (byte(*m) & 0x02) == 0x02
}
func (m *ppuMaskRegister) showSpritesLeftMost8pxlScreen() bool {
	return (byte(*m) & 0x04) == 0x04
}
func (m *ppuMaskRegister) showBG() bool {
	return (byte(*m) & 0x08) == 0x08
}
func (m *ppuMaskRegister) showSprites() bool {
	return (byte(*m) & 0x10) == 0x10
}
func (m *ppuMaskRegister) emphasizeRed() bool {
	return (byte(*m) & 0x20) == 0x20
}
func (m *ppuMaskRegister) emphasizeGreen() bool {
	return (byte(*m) & 0x40) == 0x40
}
func (m *ppuMaskRegister) emphasizeBlue() bool {
	return (byte(*m) & 0x80) == 0x80
}

/*
ref: https://www.nesdev.org/wiki/PPU_registers#PPUSTATUS

	7  bit  0
	---- ----
	VSO. ....
	|||| ||||
	|||+-++++ - Least significant bits previously written into a PPU register
	|||         (due to register not being updated for this address)
	||+------ - Sprite overflow. The intent was for this flag to be set
	||          whenever more than eight sprites appear on a scanline, but a
	||          hardware bug causes the actual behavior to be more complicated
	||          and generate false positives as well as false negatives; see
	||          PPU sprite evaluation. This flag is set during sprite
	||          evaluation and cleared at dot 1 (the second dot) of the
	||          pre-render line.
	|+------- - Sprite 0 Hit.  Set when a nonzero pixel of sprite 0 overlaps
	|           a nonzero background pixel; cleared at dot 1 of the pre-render
	|           line.  Used for raster timing.
	+-------- - Vertical blank has started (0: not in vblank; 1: in vblank).
	            Set at dot 1 of line 241 (the line *after* the post-render line);
	            cleared after reading $2002 and at dot 1 of the pre-render line.
*/
type ppuStatusRegister byte

func (s *ppuStatusRegister) setSpriteOverflow() {
	*s |= 0x20
}

func (s *ppuStatusRegister) clearSpriteOverflow() {
	*s &= 0xDF
}

func (s *ppuStatusRegister) setSprite0Hit() {
	*s |= 0x40
}

func (s *ppuStatusRegister) clearSprite0Hit() {
	*s &= 0xBF
}

func (s *ppuStatusRegister) setVBlankStarted() {
	*s |= 0x80
}

func (s *ppuStatusRegister) clearVBlankStarted() {
	*s &= 0x7F
}

func (s *ppuStatusRegister) vblankStarted() bool {
	return ((*s) & 0x80) == 0x80
}

func (s *ppuStatusRegister) get() byte {
	return byte(*s)
}

type ppuDecayRegister struct {
	val    byte
	clocks [8]int // refreshed clock with a 1 per bit
}

func (d *ppuDecayRegister) get(clock int) byte {
	// ppu_open_bus/readme.txt
	// > The PPU effectively has a "decay register", an 8-bit register. Each bit
	// > can be refreshed with a 0 or 1. If a bit isn't refreshed with a 1 for
	// > about 600 milliseconds, it will decay to 0 (some decay sooner, depending
	// > on the NES and temperature).

	// 600ms / (1000/60)ms = 36 frame
	// 1 frame = about 341*262 = 89342 PPU clocks
	// 36 * 89342 = 3216312 ppu clocks

	res := byte(0)
	for i := 0; i < 8; i++ {
		if clock-d.clocks[i] < 3216312 {
			res |= d.val & (1 << i)
		}
	}
	return res
}

// ppu_open_bus/readme.txt
// > Writing to any PPU register sets the decay register to the value written.
// So, basically, when writing, all bits are refreshed (= mask = 0xFF).
// When reading, please see ppu_open_bus/readme.txt for details.
func (d *ppuDecayRegister) refresh(v byte, mask byte, clock int) {
	for i := 0; i < 8; i++ {
		if ((1 << i) & mask) == (1 << i) {
			if (v & (1 << i)) == (1 << i) {
				d.val |= (1 << i)
				d.clocks[i] = clock
			} else {
				d.val &= ^(1 << i)
			}
		}
	}
}
