package nes

import (
	"image/color"
)

var Palette [64]color.Color

func init() {
	colors := []uint32{
		0x666666, 0x002A88, 0x1412A7, 0x3B00A4, 0x5C007E, 0x6E0040, 0x6C0600, 0x561D00,
		0x333500, 0x0B4800, 0x005200, 0x004F08, 0x00404D, 0x000000, 0x000000, 0x000000,
		0xADADAD, 0x155FD9, 0x4240FF, 0x7527FE, 0xA01ACC, 0xB71E7B, 0xB53120, 0x994E00,
		0x6B6D00, 0x388700, 0x0C9300, 0x008F32, 0x007C8D, 0x000000, 0x000000, 0x000000,
		0xFFFEFF, 0x64B0FF, 0x9290FF, 0xC676FF, 0xF36AFF, 0xFE6ECC, 0xFE8170, 0xEA9E22,
		0xBCBE00, 0x88D800, 0x5CE430, 0x45E082, 0x48CDDE, 0x4F4F4F, 0x000000, 0x000000,
		0xFFFEFF, 0xC0DFFF, 0xD3D2FF, 0xE8C8FF, 0xFBC2FF, 0xFEC4EA, 0xFECCC5, 0xF7D8A5,
		0xE4E594, 0xCFEF96, 0xBDF4AB, 0xB3F3CC, 0xB5EBF2, 0xB8B8B8, 0x000000, 0x000000,
	}
	for i, c := range colors {
		r := byte(c >> 16)
		g := byte(c >> 8)
		b := byte(c)
		Palette[i] = &color.RGBA{r, g, b, 0xFF}
	}
}

// A 6-bit value in the palette memory area corresponds to one of 64 outputs
type paletteRAM [32]byte

func (p *paletteRAM) Read(addr paletteForm) byte {
	// $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C
	if addr == 0x10 || addr == 0x14 || addr == 0x18 || addr == 0x1C {
		return p[addr-0x10]
	} else {
		return p[addr]
	}
}

func (p *paletteRAM) Write(addr paletteForm, val byte) {
	if addr == 0x10 || addr == 0x14 || addr == 0x18 || addr == 0x1C {
		p[addr-0x10] = val
	} else {
		p[addr] = val
	}
}

/*
ref: https://www.nesdev.org/wiki/PPU_palettes#Memory_Map

	43210
	|||||
	|||++ - Pixel value from tile data
	|++-- - Palette number from attribute table or OAM
	+---- - Background/Sprite select
*/
type paletteForm byte

const (
	universalBGColorPalette = paletteForm(0)
)

func (p *paletteForm) Pixel() byte {
	return byte(*p) & 0b11
}
func (p *paletteForm) PaletteNumber() byte {
	return (byte(*p) >> 2) & 0b11
}
func (p *paletteForm) IsSprite() bool {
	return (byte(*p) & 0b10000) == 0b10000
}

func newPaletteForm(isSprite bool, paletteNumber byte, pixel byte) paletteForm {
	if paletteNumber > 0b11 {
		panic("palette number must be within 2 bits")
	}
	if pixel > 0b11 {
		panic("pixel value must be within 2 bits")
	}
	if isSprite {
		return paletteForm(0b10000 | (paletteNumber << 2) | pixel)
	} else {
		return paletteForm((paletteNumber << 2) | pixel)
	}
}
