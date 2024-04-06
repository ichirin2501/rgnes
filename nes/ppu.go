package nes

import (
	"fmt"
	"image/color"
)

const (
	ScreenWidth  = 256
	ScreenHeight = 240
)

type Renderer interface {
	Render(x, y int, c color.Color)
	Refresh()
}

type ppuRAM struct {
	ram       [2048]byte
	mirroring MirroringType
}

func newPPURAM(m MirroringType) *ppuRAM {
	return &ppuRAM{mirroring: m}
}

func (m *ppuRAM) mirrorAddr(addr uint16) uint16 {
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

func (m *ppuRAM) Read(addr uint16) byte {
	return m.ram[m.mirrorAddr(addr)]
}

func (m *ppuRAM) Write(addr uint16, val byte) {
	m.ram[m.mirrorAddr(addr)] = val
}

type ppuBus struct {
	ram    *ppuRAM
	mapper Mapper
	// https://www.nesdev.org/wiki/Open_bus_behavior#PPU_open_bus
	// > The PPU has two data buses: the I/O bus, used to communicate with the CPU, and the video memory bus.
	// This is the video memory bus variable
	openbus uint16
}

func (bus *ppuBus) Read(addr uint16) byte {
	res := byte(0)
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		res = bus.mapper.Read(addr)
	case 0x2000 <= addr && addr <= 0x2FFF:
		res = bus.ram.Read(addr)
	case 0x3000 <= addr && addr <= 0x3FFF:
		// Mirrors of $2000-$2FFF
		// ref: https://www.nesdev.org/wiki/PPU_registers#The_PPUDATA_read_buffer_(post-fetch)
		// > Simultaneously, the PPU also performs a normal read from the PPU memory at the specified address, "underneath" the palette data,
		res = bus.ram.Read(addr - 0x1000)
	default:
		panic(fmt.Sprintf("read ppubus invalid addr = 0x%04x", addr))
	}
	bus.openbus = (addr & 0x3F00) | uint16(res)
	return res
}

func (bus *ppuBus) Write(addr uint16, val byte) {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		bus.mapper.Write(addr, val)
	case 0x2000 <= addr && addr <= 0x2FFF:
		bus.ram.Write(addr, val)
	case 0x3000 <= addr && addr <= 0x3EFF:
		// Mirrors of $2000-$2EFF
		bus.ram.Write(addr-0x1000, val)
	case 0x3F00 <= addr && addr <= 0x3FFF:
		// nothing
		// https://www.nesdev.org/wiki/PPU_pinout
		// > /RD and /WR specify that the PPU is reading from or writing to VRAM.
		// > As an exception, writing to the internal palette range (3F00-3FFF) will not assert /WR.
	default:
		panic(fmt.Sprintf("write ppubus invalid addr = 0x%04x", addr))
	}
	bus.openbus = (addr & 0x3F00) | uint16(val)
}

type PPU struct {
	bus        *ppuBus
	paletteRAM paletteRAM
	ctrl       ppuControlRegister
	mask       ppuMaskRegister
	status     ppuStatusRegister
	oamAddr    byte
	readBuffer byte   // internal read buffer
	v          uint16 // VRAM address (15 bits)
	t          uint16 // Temporary VRAM address (15 bits); can also be thought of as the address of the top left onscreen tile.
	x          byte   // Fine X scroll (3 bits)
	w          bool   // First or second write toggle (1 bit)
	Scanline   int
	Cycle      int
	renderer   Renderer
	oddFrame   byte // even/odd frame flag (1 bit)

	// https://www.nesdev.org/wiki/Open_bus_behavior#PPU_open_bus
	// > The PPU has two data buses: the I/O bus, used to communicate with the CPU, and the video memory bus.
	// This is the I/O bus variable
	iobus ppuDecayRegister

	suppressVBlankFlag bool

	// sprites temp variables
	// https://www.nesdev.org/wiki/PPU_OAM
	// > The OAM (Object Attribute Memory) is internal memory inside the PPU that contains a display list of up to 64 sprites,
	// > where each sprite's information occupies 4 bytes.
	primaryOAM                    [256]byte
	secondaryOAM                  [32]byte
	secondaryOAMToPrimaryOAMIndex [8]byte
	spriteSlots                   [8]spriteSlot
	spriteFounds                  int

	primaryOAMIndex       int
	secondaryOAMIndex     int
	spriteEvaluationState SpriteEvaluationState
	oamMIdx               int
	copySpriteStateCycle  int

	// background temp variables
	nameTableByte           byte
	bgPaletteNumber         byte
	bgPixelColorIndexLSBits byte
	bgPixelColorIndexMSBits byte
	// 16 bits shift register
	// curr = higher bit =  >>15
	bgPixelColorIndexLSBitsSR uint16
	bgPixelColorIndexMSBitsSR uint16
	// 8 bits shift register and latch
	// curr = higher bit =  >>7
	bgPaletteNumberLSBitsSR byte
	bgPaletteNumberMSBitsSR byte
	bgPaletteNumberLSBLatch byte
	bgPaletteNumberMSBLatch byte

	nmiLine *interruptLine

	Clock int
}

