package nes

/*
Sprite Attribute

	76543210
	||||||||
	||||||++ - Palette (4 to 7) of sprite
	|||+++-- - Unimplemented (read 0)
	||+----- - Priority (0: in front of background; 1: behind background)
	|+------ - Flip sprite horizontally
	+------- - Flip sprite vertically
*/
type spriteAttribute byte

func (s *spriteAttribute) Palette() byte {
	return byte(*s) & 0x3
}
func (s *spriteAttribute) BehindBackground() bool {
	return (byte(*s) & 0x20) == 0x20
}
func (s *spriteAttribute) FlipSpriteHorizontally() bool {
	return (byte(*s) & 0x40) == 0x40
}
func (s *spriteAttribute) FlipSpriteVertically() bool {
	return (byte(*s) & 0x80) == 0x80
}
func (s *spriteAttribute) Byte() byte {
	return byte(*s)
}

type spriteSlot struct {
	attr spriteAttribute
	x    byte
	lo   byte // pattern table low bit
	hi   byte // pattern table high bit
	idx  byte
}

func (s *spriteSlot) InRange(x int) bool {
	return int(s.x) <= x && x < int(s.x)+8
}

func (s *spriteSlot) PixelPalette(x int) paletteForm {
	dx := x - int(s.x)
	if s.attr.FlipSpriteHorizontally() {
		dx = 7 - dx
	}
	hb := (s.hi & (1 << (7 - dx))) >> (7 - dx)
	lb := (s.lo & (1 << (7 - dx))) >> (7 - dx)
	p := (hb << 1) | lb
	return newPaletteForm(true, s.attr.Palette(), p)
}

/*
getSpriteFromOAM returns (y, tile, attr, x)

ref: https://www.nesdev.org/wiki/PPU_OAM

	Byte 0: Y position of top of sprite
	Byte 1: Tile index number
	Byte 2: Attributes
	Byte 3: X position of left side of sprite
*/
func getSpriteFromOAM(oam []byte, idx byte) (byte, byte, spriteAttribute, byte) {
	y := oam[4*idx+0]
	tile := oam[4*idx+1]
	attr := spriteAttribute(oam[4*idx+2])
	x := oam[4*idx+3]
	return y, tile, attr, x
}
