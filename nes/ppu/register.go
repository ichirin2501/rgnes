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

	Latch bool
}

// TODO
func newPPURegister() *ppuRegister {
	return &ppuRegister{}
}

func (p *ppuRegister) NametableID() byte {
	return p.Controller & 0x03
}

func (p *ppuRegister) VRAMAddrIncrFlag() bool {
	return (p.Controller & vRAMAddrIncrMask) == vRAMAddrIncrMask
}

func (p *ppuRegister) SpritePatternAddrFlag() bool {
	return (p.Controller & spritePatternAddrMask) == spritePatternAddrMask
}

func (p *ppuRegister) BackgroundPatternAddrFlag() bool {
	return (p.Controller & backgroundPatternAddrMask) == backgroundPatternAddrMask
}

func (p *ppuRegister) SpriteSizeFlag() bool {
	return (p.Controller & spriteSizeMask) == spriteSizeMask
}

func (p *ppuRegister) NMIOnVBlankFlag() bool {
	return (p.Controller & nmiOnVBlankMask) == nmiOnVBlankMask
}

func (p *ppuRegister) GrayScaleEnable() bool {
	return (p.Mask & grayScaleMask) != grayScaleMask
}

func (p *ppuRegister) ShowLeftBackground() bool {
	return (p.Mask & showLeftBackgroundMask) == showLeftBackgroundMask
}

func (p *ppuRegister) ShowLeftSprites() bool {
	return (p.Mask & showLeftSpritesMask) == showLeftSpritesMask
}

func (p *ppuRegister) ShowBackground() bool {
	return (p.Mask & showBackgroundMask) == showBackgroundMask
}

func (p *ppuRegister) ShowSprites() bool {
	return (p.Mask & showSpritesMask) == showSpritesMask
}

func (p *ppuRegister) SpriteOverflow() bool {
	return (p.Status & spriteOverflowMask) == spriteOverflowMask
}

func (p *ppuRegister) SetSpriteOverflow(cond bool) {
	if cond {
		p.Status |= spriteOverflowMask
	} else {
		p.Status &= ^spriteOverflowMask
	}
}

func (p *ppuRegister) Sprite0HitFlag() bool {
	return (p.Status & sprite0HitMask) == sprite0HitMask
}

func (p *ppuRegister) SetSprite0HitFlag(cond bool) {
	if cond {
		p.Status |= sprite0HitMask
	} else {
		p.Status &= ^sprite0HitMask
	}
}

func (p *ppuRegister) VBlankStarted() bool {
	return (p.Status & vBlankStartedMask) == vBlankStartedMask
}

func (p *ppuRegister) SetVBlankStarted(cond bool) {
	if cond {
		p.Status |= vBlankStartedMask
	} else {
		p.Status &= ^vBlankStartedMask
	}
}
