package nes

import (
	"fmt"
	"image/color"
)

type Renderer interface {
	Render(x, y int, c color.Color)
	Refresh()
}

type vram struct {
	ram       [2048]byte
	mirroring MirroringType
}

func newVRAM(m MirroringType) *vram {
	return &vram{mirroring: m}
}

func (m *vram) mirrorAddr(addr uint16) uint16 {
	if 0x3000 <= addr {
		panic(fmt.Sprintf("unexpected addr 0x%04X in vram.mirrorAddr", addr))
	}
	nameIdx := (addr - 0x2000) / 0x400
	if m.mirroring.IsHorizontal() {
		// [0x2000 .. 0x2400) and [0x2400 .. 0x2800) => the first 1 KiB of VRAM
		// [0x2800 .. 0x2C00) and [0x2C00 .. 0x3000) => the second 1 KiB of VRAM
		switch nameIdx {
		case 0:
			// addr[0x2000,0x2400) => vaddr[0x000,0x400)
			return addr - 0x2000
		case 1, 2:
			// addr[0x2400,0x2800) => vaddr[0x000,0x400)
			// addr[0x2800,0x2C00) => vaddr[0x400,0x800)
			return addr - 0x2400
		case 3:
			// addr[0x2C00,0x3000) => vaddr[0x400,0x800)
			return addr - 0x2800
		default:
			panic(fmt.Sprintf("unexpected addr 0x%04X in vram.mirrorAddr", addr))
		}
	} else if m.mirroring.IsVertical() {
		// [0x2000 .. 0x2400) and [0x2800 .. 0x2C00) => the first 1 KiB of VRAM
		// [0x2400 .. 0x2800) and [0x2C00 .. 0x3000) => the second 1 KiB of VRAM
		switch nameIdx {
		case 0, 1:
			// addr[0x2000,0x2400) => vaddr[0x000,0x400)
			// addr[0x2400,0x2800) => vaddr[0x400,0x800)
			return addr - 0x2000
		case 2, 3:
			// addr[0x2800,0x2C00) => vaddr[0x000,0x400)
			// addr[0x2C00,0x3000) => vaddr[0x400,0x800)
			return addr - 0x2800
		default:
			panic(fmt.Sprintf("unexpected addr 0x%04X in vram.mirrorAddr", addr))
		}
	} else {
		panic(fmt.Sprintf("unimplemented ppu mirroing addr type: %d", m.mirroring))
	}
}

func (m *vram) Read(addr uint16) byte {
	return m.ram[m.mirrorAddr(addr)]
}

func (m *vram) Write(addr uint16, val byte) {
	m.ram[m.mirrorAddr(addr)] = val
}

type PPU struct {
	mapper       Mapper
	vram         *vram // include nametable and attribute
	paletteTable paletteRAM
	ctrl         ppuControlRegister
	mask         ppuMaskRegister
	status       ppuStatusRegister
	oamAddr      byte
	buf          byte   // internal data buffer
	v            uint16 // VRAM address (15 bits)
	t            uint16 // Temporary VRAM address (15 bits); can also be thought of as the address of the top left onscreen tile.
	x            byte   // Fine X scroll (3 bits)
	w            bool   // First or second write toggle (1 bit)
	Scanline     int
	Cycle        int
	renderer     Renderer
	f            byte // even/odd frame flag (1 bit)
	openbus      ppuDecayRegister

	suppressVBlankFlag bool

	// sprites temp variables
	primaryOAM        [256]byte
	secondaryOAM      [32]byte
	secondaryOAMIndex [8]byte
	spriteSlots       [8]spriteSlot
	spriteFounds      int

	// background temp variables
	nameTableByte        byte
	attributeTableByte   byte
	patternTableLowByte  byte
	patternTableHighByte byte
	// shift register
	// curr = higher bit =  >>15
	patternTableLowBit      uint16
	patternTableHighBit     uint16
	patternAttributeLowBit  uint16
	patternAttributeHighBit uint16

	cpu *interrupter

	nmiDelay int
	Clock    int
}

func (ppu *PPU) FetchVBlankStarted() bool {
	return ppu.status.VBlankStarted()
}
func (ppu *PPU) FetchNMIDelay() int {
	return ppu.nmiDelay
}
func (ppu *PPU) FetchV() uint16 {
	return ppu.v
}
func (ppu *PPU) FetchBuffer() byte {
	return ppu.buf
}

