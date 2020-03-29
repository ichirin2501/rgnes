package ppu

// Controller ($2000)
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

// Mask ($2001)
const (
	grayScaleMask byte = (1 << iota)
	showLeftBackgroundMask
	showLeftSpritesMask
	showBackgroundMask
	showSpritesMask
	emphaticRedMask
	emphaticGreenMask
	emphaticBlueMask
)

// Status ($2002)
const (
	_ byte = 1 << iota
	_
	_
	_
	_
	spriteOverflowMask
	sprite0HitMask
	vBlankStartedMask
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

func (p *ppuRegister) nametableID() byte {
	return p.Controller & 0x03
}

func (p *ppuRegister) vRAMAddrIncrFlag() bool {
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
	return (p.Mask & grayScaleMask) != grayScaleMask
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

func (p *ppuRegister) spriteOverflow() bool {
	return (p.Status & spriteOverflowMask) == spriteOverflowMask
}

func (p *ppuRegister) sprite0HitFlag() bool {
	return (p.Status & sprite0HitMask) == sprite0HitMask
}

func (p *ppuRegister) vBlankStarted() bool {
	return (p.Status & vBlankStartedMask) == vBlankStartedMask
}