func (ppu *PPU) FetchVBlankStarted() bool {
	return ppu.status.VBlankStarted()
}
func (ppu *PPU) FetchV() uint16 {
	return ppu.v
}
func (ppu *PPU) FetchBuffer() byte {
	return ppu.readBuffer
}

// https://www.nesdev.org/wiki/PPU_power_up_state
// > Palette unspecified
// But I will set the initial value expected by blargg_ppu_tests_2005.09.15b/power_up_palette.nes for now
var powerupPaletteRAM = [32]byte{
	0x09, 0x01, 0x00, 0x01, 0x00, 0x02, 0x02, 0x0D,
	0x08, 0x10, 0x08, 0x24, 0x00, 0x00, 0x04, 0x2C,
	0x09, 0x01, 0x34, 0x03, 0x00, 0x04, 0x00, 0x14,
	0x08, 0x3A, 0x00, 0x02, 0x00, 0x20, 0x2C, 0x08,
}

func NewPPU(renderer Renderer, mapper Mapper, mirror MirroringType, nmiLine *interruptLine) *PPU {
	ppu := &PPU{
		bus: &ppuBus{
			ram:    newPPURAM(mirror),
			mapper: mapper,
		},
		Cycle:    -1,
		renderer: renderer,
		nmiLine:  nmiLine,
	}

	_ = copy(ppu.paletteRAM[:], powerupPaletteRAM[:])

	// init
	for i := 0; i < len(ppu.primaryOAM); i++ {
		ppu.primaryOAM[i] = 0xFF
	}
	for i := 0; i < len(ppu.secondaryOAM); i++ {
		ppu.secondaryOAM[i] = 0xFF
	}
	return ppu
}

func (ppu *PPU) readData(addr uint16) (result byte, isPalette bool, busData byte) {
	addr &= 0x3FFF
	busData = ppu.bus.Read(addr)
	if 0x3F00 <= addr && addr <= 0x3FFF {
		// override
		return ppu.paletteRAM.Read(paletteAddr(addr)), true, busData
	} else {
		return busData, false, busData
	}
}

func (ppu *PPU) writeData(addr uint16, val byte) {
	addr &= 0x3FFF
	ppu.bus.Write(addr, val)
	if 0x3F00 <= addr && addr <= 0x3FFF {
		// override
		ppu.paletteRAM.Write(paletteAddr(addr), val)
	}
}