func NewPPU(renderer Renderer, mapper Mapper, mirror MirroringType, c *interrupter) *PPU {
	ppu := &PPU{
		vram:     newVRAM(mirror),
		mapper:   mapper,
		Cycle:    -1,
		renderer: renderer,
		cpu:      c,
	}
	// init
	for i := 0; i < len(ppu.primaryOAM); i++ {
		ppu.primaryOAM[i] = 0xFF
	}
	for i := 0; i < len(ppu.secondaryOAM); i++ {
		ppu.secondaryOAM[i] = 0xFF
	}
	return ppu
}

func (ppu *PPU) updateOpenBus(val byte) {
	ppu.openbus.Set(ppu.Clock, val)
}

func (ppu *PPU) getOpenBus() byte {
	return ppu.openbus.Get(ppu.Clock)
}

func (ppu *PPU) readController() byte {
	// note: $2000 write only
	return ppu.getOpenBus()
}

// PeekController is used for debugging
func (ppu *PPU) PeekController() byte {
	return ppu.getOpenBus()
}

// $2000: PPUCTRL
func (ppu *PPU) writeController(val byte) {
	ppu.updateOpenBus(val)
	beforeGeneratedVBlankNMI := ppu.ctrl.GenerateVBlankNMI()
	ppu.ctrl = ppuControlRegister(val)
	if beforeGeneratedVBlankNMI && !ppu.ctrl.GenerateVBlankNMI() {
		if ppu.Scanline == 241 && (ppu.Cycle == 1 || ppu.Cycle == 2) {
			ppu.cpu.SetNMI(false)
			ppu.nmiDelay = 0
		}
	} else if 241 <= ppu.Scanline && ppu.Scanline <= 260 && !beforeGeneratedVBlankNMI && ppu.ctrl.GenerateVBlankNMI() && ppu.status.VBlankStarted() {
		// https://www.nesdev.org/wiki/PPU_registers#Controller_($2000)_%3E_write
		// > If the PPU is currently in vertical blank, and the PPUSTATUS ($2002) vblank flag is still set (1),
		// > changing the NMI flag in bit 7 of $2000 from 0 to 1 will immediately generate an NMI.
		// vblank flagがoffになるのは Scanline=261,cycle=1 のときだけど、(261,0)でもnmiは発生させないようにする for ppu_vbl_nmi/07-nmi_on_timing.nes
		ppu.cpu.SetNMI(true)
		// https://archive.nes.science/nesdev-forums/f3/t10006.xhtml#p111038
		ppu.cpu.SetDelayNMI()
	}
	// t: ...GH.. ........ <- d: ......GH
	// <used elsewhere>    <- d: ABCDEF..
	ppu.t = (ppu.t & 0xF3FF) | (uint16(val)&0x03)<<10
}

func (ppu *PPU) readMask() byte {
	// note: $2001 write only
	return ppu.getOpenBus()
}

// PeekMask is used for debugging
func (ppu *PPU) PeekMask() byte {
	return ppu.getOpenBus()
}

// $2001: PPUMASK
func (ppu *PPU) writeMask(val byte) {
	ppu.updateOpenBus(val)
	ppu.mask = ppuMaskRegister(val)
}

// $2002: PPUSTATUS
func (ppu *PPU) readStatus() byte {
	st := ppu.status.Get()
	ppu.status.SetVBlankStarted(false)

	// https://www.nesdev.org/wiki/NMI#Race_condition
	// https://www.nesdev.org/wiki/PPU_frame_timing#VBL_Flag_Timing
	if ppu.Scanline == 241 {
		if ppu.Cycle == 0 {
			ppu.suppressVBlankFlag = true
		} else if ppu.Cycle == 1 || ppu.Cycle == 2 {
			ppu.cpu.SetNMI(false)
			ppu.nmiDelay = 0
		}
	}

	// w:                  <- 0
	ppu.w = false

	// https://www.nesdev.org/wiki/Open_bus_behavior
	// > Reading the PPU's status port loads a value onto bits 7-5 of the bus, leaving the rest unchanged.
	return st | (ppu.getOpenBus() & 0x1F)
}

