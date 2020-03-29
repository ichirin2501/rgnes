package ppu

const (
	_ byte = (1 << iota)
	_
	vRAMAddrIncrMask
	spritePatternAddrMask
	backgroundPatternAddrMask
	spriteSizeMask
	_
	nmiOnVBlankMask
)

const (
	grayscaleMask byte = (1 << iota)
	showLeftBackgroundMask
	showLeftSpritesMask
	showBackgroundMask
	showSpritesMask
	emphaticRedMask
	emphaticGreenMask
	emphaticBlueMask
)

type ppuRegister struct {
	Controller byte
	Mask       byte
	Status     byte
	OAMAddr    byte
	Scroll     uint16
	Addr       uint16
	Data       byte
}

// TODO
func newPPURegister() *ppuRegister {
	return &ppuRegister{}
}

func (p *ppuRegister) nameTableID() byte {
	return p.Controller & 0x03
}

func (p *ppuRegister) ppuAddrIncrFlag() bool {
	return (p.Controller & vRAMAddrIncrMask) == vRAMAddrIncrMask
}

func (p *ppuRegister) spritePatternAddrFlag() bool {
	return (p.Controller & spritePatternAddrMask) == spritePatternAddrMask
}

func (p *ppuRegister) backgroundPatternAddrFlag() bool {
	return (p.Controller & backgroundPatternAddrMask) == backgroundPatternAddrMask
}

func (p *ppuRegister) spriteSizeFlag() bool {
	return (p.Controller & spriteSizeMask) == spriteSizeMask
}

func (p *ppuRegister) nmiOnVBlankFlag() bool {
	return (p.Controller & nmiOnVBlankMask) == nmiOnVBlankMask
}

func (p *ppuRegister) grayScaleEnable() bool {
	return (p.Mask & grayscaleMask) != grayscaleMask
}

func (p *ppuRegister) showLeftBackground() bool {
	return (p.Mask & showLeftBackgroundMask) == showLeftBackgroundMask
}

func (p *ppuRegister) showLeftSprites() bool {
	return (p.Mask & showLeftSpritesMask) == showLeftSpritesMask
}

func (p *ppuRegister) showBackground() bool {
	return (p.Mask & showBackgroundMask) == showBackgroundMask
}

func (p *ppuRegister) showSprites() bool {
	return (p.Mask & showSpritesMask) == showSpritesMask
}