// ReadRegister is called from CPU Memory Mapped I/O
func (ppu *PPU) ReadRegister(addr uint16) byte {
	// https://www.nesdev.org/wiki/PPU_pinout
	// > CPU A2-A0 are tied to the corresponding CPU address pins and select the PPU register (0-7).
	switch addr & 0x07 {
	case 0:
		// $2000 PPUCTRL write only
	case 1:
		// $2001 PPUMASK write only
	case 2:
		// $2002 PPUSTATUS
		// https://www.nesdev.org/wiki/Open_bus_behavior
		// > Reading the PPU's status port loads a value onto bits 7-5 of the bus, leaving the rest unchanged.
		st := ppu.readStatus()
		ppu.iobus.refresh(st, 0xE0, ppu.Clock)
	case 3:
		// $2003 OAMADDR write only
	case 4:
		// $2004 OAMDATA
		d := ppu.readOAMData()
		ppu.iobus.refresh(d, 0xFF, ppu.Clock)
	case 5:
		// $2005 PPUSCROLL write only
	case 6:
		// $2006 PPUADDR write only
	case 7:
		// $2007 PPUDATA
		// ppu_open_bus/readme.txt
		// D = openbus bit
		// DD-- ----   palette
		result, isPalette, busData := ppu.readPPUData()
		if isPalette {
			// ref: https://www.nesdev.org/wiki/PPU_registers#The_PPUDATA_read_buffer_(post-fetch)
			// > The referenced 6-bit palette data is returned immediately instead of going to the internal read buffer,
			// > and hence no priming read is required.
			ppu.iobus.refresh(result, 0x3F, ppu.Clock)
		} else {
			ppu.iobus.refresh(ppu.readBuffer, 0xFF, ppu.Clock)
		}
		// The buffered value corresponds to bus read data, not palette data
		ppu.readBuffer = busData
	default:
		panic("unreachable")
	}

	return ppu.iobus.get(ppu.Clock)
}

func (ppu *PPU) PeekRegister(addr uint16) byte {
	switch addr & 0x07 {
	case 0:
		return ppu.iobus.get(ppu.Clock)
	case 1:
		return ppu.iobus.get(ppu.Clock)
	case 2:
		return (ppu.status.Get() & 0xE0) | (ppu.iobus.get(ppu.Clock) & 0x1F)
	case 3:
		return ppu.iobus.get(ppu.Clock)
	case 4:
		return ppu.primaryOAM[ppu.oamAddr]
	case 5:
		return ppu.iobus.get(ppu.Clock)
	case 6:
		return ppu.iobus.get(ppu.Clock)
	case 7:
		return ppu.peekPPUData()
	default:
		panic("unreachable")
	}
}

// WriteRegister is called from CPU Memory Mapped I/O
func (ppu *PPU) WriteRegister(addr uint16, val byte) {
	ppu.iobus.refresh(val, 0xFF, ppu.Clock)
	switch addr & 0x07 {
	case 0:
		ppu.writeController(val)
	case 1:
		ppu.writeMask(val)
	case 2:
		// $2002 read only
	case 3:
		ppu.writeOAMAddr(val)
	case 4:
		ppu.writeOAMData(val)
	case 5:
		ppu.writeScroll(val)
	case 6:
		ppu.writePPUAddr(val)
	case 7:
		ppu.writePPUData(val)
	default:
		panic("unreachable")
	}
}

// $2000: PPUCTRL
func (ppu *PPU) writeController(val byte) {
	beforeGeneratedVBlankNMI := ppu.ctrl.GenerateVBlankNMI()
	ppu.ctrl = ppuControlRegister(val)

	// https://www.nesdev.org/wiki/PPU_registers#PPUCTRL
	// > If the PPU is currently in vertical blank, and the PPUSTATUS ($2002) vblank flag is still set (1),
	// > changing the NMI flag in bit 7 of $2000 from 0 to 1 will immediately generate an NMI.
	if ppu.ctrl.GenerateVBlankNMI() && ppu.status.VBlankStarted() {
		if !beforeGeneratedVBlankNMI {
			ppu.nmiLine.SetLow()
		}
	} else {
		ppu.nmiLine.SetHigh()
	}

	// t: ...GH.. ........ <- d: ......GH
	// <used elsewhere>    <- d: ABCDEF..
	ppu.t = (ppu.t & 0xF3FF) | (uint16(val)&0x03)<<10
}

// $2001: PPUMASK
func (ppu *PPU) writeMask(val byte) {
	ppu.mask = ppuMaskRegister(val)
}

// $2002: PPUSTATUS
func (ppu *PPU) readStatus() byte {
	// https://www.nesdev.org/wiki/PPU_frame_timing#VBL_Flag_Timing
	// > Reading $2002 within a few PPU clocks of when VBL is set results in special-case behavior.
	// > Reading one PPU clock before reads it as clear and never sets the flag or generates NMI for that frame.
	// > Reading on the same PPU clock or one later reads it as set, clears it, and suppresses the NMI for that frame.
	// > Reading two or more PPU clocks before/after it's set behaves normally (reads flag's value, clears it, and doesn't affect NMI operation).
	var st byte
	if ppu.Scanline == 241 && ppu.Cycle == 0 {
		ppu.status.ClearVBlankStarted()
		st = ppu.status.Get()
		ppu.suppressVBlankFlag = true
	} else {
		st = ppu.status.Get()
		ppu.status.ClearVBlankStarted()
		ppu.nmiLine.SetHigh()
	}

	// w:                  <- 0
	ppu.w = false

	return st
}