// PeekStatus is used for debugging
func (ppu *PPU) PeekStatus() byte {
	return ppu.status.Get() | (ppu.getOpenBus() & 0x1F)
}

func (ppu *PPU) writeStatus(val byte) {
	// note: $2002 read only
	ppu.updateOpenBus(val)
}

func (ppu *PPU) readOAMAddr() byte {
	// note: $2003 write only
	return ppu.getOpenBus()
}

// PeekOAMAddr is used for debugging
func (ppu *PPU) PeekOAMAddr() byte {
	return ppu.getOpenBus()
}

// $2003: OAMADDR
func (ppu *PPU) writeOAMAddr(val byte) {
	ppu.updateOpenBus(val)
	ppu.oamAddr = val
}

// $2004: OAMDATA read
func (ppu *PPU) readOAMData() byte {
	res := ppu.primaryOAM[ppu.oamAddr]
	ppu.updateOpenBus(res)
	return res
}

// PeekOAMData is used for debuggin
func (ppu *PPU) PeekOAMData() byte {
	return ppu.primaryOAM[ppu.oamAddr]
}

// $2004: OAMDATA write
func (ppu *PPU) writeOAMData(val byte) {
	ppu.updateOpenBus(val)

	// https://www.nesdev.org/wiki/PPU_registers#OAM_data_($2004)_%3C%3E_read/write
	// > Writes to OAMDATA during rendering (on the pre-render line and the visible lines 0-239, provided either sprite or background rendering is enabled) do not modify values in OAM
	if !((ppu.Scanline < 240 || ppu.Scanline == 261) && (ppu.mask.ShowBackground() || ppu.mask.ShowSprites())) {
		// https://www.nesdev.org/wiki/PPU_OAM#Byte_2
		// > The three unimplemented bits of each sprite's byte 2 do not exist in the PPU and always read back as 0 on PPU revisions that allow reading PPU OAM through OAMDATA ($2004).
		// > This can be emulated by ANDing byte 2 with $E3 either when writing to or when reading from OAM.
		if (ppu.oamAddr & 0x03) == 0x02 {
			val &= 0xE3
		}
		ppu.primaryOAM[ppu.oamAddr] = val
		ppu.oamAddr++
	} else {
		// https://www.nesdev.org/wiki/PPU_registers#OAM_data_($2004)_%3C%3E_read/write
		// > but do perform a glitchy increment of OAMADDR, bumping only the high 6 bits (i.e., it bumps the [n] value in PPU sprite evaluation
		// https://forums.nesdev.org/viewtopic.php?t=14140
		// 今はeval spriteの処理をまとめてやってるのでこれは影響ないはず
		ppu.oamAddr += 4
	}
}

func (ppu *PPU) readScroll() byte {
	// note: $2005 write only
	return ppu.getOpenBus()
}

// ReadScroll is used for debugging
func (ppu *PPU) PeekScroll() byte {
	return ppu.getOpenBus()
}

// $2005: PPUSCROLL
func (ppu *PPU) writeScroll(val byte) {
	ppu.updateOpenBus(val)
	if !ppu.w {
		// first write
		// t: ....... ...ABCDE <- d: ABCDE...
		// x:              FGH <- d: .....FGH
		// w:                  <- 1
		ppu.t = (ppu.t & 0xFFE0) | (uint16(val) >> 3) // ABCDE
		ppu.x = val & 0x07                            // FGH
		ppu.w = true
	} else {
		// second write
		// t: FGH..AB CDE..... <- d: ABCDEFGH
		// w:                  <- 0
		ppu.t = (ppu.t & 0x8FFF) | ((uint16(val) & 0x07) << 12)
		ppu.t = (ppu.t & 0xFC1F) | ((uint16(val) & 0xF8) << 2)
		ppu.w = false
	}
}

func (ppu *PPU) readPPUAddr() byte {
	return ppu.getOpenBus()
}

// ReadPPUAddr is used for debugging
func (ppu *PPU) PeekPPUAddr() byte {
	return ppu.getOpenBus()
}

