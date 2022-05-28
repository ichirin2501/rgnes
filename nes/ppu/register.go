package ppu

// ControlRegister
// 7  bit  0
// ---- ----
// VPHB SINN
// |||| ||||
// |||| ||++ - Base nametable address
// |||| ||     (0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00)
// |||| |+-- - VRAM address increment per CPU read/write of PPUDATA
// |||| |      (0: add 1, going across; 1: add 32, going down)
// |||| +--- - Sprite pattern table address for 8x8 sprites
// ||||        (0: $0000; 1: $1000; ignored in 8x16 mode)
// |||+----- - Background pattern table address (0: $0000; 1: $1000)
// ||+------ - Sprite size (0: 8x8 pixels; 1: 8x16 pixels)
// |+------- - PPU master/slave select
// |           (0: read backdrop from EXT pins; 1: output color on EXT pins)
// +-------- - Generate an NMI at the start of the
//             vertical blanking interval (0: off; 1: on)
type ControlRegister byte

func (c *ControlRegister) BaseNameTableAddr() uint16 {
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

func (c *ControlRegister) IncrementalVRAMAddr() byte {
	if (byte(*c) & 0x04) == 0 {
		return 1
	} else {
		return 32
	}
}

func (c *ControlRegister) SpritePatternAddr() uint16 {
	if (byte(*c) & 0x08) == 0 {
		return 0
	} else {
		return 0x1000
	}
}

func (c *ControlRegister) BackgroundPatternAddr() uint16 {
	if (byte(*c) & 0x10) == 0 {
		return 0
	} else {
		return 0x1000
	}
}

func (c *ControlRegister) SpriteSize() byte {
	if (byte(*c) & 0x20) == 0 {
		return 8
	} else {
		return 16
	}
}

func (c *ControlRegister) MasterSlaveSelect() byte {
	if (byte(*c) & 0x40) == 0 {
		return 0
	} else {
		return 1
	}
}

func (c *ControlRegister) GenerateVBlankNMI() bool {
	return (byte(*c) & 0x80) == 0x80
}

// Mask Register
// 7  bit  0
// ---- ----
// BGRs bMmG
// |||| ||||
// |||| |||+ - Greyscale (0: normal color, 1: produce a greyscale display)
// |||| ||+- - 1: Show background in leftmost 8 pixels of screen, 0: Hide
// |||| |+-- - 1: Show sprites in leftmost 8 pixels of screen, 0: Hide
// |||| +--- - 1: Show background
// |||+----- - 1: Show sprites
// ||+------ - Emphasize red (green on PAL/Dendy)
// |+------- - Emphasize green (red on PAL/Dendy)
// +-------- - Emphasize blue
type MaskRegister byte

func (m *MaskRegister) IsGreyscale() bool {
	return (byte(*m) & 0x01) == 0x01
}
func (m *MaskRegister) ShowBackgroundLeftMost8pxlScreen() bool {
	return (byte(*m) & 0x02) == 0x02
}
func (m *MaskRegister) ShowSpritesLeftMost8pxlScreen() bool {
	return (byte(*m) & 0x04) == 0x04
}
func (m *MaskRegister) ShowBackground() bool {
	return (byte(*m) & 0x08) == 0x08
}
func (m *MaskRegister) ShowSprites() bool {
	return (byte(*m) & 0x10) == 0x10
}
func (m *MaskRegister) EmphasizeRed() bool {
	return (byte(*m) & 0x20) == 0x20
}
func (m *MaskRegister) EmphasizeGreen() bool {
	return (byte(*m) & 0x40) == 0x40
}
func (m *MaskRegister) EmphasizeBlue() bool {
	return (byte(*m) & 0x80) == 0x80
}

// Status Register
// 7  bit  0
// ---- ----
// VSO. ....
// |||| ||||
// |||+-++++ - Least significant bits previously written into a PPU register
// |||         (due to register not being updated for this address)
// ||+------ - Sprite overflow. The intent was for this flag to be set
// ||          whenever more than eight sprites appear on a scanline, but a
// ||          hardware bug causes the actual behavior to be more complicated
// ||          and generate false positives as well as false negatives; see
// ||          PPU sprite evaluation. This flag is set during sprite
// ||          evaluation and cleared at dot 1 (the second dot) of the
// ||          pre-render line.
// |+------- - Sprite 0 Hit.  Set when a nonzero pixel of sprite 0 overlaps
// |           a nonzero background pixel; cleared at dot 1 of the pre-render
// |           line.  Used for raster timing.
// +-------- - Vertical blank has started (0: not in vblank; 1: in vblank).
//             Set at dot 1 of line 241 (the line *after* the post-render
//             line); cleared after reading $2002 and at dot 1 of the
//             pre-render line.
type StatusRegister byte

func (s *StatusRegister) SetSpriteOverflow(val bool) {
	if val {
		*s |= 0x20
	} else {
		*s &= 0xDF
	}
}

func (s *StatusRegister) SetSprite0Hit(val bool) {
	if val {
		*s |= 0x40
	} else {
		*s &= 0xBF
	}
}

func (s *StatusRegister) SetVBlankStarted(val bool) {
	if val {
		*s |= 0x80
	} else {
		*s &= 0x7F
	}
}

func (s *StatusRegister) VBlankStarted() bool {
	return ((*s) & 0x80) == 0x80
}

func (s *StatusRegister) Get() byte {
	return byte(*s)
}

type DecayRegister struct {
	val byte
	tc  int // updated time clock
}

func (d *DecayRegister) Get(currClock int) byte {
	// ppu_open_bus/readme.txt
	// > The PPU effectively has a "decay register", an 8-bit register. Each bit
	// > can be refreshed with a 0 or 1. If a bit isn't refreshed with a 1 for
	// > about 600 milliseconds, it will decay to 0 (some decay sooner, depending
	// > on the NES and temperature).

	// 600ms / 16ms = 37.5 frame
	// 1 frame = about 341*262 = 89342 PPU clocks
	// 37.5 * 89342 = 3350325 ppu clocks
	if currClock-d.tc < 3350325 {
		return d.val
	} else {
		return 0
	}
}

func (d *DecayRegister) Set(currClock int, val byte) {
	d.tc = currClock
	d.val = val
}
