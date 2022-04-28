package ppu

import (
	"fmt"
	"image/color"

	"github.com/ichirin2501/rgnes/nes/cassette"
)

type Renderer interface {
	Render(x, y int, c color.Color)
}

type Trace interface {
	SetPPUX(uint16)
	SetPPUY(uint16)
	SetPPUVBlankState(bool)
}

type Interrupter interface {
	TriggerNMI()
}

type vram struct {
	ram       [2048]byte
	mirroring cassette.MirroringType
}

func newVRAM(m cassette.MirroringType) *vram {
	return &vram{mirroring: m}
}

func (m *vram) mirrorAddr(addr uint16) uint16 {
	if 0x3000 <= addr {
		panic(fmt.Sprintf("unexpected addr 0x%04X in vram.mirrorAddr", addr))
	}
	nameIdx := (addr - 0x2000) / 0x400
	if m.mirroring == cassette.MirroringHorizontal {
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
	} else if m.mirroring == cassette.MirroringVertical {
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
	mapper       cassette.Mapper
	vram         *vram // include nametable and attribute
	paletteTable [32]byte
	oamData      [256]byte
	ctrl         ControlRegister
	mask         MaskRegister
	status       StatusRegister
	oamAddr      byte
	buf          byte   // internal data buffer
	v            uint16 // VRAM address
	t            uint16 // Temporary VRAM address (15 bits); can also be thought of as the address of the top left onscreen tile.
	x            byte   // Fine X scroll (3 bits)
	w            bool   // First or second write toggle (1 bit)
	scanLine     int
	Cycle        int
	mirroring    cassette.MirroringType
	renderer     Renderer

	//suppressVBlankFlag bool

	nameTableByte        byte
	attributeTableByte   byte
	patternTableLowByte  byte
	patternTableHighByte byte

	f byte // even/odd frame flag (1 bit)

	// shift register
	// curr = higher bit =  >>15
	patternTableLowBit      uint16
	patternTableHighBit     uint16
	patternAttributeLowBit  uint16
	patternAttributeHighBit uint16

	trace       Trace
	interrupter Interrupter
}

func NewPPU(renderer Renderer, mapper cassette.Mapper, mirroring cassette.MirroringType, i Interrupter, trace Trace) *PPU {
	ppu := &PPU{
		vram:        newVRAM(mirroring),
		mapper:      mapper,
		mirroring:   mirroring,
		renderer:    renderer,
		interrupter: i,
		trace:       trace,
	}
	return ppu
}

// $2000: PPUCTRL
func (ppu *PPU) WriteController(val byte) {
	beforeGeneratedVBlankNMI := ppu.ctrl.GenerateVBlankNMI()
	ppu.ctrl = ControlRegister(val)
	// https://www.nesdev.org/wiki/PPU_registers#Controller_($2000)_%3E_write
	// > If the PPU is currently in vertical blank, and the PPUSTATUS ($2002) vblank flag is still set (1),
	// > changing the NMI flag in bit 7 of $2000 from 0 to 1 will immediately generate an NMI.
	if 241 <= ppu.scanLine && ppu.scanLine <= 260 && !beforeGeneratedVBlankNMI {
		ppu.triggerNMI()
	}
	// t: ...GH.. ........ <- d: ......GH
	// <used elsewhere>    <- d: ABCDEF..
	ppu.t = (ppu.t & 0xF3FF) | (uint16(val)&0x03)<<10
}

// $2001: PPUMASK
func (ppu *PPU) WriteMask(val byte) {
	ppu.mask = MaskRegister(val)
}

// $2002: PPUSTATUS
func (ppu *PPU) ReadStatus() byte {
	st := ppu.status.Get()
	ppu.status.SetVBlankStarted(false)

	// TODO: suppress NMI, scanlineとcycleを見て先読みする
	// https://www.nesdev.org/wiki/NMI#Race_condition

	// w:                  <- 0
	ppu.w = false
	return st
}

// $2003: OAMADDR
func (ppu *PPU) WriteOAMAddr(val byte) {
	ppu.oamAddr = val
}

// $2004: OAMDATA read
func (ppu *PPU) ReadOAMData() byte {
	if ppu.visibleScanLine() && 1 <= ppu.Cycle && ppu.Cycle <= 64 {
		// https://www.nesdev.org/wiki/PPU_sprite_evaluation#Details
		// > Cycles 1-64: Secondary OAM (32-byte buffer for current sprites on scanline) is initialized to $FF - attempting to read $2004 will return $FF.
		return 0xFF
	}
	// https://www.nesdev.org/wiki/PPU_OAM#Byte_2
	// > The three unimplemented bits of each sprite's byte 2 do not exist in the PPU and always read back as 0 on PPU revisions that allow reading PPU OAM through OAMDATA ($2004).
	// > This can be emulated by ANDing byte 2 with $E3 either when writing to or when reading from OAM.
	// > Bits that have decayed also read back as 0 through OAMDATA.
	res := ppu.oamData[ppu.oamAddr]
	if (ppu.oamAddr & 0x03) == 0x02 {
		res &= 0xE3
	}
	return res
}

// $2004: OAMDATA write
func (ppu *PPU) WriteOAMData(val byte) {
	ppu.oamData[ppu.oamAddr] = val
	ppu.oamAddr++
}

// $2005: PPUSCROLL
func (ppu *PPU) WriteScroll(val byte) {
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
func (ppu *PPU) WritePPUAddr(val byte) {
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
func (ppu *PPU) ReadPPUData() byte {
	addr := ppu.v
	ppu.v += uint16(ppu.ctrl.IncrementalVRAMAddr())
	return ppu.readPPUData(addr)
}

func (ppu *PPU) readPPUData(addr uint16) byte {
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
		return ppu.readPPUData(addr - 0x1000)
	case 0x3F00 <= addr && addr <= 0x3F1F:
		// https://www.nesdev.org/wiki/PPU_registers#The_PPUDATA_read_buffer_(post-fetch)
		// > Reading palette data from $3F00-$3FFF works differently. The palette data is placed immediately on the data bus, and hence no priming read is required.
		// > Reading the palettes still updates the internal buffer though, but the data placed in it is the mirrored nametable data that would appear "underneath" the palette.
		// > (Checking the PPU memory map should make this clearer.)
		if addr == 0x3F10 || addr == 0x3F14 || addr == 0x3F18 || addr == 0x3F1C {
			// $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C
			res := ppu.paletteTable[addr-0x3F00-0x10]
			ppu.buf = ppu.readPPUData(addr - 0x1000)
			return res
		} else {
			res := ppu.paletteTable[addr-0x3F00]
			ppu.buf = ppu.readPPUData(addr - 0x1000)
			return res
		}
	case 0x3F20 <= addr && addr <= 0x3FFF:
		// Mirrors of $3F00-$3F1F
		return ppu.readPPUData(0x3F00 + ((addr - 0x3F20) % 32))
	default:
		panic(fmt.Sprintf("readPPUData invalid addr = 0x%04x", addr))
	}
}

// $2007: PPUDATA write
func (ppu *PPU) WritePPUData(val byte) {
	addr := ppu.v
	ppu.v += uint16(ppu.ctrl.IncrementalVRAMAddr())
	ppu.writePPUData(addr, val)
}

func (ppu *PPU) writePPUData(addr uint16, val byte) {
	addr &= 0x3FFF
	switch {
	case 0x000 <= addr && addr <= 0x1FFF:
		ppu.mapper.Write(addr, val)
	case 0x2000 <= addr && addr <= 0x2FFF:
		ppu.vram.Write(addr, val)
	case 0x3000 <= addr && addr <= 0x3EFF:
		// Mirrors of $2000-$2EFF
		ppu.writePPUData(addr-0x1000, val)
	case 0x3F00 <= addr && addr <= 0x3F1F:
		if addr == 0x3F10 || addr == 0x3F14 || addr == 0x3F18 || addr == 0x3F1C {
			// $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C
			ppu.paletteTable[addr-0x3F00-0x10] = val
		} else {
			ppu.paletteTable[addr-0x3F00] = val
		}
	case 0x3F20 <= addr && addr <= 0x3FFF:
		// Mirrors of $3F00-$3F1F
		ppu.writePPUData(0x3F00+((addr-0x3F20)%32), val)
	default:
		panic("uaaaaaaaaaaaaaaa")
	}
}

// $4014: OAMDMA write
func (ppu *PPU) WriteOAMDMA(data []byte) {
	for i := 0; i < len(data); i++ {
		ppu.oamData[ppu.oamAddr] = data[i]
		ppu.oamAddr++
	}
	// todo: cycle
}

func (ppu *PPU) triggerNMI() {
	if ppu.status.VBlankStarted() && ppu.ctrl.GenerateVBlankNMI() {
		ppu.interrupter.TriggerNMI()
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
	addr := ppu.ctrl.BackgroundPatternAddr() + uint16(ppu.nameTableByte)*16 + fineY
	ppu.patternTableLowByte = ppu.mapper.Read(addr)
}

func (ppu *PPU) fetchPatternTableHighByte() {
	fineY := (ppu.v >> 12) & 7
	addr := ppu.ctrl.BackgroundPatternAddr() + uint16(ppu.nameTableByte)*16 + fineY
	ppu.patternTableHighByte = ppu.mapper.Read(addr + 8)
}

func (ppu *PPU) visibleFrame() bool {
	return ppu.scanLine == 261 || 0 <= ppu.scanLine && ppu.scanLine < 240
}

func (ppu *PPU) visibleScanLine() bool {
	return ppu.scanLine < 240
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

func (ppu *PPU) backgroundPixel() byte {
	if !ppu.mask.ShowBackground() {
		return 0
	}
	addr := byte(ppu.patternAttributeHighBit>>15)<<3 |
		byte(ppu.patternAttributeLowBit>>15)<<2 |
		byte(ppu.patternTableHighBit>>15)<<1 |
		byte(ppu.patternTableLowBit>>15)
	return addr
}

func (ppu *PPU) spritePixel() byte {
	if !ppu.mask.ShowSprites() {
		return 0
	}
	// TODO
	return 0
}

func (ppu *PPU) renderPixel() {
	x := ppu.Cycle - 1 // visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	y := ppu.scanLine

	addr := byte(0)

	// todo
	bg := ppu.backgroundPixel()
	//sp := ppu.spritePixel()
	addr = bg

	c := Palette[(ppu.paletteTable[addr])%64]
	ppu.renderer.Render(x, y, c)
}

func (ppu *PPU) tick() {
	if ppu.mask.ShowBackground() || ppu.mask.ShowSprites() {
		if ppu.f == 1 && ppu.scanLine == 261 && ppu.Cycle == 339 {
			ppu.Cycle = 0
			ppu.scanLine = 0
			ppu.f ^= 1
			return
		}
	}
	ppu.Cycle++
	if ppu.Cycle > 340 {
		ppu.Cycle = 0
		ppu.scanLine++
		if ppu.scanLine > 261 {
			ppu.scanLine = 0
			//ppu.suppressVBlankFlag = false
			ppu.f ^= 1
		}
	}
}

// ref: http://wiki.nesdev.com/w/images/4/4f/Ppu.svg
func (ppu *PPU) Step() {
	if ppu.trace != nil {
		ppu.trace.SetPPUX(uint16(ppu.Cycle))
		ppu.trace.SetPPUY(uint16(ppu.scanLine))
		ppu.trace.SetPPUVBlankState((ppu.status.Get() & 0x80) == 0x80)
	}
	rendering := ppu.mask.ShowBackground() || ppu.mask.ShowSprites()
	preLine := ppu.scanLine == 261
	renderLine := preLine || ppu.visibleScanLine()
	visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	preFetchCycle := ppu.Cycle >= 321 && ppu.Cycle <= 336
	fetchCycle := preFetchCycle || visibleCycle

	if rendering {
		if ppu.visibleScanLine() && visibleCycle {
			ppu.renderPixel()
		}
		if renderLine && fetchCycle {
			// shift
			ppu.patternAttributeHighBit <<= 1
			ppu.patternAttributeLowBit <<= 1
			ppu.patternTableHighBit <<= 1
			ppu.patternTableLowBit <<= 1
			switch ppu.Cycle % 8 {
			case 1:
				ppu.fetchNameTableByte()
			case 3:
				ppu.fetchAttributeTableByte()
			case 5:
				ppu.fetchPatternTableLowByte()
			case 7:
				ppu.fetchPatternTableHighByte()
			case 0:
				ppu.loadNextPixelData()
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
	if ppu.scanLine == 241 && ppu.Cycle == 1 {
		ppu.status.SetVBlankStarted(true)
		ppu.triggerNMI()
	}

	// Pre-render line
	if preLine && ppu.Cycle == 1 {
		ppu.status.SetVBlankStarted(false)
		ppu.status.SetSprite0Hit(false)
		ppu.status.SetSpriteOverflow(false)
	}

	ppu.tick()
}