// $2006: PPUADDR
func (ppu *PPU) writePPUAddr(val byte) {
	ppu.updateOpenBus(val)
	if !ppu.w {
		// first write
		// t: .CDEFGH ........ <- d: ..CDEFGH
		//        <unused>     <- d: AB......
		// t: Z...... ........ <- 0 (bit Z is cleared)
		// w:                  <- 1
		ppu.t = (ppu.t & 0x80FF) | (uint16(val)&0x3F)<<8
		ppu.w = true
	} else {
		// second write
		// t: ....... ABCDEFGH <- d: ABCDEFGH
		// v: <...all bits...> <- t: <...all bits...>
		// w:                  <- 0
		ppu.t = (ppu.t & 0xFF00) | uint16(val)
		ppu.v = ppu.t
		ppu.w = false
	}
}

// $2007: PPUDATA read
func (ppu *PPU) readPPUData() byte {
	res := ppu._readPPUData(ppu.v)
	ppu.updateOpenBus(res)

	if (ppu.Scanline < 240 || ppu.Scanline == 261) && (ppu.mask.ShowBackground() || ppu.mask.ShowSprites()) {
		// https://www.nesdev.org/wiki/PPU_scrolling#$2007_reads_and_writes
		// > During rendering (on the pre-render line and the visible lines 0-239, provided either background or sprite rendering is enabled),
		// > it will update v in an odd way, triggering a coarse X increment and a Y increment simultaneously (with normal wrapping behavior).
		ppu.incrementX()
		ppu.incrementY()
	} else {
		// normal
		ppu.v += uint16(ppu.ctrl.IncrementalVRAMAddr())
	}
	return res
}

// PeekPPUData is used for debugging
func (ppu *PPU) PeekPPUData() byte {
	addr := ppu.v
	addr &= 0x3FFF
	switch {
	case 0x0000 <= addr && addr <= 0x3EFF:
		// include mirrors of $2000-$2EFF
		return ppu.buf
	case 0x3F00 <= addr && addr <= 0x3FFF:
		// ppu_open_bus/readme.txt
		// D = openbus bit
		// DD-- ----   palette
		return (ppu.getOpenBus() & 0b11000000) | ppu.paletteTable.Read(paletteForm(addr%0x20))
	default:
		panic(fmt.Sprintf("PeekPPUData invalid addr = 0x%04x", addr))
	}
}

func (ppu *PPU) _readPPUData(addr uint16) byte {
	// https://www.nesdev.org/wiki/PPU_scrolling#PPU_internal_registers
	// > Note that while the v register has 15 bits, the PPU memory space is only 14 bits wide.
	// > The highest bit is unused for access through $2007.
	addr &= 0x3FFF
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		res := ppu.buf
		ppu.buf = ppu.mapper.Read(addr)
		return res
	case 0x2000 <= addr && addr <= 0x2FFF:
		res := ppu.buf
		ppu.buf = ppu.vram.Read(addr)
		return res
	case 0x3000 <= addr && addr <= 0x3EFF:
		// Mirrors of $2000-$2EFF
		res := ppu.buf
		ppu.buf = ppu.vram.Read(addr - 0x1000)
		return res
	case 0x3F00 <= addr && addr <= 0x3FFF:
		// ppu_open_bus/readme.txt
		// D = openbus bit
		// DD-- ----   palette

		// note: [0x3F20,0x3FFF] => Mirrors $3F00-$3F1F
		res := (ppu.getOpenBus() & 0b11000000) | ppu.paletteTable.Read(paletteForm(addr%0x20))
		ppu.buf = ppu.vram.Read(addr - 0x1000)
		return res
	default:
		panic(fmt.Sprintf("readPPUData invalid addr = 0x%04x", addr))
	}
}

// $2007: PPUDATA write
func (ppu *PPU) writePPUData(val byte) {
	ppu.updateOpenBus(val)
	ppu._writePPUData(ppu.v, val)

	if (ppu.Scanline < 240 || ppu.Scanline == 261) && (ppu.mask.ShowBackground() || ppu.mask.ShowSprites()) {
		// https://www.nesdev.org/wiki/PPU_scrolling#$2007_reads_and_writes
		// > During rendering (on the pre-render line and the visible lines 0-239, provided either background or sprite rendering is enabled),
		// > it will update v in an odd way, triggering a coarse X increment and a Y increment simultaneously (with normal wrapping behavior).
		ppu.incrementX()
		ppu.incrementY()
	} else {
		// normal
		ppu.v += uint16(ppu.ctrl.IncrementalVRAMAddr())
	}
}