// $2003: OAMADDR
func (ppu *PPU) writeOAMAddr(val byte) {
	ppu.oamAddr = val
}

// $2004: OAMDATA read
func (ppu *PPU) readOAMData() byte {
	return ppu.primaryOAM[ppu.oamAddr]
}

// $2004: OAMDATA write
func (ppu *PPU) writeOAMData(val byte) {
	// https://www.nesdev.org/wiki/PPU_registers#OAM_data_($2004)_%3C%3E_read/write
	// > Writes to OAMDATA during rendering (on the pre-render line and the visible lines 0-239, provided either sprite or background rendering is enabled) do not modify values in OAM
	if !(ppu.isRenderLine() && ppu.isRenderingEnabled()) {
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
		// I' dont' understand the meaning of "bumping only the high 6 bits".
		// On another page, OAM is expressed as follows:
		// > OAM[n][m] below refers to the byte at offset 4*n + m within OAM, i.e. OAM byte m (0-3) of sprite n (0-63).
		// So, I interpreted "bumps the [n] value" as an increase in the address equivalent to one sprite.
		// OAMADDR: [ 0x00 ][ 0x01 ][ 0x02 ][ 0x03 ][ 0x04 ][ 0x05 ][ 0x06 ][ 0x07 ]
		// OAMDATA: [Byte 0][Byte 1][Byte 2][Byte 3][Byte 0][Byte 1][Byte 2][Byte 3]
		//  Sprite: [           Sprite 0           ][           Sprite 1           ]
		// Also, the increase in address for one sprite is equivalent to 32 bits on OAMDATA bits.
		// 32 = 0b100000, <- "bumping only the high 6 bits" ?
		ppu.oamAddr += 4
	}
}

