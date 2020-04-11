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

	t uint16 // Temporary VRAM address (15 bits); can also be thought of as the address of the top left onscreen tile.
	x byte   // Fine X scroll (3 bits)
	w bool   // First or second write toggle (1 bit)
}

// TODO
func newPPURegister() *ppuRegister {
	return &ppuRegister{}
}

func (r *ppuRegister) NametableID() byte {
	return r.Controller & 0x03
}

func (r *ppuRegister) VRAMAddrIncrFlag() bool {
	return (r.Controller & vRAMAddrIncrMask) == vRAMAddrIncrMask
}

func (r *ppuRegister) SpritePatternAddrFlag() bool {
	return (r.Controller & spritePatternAddrMask) == spritePatternAddrMask
}

func (r *ppuRegister) BackgroundPatternAddrFlag() bool {
	return (r.Controller & backgroundPatternAddrMask) == backgroundPatternAddrMask
}

func (r *ppuRegister) SpriteSizeFlag() bool {
	return (r.Controller & spriteSizeMask) == spriteSizeMask
}

func (r *ppuRegister) NMIOnVBlankFlag() bool {
	return (r.Controller & nmiOnVBlankMask) == nmiOnVBlankMask
}

func (r *ppuRegister) GrayScaleEnable() bool {
	return (r.Mask & grayScaleMask) != grayScaleMask
}

func (r *ppuRegister) ShowLeftBackground() bool {
	return (r.Mask & showLeftBackgroundMask) == showLeftBackgroundMask
}

func (r *ppuRegister) ShowLeftSprites() bool {
	return (r.Mask & showLeftSpritesMask) == showLeftSpritesMask
}

func (r *ppuRegister) ShowBackground() bool {
	return (r.Mask & showBackgroundMask) == showBackgroundMask
}

func (r *ppuRegister) ShowSprites() bool {
	return (r.Mask & showSpritesMask) == showSpritesMask
}

func (r *ppuRegister) SpriteOverflow() bool {
	return (r.Status & spriteOverflowMask) == spriteOverflowMask
}

func (r *ppuRegister) SetSpriteOverflow(cond bool) {
	if cond {
		r.Status |= spriteOverflowMask
	} else {
		r.Status &= ^spriteOverflowMask
	}
}

func (r *ppuRegister) Sprite0HitFlag() bool {
	return (r.Status & sprite0HitMask) == sprite0HitMask
}

func (r *ppuRegister) SetSprite0HitFlag(cond bool) {
	if cond {
		r.Status |= sprite0HitMask
	} else {
		r.Status &= ^sprite0HitMask
	}
}

func (r *ppuRegister) VBlankStarted() bool {
	return (r.Status & vBlankStartedMask) == vBlankStartedMask
}

func (r *ppuRegister) SetVBlankStarted(cond bool) {
	if cond {
		r.Status |= vBlankStartedMask
	} else {
		r.Status &= ^vBlankStartedMask
	}
}