func (ppu *PPU) _writePPUData(addr uint16, val byte) {
	addr &= 0x3FFF
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		ppu.mapper.Write(addr, val)
	case 0x2000 <= addr && addr <= 0x2FFF:
		ppu.vram.Write(addr, val)
	case 0x3000 <= addr && addr <= 0x3EFF:
		// Mirrors of $2000-$2EFF
		ppu.vram.Write(addr-0x1000, val)
	case 0x3F00 <= addr && addr <= 0x3FFF:
		// note: [0x3F20,0x3FFF] => Mirrors $3F00-$3F1F
		ppu.paletteTable.Write(paletteForm(addr%0x20), val)
	default:
		panic("uaaaaaaaaaaaaaaa")
	}
}

// writeOAMDMAByte is used $4014 from cpubus.
// control cpu clock on cpubus side
func (ppu *PPU) writeOAMDMAByte(val byte) {
	ppu.primaryOAM[ppu.oamAddr] = val
	ppu.oamAddr++
}

// // copyX() is `hori(v) = hori(t)` in NTSC PPU Frame Timing
func (ppu *PPU) copyX() {
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	ppu.v = (ppu.v & 0xFBE0) | (ppu.t & 0x041F)
}

// // copyY() is `vert(v) = vert(t)` in NTSC PPU Frame Timing
func (ppu *PPU) copyY() {
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	ppu.v = (ppu.v & 0x841F) | (ppu.t & 0x7BE0)
}

// incrementX() is `incr hori(v)` in NTSC PPU Frame Timing
// Coarse X increment
func (ppu *PPU) incrementX() {
	v := ppu.v
	if (v & 0x001F) == 31 {
		v &= 0xFFE0
		v ^= 0x0400
	} else {
		v++
	}
	ppu.v = v
}

// incrementY() is `incr vert(v)` in NTSC PPU Frame Timing
func (ppu *PPU) incrementY() {
	v := ppu.v
	if (v & 0x7000) != 0x7000 {
		v += 0x1000
	} else {
		v &= 0x8FFF
		y := (v & 0x03E0) >> 5
		if y == 29 {
			y = 0
			v ^= 0x0800
		} else if y == 31 {
			y = 0
		} else {
			y++
		}
		v = (v & 0xFC1F) | (y << 5)
	}
	ppu.v = v
}

func (ppu *PPU) fetchNameTableByte() {
	v := ppu.v
	addr := 0x2000 | (v & 0x0FFF)
	ppu.nameTableByte = ppu.vram.Read(addr)
}

func (ppu *PPU) fetchAttributeTableByte() {
	v := ppu.v
	addr := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	b := ppu.vram.Read(addr)
	//
	// b
	// 7654 3210
	// |||| ||++ - Color bits 3-2 for top left quadrant of this byte
	// |||| ++-- - Color bits 3-2 for top right quadrant of this byte
	// ||++----- - Color bits 3-2 for bottom left quadrant of this byte
	// ++------- - Color bits 3-2 for bottom right quadrant of this byte

	// coarse X,Y は画面全体から見たTile(8x8 pixel)のindexを表す
	// Goal: coarse X,Y の情報から、マッピングされる属性テーブルの1byte中の2bitを求めること (その2bitはpallete numberに相当する)
	// ここで属性テーブルの1byteの情報は 32x32 pixel(= 4x4 tile) までの範囲の情報となっていることを思い出そう
	// 例えば、coarse X=[0,1,2,3],[4,5,6,7],... という分け方になる
	// そして属性テーブルの1byte内の表現は上記bのことを指す
	// 対象Tile(8x8)が、上左(16x16)、上右(16x16)、下左(16x16)、下右(16x16)のうち、いずれにマッピングされるかを導出する
	// coarse X,Y の値から、4つに面に対応する2bit毎の位置(shift区切り)を算出するときに、bitテクニックを使うと以下のようになる
	shift := ((v >> 4) & 4) | (v & 2)
	// Byteという文字がよくない、既に正確な位置は求めてるので、単にattributeTableで良い？...いやぁ、もう少し細かい
	// shift & 0x3 した時点で、意味が異なっている。パレットテーブルの情報では？？？？ => Palette Attribute
	ppu.attributeTableByte = (b >> shift) & 0x3
}