// $2005: PPUSCROLL
func (ppu *PPU) writeScroll(val byte) {
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

// $2006: PPUADDR
func (ppu *PPU) writePPUAddr(val byte) {
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
func (ppu *PPU) readPPUData() (result byte, isPalette bool, busData byte) {
	result, isPalette, busData = ppu.readData(ppu.v)

	if ppu.isRenderLine() && ppu.isRenderingEnabled() {
		// https://www.nesdev.org/wiki/PPU_scrolling#$2007_reads_and_writes
		// > During rendering (on the pre-render line and the visible lines 0-239, provided either background or sprite rendering is enabled),
		// > it will update v in an odd way, triggering a coarse X increment and a Y increment simultaneously (with normal wrapping behavior).
		ppu.incrementX()
		ppu.incrementY()
	} else {
		// normal
		ppu.v += uint16(ppu.ctrl.IncrementalVRAMAddr())
	}

	return result, isPalette, busData
}

// peekPPUData is used for debugging
func (ppu *PPU) peekPPUData() byte {
	addr := ppu.v
	addr &= 0x3FFF
	switch {
	case 0x0000 <= addr && addr <= 0x3EFF:
		// include mirrors of $2000-$2EFF
		return ppu.readBuffer
	case 0x3F00 <= addr && addr <= 0x3FFF:
		// ppu_open_bus/readme.txt
		// D = openbus bit
		// DD-- ----   palette
		return (ppu.iobus.get(ppu.Clock) & 0xC0) | (ppu.paletteRAM.Read(paletteAddr(addr)) & 0x3F)
	default:
		panic(fmt.Sprintf("PeekPPUData invalid addr = 0x%04x", addr))
	}
}

// $2007: PPUDATA write
func (ppu *PPU) writePPUData(val byte) {
	ppu.writeData(ppu.v, val)

	if ppu.isRenderLine() && ppu.isRenderingEnabled() {
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

func (ppu *PPU) fetchNT() {
	v := ppu.v
	addr := 0x2000 | (v & 0x0FFF)
	ppu.nameTableByte, _, _ = ppu.readData(addr)
}

func (ppu *PPU) fetchAT() {
	v := ppu.v
	addr := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	b, _, _ := ppu.readData(addr)
	//
	// b
	// 7654 3210
	// |||| ||++ - Color bits 1-0 for top left quadrant of this byte
	// |||| ++-- - Color bits 3-2 for top right quadrant of this byte
	// ||++----- - Color bits 5-4 for bottom left quadrant of this byte
	// ++------- - Color bits 7-6 for bottom right quadrant of this byte

	// coarse X,Y は画面全体から見たTile(8x8 pixel)のindexを表す
	// Goal: coarse X,Y の情報から、マッピングされる属性テーブルの1byte中の2bitを求めること (その2bitはpallete numberに相当する)
	// ここで属性テーブルの1byteの情報は 32x32 pixel(= 4x4 tile) までの範囲の情報となっていることを思い出そう
	// 例えば、coarse X=[0,1,2,3],[4,5,6,7],... という分け方になる
	// そして属性テーブルの1byte内の表現は上記bのことを指す
	// 対象Tile(8x8)が、上左(16x16)、上右(16x16)、下左(16x16)、下右(16x16)のうち、いずれにマッピングされるかを導出する
	// coarse X,Y の値から、4つに面に対応する2bit毎の位置(shift区切り)を算出するときに、bitテクニックを使うと以下のようになる

	shift := ((v >> 4) & 4) | (v & 2)
	ppu.bgPaletteNumber = (b >> shift) & 0x3
}

func (ppu *PPU) fetchBGLSBits() {
	fineY := (ppu.v >> 12) & 7
	addr := ppu.ctrl.BackgroundPatternAddr() | uint16(ppu.nameTableByte)<<4 | fineY
	ppu.bgPixelColorIndexLSBits, _, _ = ppu.readData(addr)
}

func (ppu *PPU) fetchBGMSBits() {
	fineY := (ppu.v >> 12) & 7
	addr := ppu.ctrl.BackgroundPatternAddr() | uint16(ppu.nameTableByte)<<4 | fineY
	ppu.bgPixelColorIndexMSBits, _, _ = ppu.readData(addr + 8)
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
				// As a result of fixing another bug, maybe there is no problem now?
				return 0, 0, false
			}
			if sattr.FlipSpriteVertically() {
				y = 7 - y
			}
			addr := ppu.ctrl.SpritePatternAddr() | (uint16(stile) << 4) | y
			lo, _, _ := ppu.readData(addr)
			hi, _, _ := ppu.readData(addr + 8)

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
			lo, _, _ := ppu.readData(addr)
			hi, _, _ := ppu.readData(addr + 8)

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
		idx:  ppu.secondaryOAMToPrimaryOAMIndex[byte(sidx)],
	}
	ppu.spriteFounds++
}

type SpriteEvaluationState byte

const (
	InitAndScanPrimaryOAMState SpriteEvaluationState = iota
	ScanPrimaryOAMState
	CopySpriteState
	CheckSpriteOverflowState
	DoneSpriteEvaluationState
)

func (ppu *PPU) evalSpriteForNextScanline() {
	switch ppu.spriteEvaluationState {
	case InitAndScanPrimaryOAMState:
		ppu.primaryOAMIndex = 0
		ppu.secondaryOAMIndex = 0
		ppu.oamMIdx = 0
		ppu.copySpriteStateCycle = -1
		ppu.spriteEvaluationState = ScanPrimaryOAMState
		fallthrough
	case ScanPrimaryOAMState:
		if ppu.Cycle%2 == 0 {
			// > On odd cycles, data is read from (primary) OAM
			// > On even cycles, data is written to secondary OAM (unless secondary OAM is full, in which case it will read the value in secondary OAM instead)
			// > 1. Starting at n = 0, read a sprite's Y-coordinate (OAM[n][0], copying it to the next open slot in secondary OAM (unless 8 sprites have been found, in which case the write is ignored).
			// Regardless of whether the OAM[n][0] is included in the y-axis range, if there is a secondaryOAM slot, it will be copied, so the process will perform in even ppu cycles.
			y := ppu.primaryOAM[4*ppu.primaryOAMIndex+0]
			// Does not transition to this ScanPrimaryOAMState when 8 sprites are found
			ppu.secondaryOAM[4*ppu.secondaryOAMIndex+0] = y

			d := ppu.ctrl.SpriteSize()
			if uint(y) <= uint(ppu.Scanline) && uint(ppu.Scanline) < uint(y)+uint(d) {
				ppu.spriteEvaluationState = CopySpriteState
				// > 1a. If Y-coordinate is in range, copy remaining bytes of sprite data (OAM[n][1] thru OAM[n][3]) into secondary OAM.
				// It takes a total of 2 ppu cycles to read from primaryOAM(odd cycles) and write to secondaryOAM(even cycles).
				// In other words, copying one sprite requires 8 ppu cycles. The remaining three copies will be completed after 6 ppu cycles
				ppu.copySpriteStateCycle = ppu.Cycle + 6
			} else {
				ppu.primaryOAMIndex++
				if ppu.primaryOAMIndex == 64 {
					ppu.spriteEvaluationState = DoneSpriteEvaluationState
				} else {
					ppu.spriteEvaluationState = ScanPrimaryOAMState
				}
			}
		}
	case CopySpriteState:
		if ppu.Cycle == ppu.copySpriteStateCycle {
			ppu.secondaryOAM[4*ppu.secondaryOAMIndex+1] = ppu.primaryOAM[4*ppu.primaryOAMIndex+1]
			ppu.secondaryOAM[4*ppu.secondaryOAMIndex+2] = ppu.primaryOAM[4*ppu.primaryOAMIndex+2]
			ppu.secondaryOAM[4*ppu.secondaryOAMIndex+3] = ppu.primaryOAM[4*ppu.primaryOAMIndex+3]
			ppu.secondaryOAMToPrimaryOAMIndex[ppu.secondaryOAMIndex] = byte(ppu.primaryOAMIndex)
			ppu.secondaryOAMIndex++
			ppu.primaryOAMIndex++
			if ppu.primaryOAMIndex == 64 {
				ppu.spriteEvaluationState = DoneSpriteEvaluationState
			} else if ppu.secondaryOAMIndex < 8 {
				ppu.spriteEvaluationState = ScanPrimaryOAMState
			} else if ppu.secondaryOAMIndex == 8 {
				ppu.spriteEvaluationState = CheckSpriteOverflowState
			}
		}
	case CheckSpriteOverflowState:
		if ppu.Cycle%2 == 1 {
			y1 := ppu.primaryOAM[4*ppu.primaryOAMIndex+int(ppu.oamMIdx)]
			d := ppu.ctrl.SpriteSize()
			if uint(y1) <= uint(ppu.Scanline) && uint(ppu.Scanline) < uint(y1)+uint(d) {
				ppu.status.SetSpriteOverflow()
				ppu.spriteEvaluationState = DoneSpriteEvaluationState
			} else {
				ppu.primaryOAMIndex++
				ppu.oamMIdx = (ppu.oamMIdx + 1) % 4
				if ppu.primaryOAMIndex == 64 {
					ppu.spriteEvaluationState = DoneSpriteEvaluationState
				}
			}
		}
	case DoneSpriteEvaluationState:
		// nothing
	}
}

func (ppu *PPU) isRenderLine() bool {
	return ppu.isPreLine() || ppu.isVisibleScanlines()
}

func (ppu *PPU) isVisibleScanlines() bool {
	return ppu.Scanline < 240
}

func (ppu *PPU) isPreLine() bool {
	return ppu.Scanline == 261
}

func (ppu *PPU) isRenderingEnabled() bool {
	return ppu.mask.ShowBackground() || ppu.mask.ShowSprites()
}

func (ppu *PPU) loadNextBackgroundPaletteData() {
	ppu.bgPixelColorIndexMSBitsSR |= uint16(ppu.bgPixelColorIndexMSBits)
	ppu.bgPixelColorIndexLSBitsSR |= uint16(ppu.bgPixelColorIndexLSBits)
	ppu.bgPaletteNumberLSBLatch = ppu.bgPaletteNumber & 0b1
	ppu.bgPaletteNumberMSBLatch = (ppu.bgPaletteNumber >> 1) & 0b1
}

func (ppu *PPU) getCandidateBackgroundPaletteAddr(x int) paletteAddr {
	if !ppu.mask.ShowBackground() {
		return universalBGColor
	}
	if x < 8 && !ppu.mask.ShowBackgroundLeftMost8pxlScreen() {
		return universalBGColor
	}

	talo := byte((ppu.bgPaletteNumberLSBitsSR >> (7 - ppu.x)) & 0b1)
	tahi := byte((ppu.bgPaletteNumberMSBitsSR >> (7 - ppu.x)) & 0b1)

	ttlo := byte((ppu.bgPixelColorIndexLSBitsSR >> (15 - ppu.x)) & 0b1)
	tthi := byte((ppu.bgPixelColorIndexMSBitsSR >> (15 - ppu.x)) & 0b1)

	return newPaletteAddr(false, (tahi<<1)|talo, (tthi<<1)|ttlo)
}

func (ppu *PPU) getCandidateSpritePaletteAddrAndSlotIndex(x int) (paletteAddr, int) {
	if !ppu.mask.ShowSprites() {
		return universalBGColor, -1
	}
	if x < 8 && !ppu.mask.ShowSpritesLeftMost8pxlScreen() {
		return universalBGColor, -1
	}
	for i := 0; i < ppu.spriteFounds; i++ {
		s := ppu.spriteSlots[i]
		if s.InRange(x) {
			p := s.PaletteAddr(x)
			if p.PixelColorIndex() == 0 {
				// transparent pixel
				// https://www.nesdev.org/wiki/PPU_rendering#Preface
				// > The current pixel for each "active" sprite is checked (from highest to lowest priority),
				// > and the first non-transparent pixel moves on to a multiplexer, where it joins the BG pixel.
				continue
			}
			return p, i
		}
	}
	return universalBGColor, -1
}

func (ppu *PPU) multiplexPaletteAddr(x int) paletteAddr {
	bp := ppu.getCandidateBackgroundPaletteAddr(x)
	sp, slotIdx := ppu.getCandidateSpritePaletteAddrAndSlotIndex(x)

	if bp.PixelColorIndex() == 0 && sp.PixelColorIndex() == 0 {
		// 0x3F00, universal background color
		return universalBGColor
	} else if bp.PixelColorIndex() == 0 && sp.PixelColorIndex() != 0 {
		return sp
	} else if bp.PixelColorIndex() != 0 && sp.PixelColorIndex() == 0 {
		return bp
	} else {
		// bp != 0 && sp != 0
		s := ppu.spriteSlots[slotIdx]

		if x < 255 && s.idx == 0 {
			ppu.status.SetSprite0Hit()
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

	var c color.Color
	if ppu.isRenderingEnabled() {
		addr := ppu.multiplexPaletteAddr(x)
		c = Palette[ppu.paletteRAM.Read(addr)%64]
	} else {
		// https://www.nesdev.org/wiki/PPU_rendering#Rendering_disabled
		// > When the PPU isn't rendering, its v register specifies the current VRAM address (and is output on the PPU's address pins).
		// > Whenever the low 14 bits of v point into palette RAM ($3F00-$3FFF), the PPU will continuously draw the color at that address instead of the EXT input,
		// > overriding the backdrop color.
		if (ppu.v & 0x3F00) == 0x3F00 {
			c = Palette[ppu.paletteRAM.Read(paletteAddr(ppu.v))%64]
		} else {
			c = Palette[ppu.paletteRAM.Read(universalBGColor)%64]
		}
	}
	ppu.renderer.Render(x, y, c)
}

// ref: http://wiki.nesdev.com/w/images/4/4f/Ppu.svg
func (ppu *PPU) Step() {
	ppu.Clock++

	if ppu.isRenderingEnabled() {
		if ppu.oddFrame == 1 && ppu.isPreLine() && ppu.Cycle == 339 {
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
			ppu.oddFrame ^= 1
		}
	}

	visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	preFetchCycle := ppu.Cycle >= 321 && ppu.Cycle <= 336
	fetchCycle := preFetchCycle || visibleCycle

	// https://www.nesdev.org/wiki/File:Ntsc_timing.png
	// > The background shift registers shift during each of dots 2...257 and 322...337, inclusive.
	if ppu.isRenderingEnabled() && ppu.isRenderLine() && ((2 <= ppu.Cycle && ppu.Cycle <= 257) || (322 <= ppu.Cycle && ppu.Cycle <= 337)) {
		// shift
		ppu.bgPaletteNumberMSBitsSR = (ppu.bgPaletteNumberMSBitsSR << 1) | ppu.bgPaletteNumberMSBLatch
		ppu.bgPaletteNumberLSBitsSR = (ppu.bgPaletteNumberLSBitsSR << 1) | ppu.bgPaletteNumberLSBLatch
		ppu.bgPixelColorIndexMSBitsSR <<= 1
		ppu.bgPixelColorIndexLSBitsSR <<= 1
	}
	// https://www.nesdev.org/wiki/File:Ntsc_timing.png
	// > the lower 8bits are then reloaded at ticks 9, 17, 25, ..., 257 and ticks 329 and 337
	if ppu.isRenderingEnabled() && ppu.isRenderLine() && ((9 <= ppu.Cycle && ppu.Cycle <= 257 && ppu.Cycle%8 == 1) || (ppu.Cycle == 329 || ppu.Cycle == 337)) {
		ppu.loadNextBackgroundPaletteData()
	}
	// The screen draws regardless of the PPU rendering mode
	if ppu.isVisibleScanlines() && visibleCycle {
		ppu.renderPixel()
	}
	if ppu.isRenderingEnabled() && ppu.isRenderLine() && fetchCycle {
		switch ppu.Cycle % 8 {
		case 2:
			ppu.fetchNT()
		case 4:
			ppu.fetchAT()
		case 6:
			ppu.fetchBGLSBits()
		case 0:
			ppu.fetchBGMSBits()
		}
	}

	// secondary OAM clear
	if ppu.isRenderingEnabled() && 1 <= ppu.Cycle && ppu.Cycle <= 64 && ppu.isVisibleScanlines() {
		if ppu.Cycle%2 == 1 {
			addr := ppu.Cycle / 2
			ppu.secondaryOAM[addr] = 0xFF
		}
	}

	// sprite eval for next Scanline
	if 65 <= ppu.Cycle && ppu.Cycle <= 256 && ppu.isVisibleScanlines() {
		if ppu.Cycle == 65 {
			ppu.spriteEvaluationState = InitAndScanPrimaryOAMState
		}
		if ppu.isRenderingEnabled() {
			ppu.evalSpriteForNextScanline()
		}
	}

	if ppu.Cycle == 257 && ppu.isRenderLine() {
		// https://www.nesdev.org/wiki/PPU_registers#OAM_address_($2003)_%3E_write
		// > Values during rendering
		// > OAMADDR is set to 0 during each of ticks 257-320 (the sprite tile loading interval) of the pre-render and visible Scanlines.
		// init
		if ppu.isRenderingEnabled() {
			ppu.oamAddr = 0
		}
		ppu.spriteFounds = 0
	}

	// sprite fetch
	if ppu.isRenderingEnabled() && 257 <= ppu.Cycle && ppu.Cycle <= 320 && ppu.isVisibleScanlines() {
		switch ppu.Cycle % 8 {
		case 2:
			// garbage NT byte
		case 4:
			// garbage AT byte
		case 6:
			// fetch sprite pattern table low byte
			// this process is included in fetchSpriteForNextScanline
		case 0:
			ppu.fetchSpriteForNextScanline()
		}
	}

	// update vram register
	if ppu.isRenderingEnabled() {
		if ppu.isRenderLine() {
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
		if ppu.isPreLine() && ppu.Cycle >= 280 && ppu.Cycle <= 304 {
			ppu.copyY()
		}
	}

	// vblank
	if ppu.Scanline == 241 && ppu.Cycle == 1 {
		ppu.renderer.Refresh()
		if !ppu.suppressVBlankFlag {
			ppu.status.SetVBlankStarted()
			if ppu.status.VBlankStarted() && ppu.ctrl.GenerateVBlankNMI() {
				ppu.nmiLine.SetLow()
			}
		}
	}

	// Pre-render line
	if ppu.isPreLine() && ppu.Cycle == 1 {
		ppu.status.ClearVBlankStarted()
		ppu.nmiLine.SetHigh()
		ppu.status.ClearSprite0Hit()
		ppu.status.ClearSpriteOverflow()
	}

}