func (ppu *PPU) fetchPatternTableLowByte() {
	fineY := (ppu.v >> 12) & 7
	addr := ppu.ctrl.BackgroundPatternAddr() | uint16(ppu.nameTableByte)<<4 | fineY
	ppu.patternTableLowByte = ppu.mapper.Read(addr)
}

func (ppu *PPU) fetchPatternTableHighByte() {
	fineY := (ppu.v >> 12) & 7
	addr := ppu.ctrl.BackgroundPatternAddr() | uint16(ppu.nameTableByte)<<4 | fineY
	ppu.patternTableHighByte = ppu.mapper.Read(addr + 8)
}

func (ppu *PPU) fetchSpriteForNextScanline() {
	// called cycle: 264, 272, ..., 320
	sidx := (ppu.Cycle - 264) / 8
	sy, stile, sattr, sx := getSpriteFromOAM(ppu.secondaryOAM[:], byte(sidx))
	lo, hi, ok := func() (byte, byte, bool) {
		if ppu.ctrl.SpriteSize() == 8 {
			y := uint16(ppu.Scanline) - uint16(sy)
			if y > 7 {
				// eval時点で範囲内しか見ない&0xFFで初期化されるが、初期化途中でsprite&bgともにdisableされて前回の状態が残ることがある
				// eval時点だけでなくこのfetchタイミングでも範囲内か確認して少なくとも場外のspriteは表示させないようにしておく
				return 0, 0, false
			}
			if sattr.FlipSpriteVertically() {
				y = 7 - y
			}
			addr := ppu.ctrl.SpritePatternAddr() | (uint16(stile) << 4) | y
			lo := ppu.mapper.Read(addr)
			hi := ppu.mapper.Read(addr + 8)
			return lo, hi, true
		} else {
			// 8x16
			// https://www.nesdev.org/wiki/PPU_OAM#Byte_1
			// > For 8x16 sprites, the PPU ignores the pattern table selection and selects a pattern table from bit 0 of this number.
			bankTile := (uint16(stile) & 0b1) * 0x1000
			tileIndex := uint16(stile) & 0b11111110
			y := uint16(ppu.Scanline) - uint16(sy)
			if y > 15 {
				return 0, 0, false
			}
			if sattr.FlipSpriteVertically() {
				y = 15 - y
			}
			if y > 7 {
				tileIndex++
				y -= 8
			}
			addr := bankTile | tileIndex<<4 | y
			lo := ppu.mapper.Read(addr)
			hi := ppu.mapper.Read(addr + 8)
			return lo, hi, true
		}
	}()
	if !ok {
		return
	}
	ppu.spriteSlots[ppu.spriteFounds] = spriteSlot{
		x:    sx,
		attr: sattr,
		lo:   lo,
		hi:   hi,
		idx:  ppu.secondaryOAMIndex[byte(sidx)],
	}
	ppu.spriteFounds++
}

func (ppu *PPU) evalSpriteForNextScanline() {
	if ppu.Scanline >= 240 {
		panic("uaaaaaaaaaaaaxxxxxaaaa")
	}

	sidx := byte(0)
	for i := byte(0); i < 64; i++ {
		y, tile, attr, x := getSpriteFromOAM(ppu.primaryOAM[:], i)
		// in y range
		d := ppu.ctrl.SpriteSize()
		if uint(y) <= uint(ppu.Scanline) && uint(ppu.Scanline) < uint(y)+uint(d) {
			if sidx < 8 {
				ppu.secondaryOAM[4*sidx+0] = y
				ppu.secondaryOAM[4*sidx+1] = tile
				ppu.secondaryOAM[4*sidx+2] = byte(attr)
				ppu.secondaryOAM[4*sidx+3] = x
				ppu.secondaryOAMIndex[sidx] = i
			}
			sidx++
		}
		if sidx > 8 {
			ppu.status.SetSpriteOverflow(true)
			return
		}
	}
}

func (ppu *PPU) visibleFrame() bool {
	return ppu.Scanline == 261 || ppu.Scanline < 240
}

func (ppu *PPU) visibleScanline() bool {
	return ppu.Scanline < 240
}

func (ppu *PPU) loadNextPixelData() {
	ppu.patternTableHighBit |= uint16(ppu.patternTableHighByte)
	ppu.patternTableLowBit |= uint16(ppu.patternTableLowByte)
	if ppu.attributeTableByte&0x2 == 0x2 {
		ppu.patternAttributeHighBit |= 0xFF
	}
	if ppu.attributeTableByte&0x1 == 0x1 {
		ppu.patternAttributeLowBit |= 0xFF
	}
}

func (ppu *PPU) backgroundPixelPaletteAddr() paletteForm {
	if !ppu.mask.ShowBackground() {
		return universalBGColorPalette
	}
	x := ppu.Cycle - 1
	if x < 8 && !ppu.mask.ShowBackgroundLeftMost8pxlScreen() {
		return universalBGColorPalette
	}

	talo := byte((ppu.patternAttributeLowBit >> (15 - ppu.x)) & 0b1)
	tahi := byte((ppu.patternAttributeHighBit >> (15 - ppu.x)) & 0b1)

	ttlo := byte((ppu.patternTableLowBit >> (15 - ppu.x)) & 0b1)
	tthi := byte((ppu.patternTableHighBit >> (15 - ppu.x)) & 0b1)

	return newPaletteForm(false, (tahi<<1)|talo, (tthi<<1)|ttlo)
}

func (ppu *PPU) spritePixelPaletteAddrAndSlotIndex() (paletteForm, int) {
	if !ppu.mask.ShowSprites() {
		return universalBGColorPalette, -1
	}
	x := ppu.Cycle - 1
	if x < 8 && !ppu.mask.ShowSpritesLeftMost8pxlScreen() {
		return universalBGColorPalette, -1
	}
	for i := 0; i < ppu.spriteFounds; i++ {
		s := ppu.spriteSlots[i]
		if s.InRange(x) {
			p := s.PixelPalette(x)
			if p.Pixel() == 0 {
				// transparent pixel
				// https://www.nesdev.org/wiki/PPU_rendering#Preface
				// > The current pixel for each "active" sprite is checked (from highest to lowest priority),
				// > and the first non-transparent pixel moves on to a multiplexer, where it joins the BG pixel.
				continue
			}
			return p, i
		}
	}
	return universalBGColorPalette, -1
}

func (ppu *PPU) multiplexPixelPaletteAddr() paletteForm {
	bp := ppu.backgroundPixelPaletteAddr()
	sp, slotIdx := ppu.spritePixelPaletteAddrAndSlotIndex()

	if bp.Pixel() == 0 && sp.Pixel() == 0 {
		// 0x3F00, universal background color
		return universalBGColorPalette
	} else if bp.Pixel() == 0 && sp.Pixel() != 0 {
		return sp
	} else if bp.Pixel() != 0 && sp.Pixel() == 0 {
		return bp
	} else {
		// bp != 0 && sp != 0
		s := ppu.spriteSlots[slotIdx]
		x := ppu.Cycle - 1

		if x < 255 && s.idx == 0 {
			ppu.status.SetSprite0Hit(true)
		}
		if s.attr.BehindBackground() {
			return bp
		} else {
			return sp
		}
	}
}

func (ppu *PPU) renderPixel() {
	x := ppu.Cycle - 1 // visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	y := ppu.Scanline

	addr := ppu.multiplexPixelPaletteAddr()
	c := Palette[ppu.paletteTable.Read(addr)%64]
	ppu.renderer.Render(x, y, c)
}

// ref: http://wiki.nesdev.com/w/images/4/4f/Ppu.svg
func (ppu *PPU) Step() {
	ppu.Clock++

	if ppu.nmiDelay > 0 {
		ppu.nmiDelay--
		if ppu.nmiDelay == 0 && ppu.status.VBlankStarted() && ppu.ctrl.GenerateVBlankNMI() {
			ppu.cpu.SetNMI(true)
		}
	}

	if ppu.mask.ShowBackground() || ppu.mask.ShowSprites() {
		if ppu.f == 1 && ppu.Scanline == 261 && ppu.Cycle == 339 {
			// skip 1 cycle
			ppu.Cycle = 340
		}
	}

	ppu.Cycle++
	if ppu.Cycle > 340 {
		ppu.Cycle = 0
		ppu.Scanline++
		if ppu.Scanline > 261 {
			ppu.Scanline = 0
			ppu.suppressVBlankFlag = false
			ppu.f ^= 1
		}
	}

	rendering := ppu.mask.ShowBackground() || ppu.mask.ShowSprites()
	preLine := ppu.Scanline == 261
	renderLine := preLine || ppu.visibleScanline()
	visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	preFetchCycle := ppu.Cycle >= 321 && ppu.Cycle <= 336
	fetchCycle := preFetchCycle || visibleCycle

	if rendering {
		// https://www.nesdev.org/wiki/File:Ntsc_timing.png
		// > The background shift registers shift during each of dots 2...257 and 322...337, inclusive.
		if renderLine && ((2 <= ppu.Cycle && ppu.Cycle <= 257) || (322 <= ppu.Cycle && ppu.Cycle <= 337)) {
			// shift
			ppu.patternAttributeHighBit <<= 1
			ppu.patternAttributeLowBit <<= 1
			ppu.patternTableHighBit <<= 1
			ppu.patternTableLowBit <<= 1
		}
		// https://www.nesdev.org/wiki/File:Ntsc_timing.png
		// > the lower 8bits are then reloaded at ticks 9, 17, 25, ..., 257 and ticks 329 and 337
		if renderLine && ((9 <= ppu.Cycle && ppu.Cycle <= 257 && ppu.Cycle%8 == 1) || (ppu.Cycle == 329 || ppu.Cycle == 337)) {
			ppu.loadNextPixelData()
		}
		if ppu.visibleScanline() && visibleCycle {
			ppu.renderPixel()
		}
		if renderLine && fetchCycle {
			switch ppu.Cycle % 8 {
			case 1:
				ppu.fetchNameTableByte()
			case 3:
				ppu.fetchAttributeTableByte()
			case 5:
				ppu.fetchPatternTableLowByte()
			case 7:
				ppu.fetchPatternTableHighByte()
			}
		}

		// secondary OAM clear
		if 1 <= ppu.Cycle && ppu.Cycle <= 64 && ppu.visibleScanline() {
			if ppu.Cycle%2 == 1 {
				addr := ppu.Cycle / 2
				ppu.secondaryOAM[addr] = 0xFF
			}
		}

		// sprite eval for next Scanline
		// 65 <= ppu.Cycle <= 256
		if ppu.Cycle == 256 && ppu.visibleScanline() {
			ppu.evalSpriteForNextScanline()
		}
		// sprite fetch
		if 257 <= ppu.Cycle && ppu.Cycle <= 320 && (preLine || ppu.visibleScanline()) {
			// https://www.nesdev.org/wiki/PPU_registers#OAM_address_($2003)_%3E_write
			// > Values during rendering
			// > OAMADDR is set to 0 during each of ticks 257-320 (the sprite tile loading interval) of the pre-render and visible Scanlines.
			ppu.oamAddr = 0

			if ppu.Cycle == 257 {
				ppu.spriteFounds = 0
			}

			switch ppu.Cycle % 8 {
			case 1:
				// garbage NT byte
			case 3:
				// garbage AT byte
			case 5:
				// fetch sprite pattern table low byte
				// this process is included in fetchSpriteForNextScanline
			case 7:
				// fetch sprite pattern table high byte
				// this process is included in fetchSpriteForNextScanline
			case 0:
				if ppu.visibleScanline() {
					ppu.fetchSpriteForNextScanline()
				}
			}
		}

		if preLine && ppu.Cycle >= 280 && ppu.Cycle <= 304 {
			ppu.copyY()
		}

		if renderLine {
			if fetchCycle && ppu.Cycle%8 == 0 {
				ppu.incrementX()
			}
			if ppu.Cycle == 256 {
				ppu.incrementY()
			}
			if ppu.Cycle == 257 {
				ppu.copyX()
			}
		}
	}

	// vblank
	if ppu.Scanline == 241 && ppu.Cycle == 1 {
		ppu.renderer.Refresh()
		if !ppu.suppressVBlankFlag {
			ppu.status.SetVBlankStarted(true)
			// hack for vbl_nmi_timing/7.nmi_timing.nes and ppu_vbl_nmi/05-nmi_timing.nes
			ppu.nmiDelay = 2
		}
	}

	// Pre-render line
	if preLine && ppu.Cycle == 1 {
		ppu.status.SetVBlankStarted(false)
		ppu.status.SetSprite0Hit(false)
		ppu.status.SetSpriteOverflow(false)
	}

}
