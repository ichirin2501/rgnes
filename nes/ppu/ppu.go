package ppu

import (
	"fmt"
	"image/color"

	"github.com/ichirin2501/rgnes/nes/cassette"
)

type Renderer interface {
	Render(x, y int, c color.Color)
	Refresh()
}

type Trace interface {
	SetPPUX(uint16)
	SetPPUY(uint16)
	SetPPUVBlankState(bool)
}

type CPU interface {
	SetDelayNMI()
	DMASuspend()
	SetNMI(val bool)
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

// 76543210
// ||||||||
// ||||||++ - Palette (4 to 7) of sprite
// |||+++-- - Unimplemented (read 0)
// ||+----- - Priority (0: in front of background; 1: behind background)
// |+------ - Flip sprite horizontally
// +------- - Flip sprite vertically
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

// https://www.nesdev.org/wiki/PPU_OAM
type sprite struct {
	y          byte            // Byte 0
	tileIndex  byte            // Byte 1
	attributes spriteAttribute // Byte 2
	x          byte            // Byte 3
}

type oam []*sprite

func (o oam) GetSpriteByAddr(addr byte) *sprite {
	return o[addr/4]
}
func (o oam) GetSpriteByIndex(idx byte) *sprite {
	return o[idx]
}
func (o oam) SetSpriteByIndex(idx byte, src sprite) {
	o[idx] = &src
}
func (o oam) GetByte(addr byte) byte {
	e := o[addr/4]
	switch addr % 4 {
	case 0:
		return e.y
	case 1:
		return e.tileIndex
	case 2:
		return e.attributes.Byte()
	case 3:
		return e.x
	}
	panic("aaaaaaaaaaaaaaa")
}
func (o oam) SetByte(addr byte, val byte) {
	e := o[addr/4]
	switch addr % 4 {
	case 0:
		e.y = val
	case 1:
		e.tileIndex = val
	case 2:
		e.attributes = spriteAttribute(val)
	case 3:
		e.x = val
	}
}

type spriteSlot struct {
	attributes spriteAttribute
	x          byte
	lo         byte // pattern table low bit
	hi         byte // pattern table high bit
}

type PPU struct {
	mapper       cassette.Mapper
	vram         *vram // include nametable and attribute
	paletteTable paletteRAM
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
	f            byte // even/odd frame flag (1 bit)
	openbus      byte

	suppressVBlankFlag bool

	// sprites temp variables
	primaryOAM   oam
	secondaryOAM oam
	spriteSlots  [8]spriteSlot
	spriteFounds int

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

	trace Trace
	cpu   CPU

	nmiDelay int
	Clock    int
}

func (ppu *PPU) FetchVBlankStarted() bool {
	return ppu.status.VBlankStarted()
}
func (ppu *PPU) FetchNMIDelay() int {
	return ppu.nmiDelay
}
func (ppu *PPU) FetchScanline() int {
	return ppu.scanLine
}

func NewPPU(renderer Renderer, mapper cassette.Mapper, mirroring cassette.MirroringType, c CPU, trace Trace) *PPU {
	po := make([]*sprite, 64)
	for i := 0; i < 64; i++ {
		po[i] = &sprite{0xFF, 0xFF, 0xFF, 0xFF}
	}
	so := make([]*sprite, 8)
	for i := 0; i < 8; i++ {
		so[i] = &sprite{0xFF, 0xFF, 0xFF, 0xFF}
	}
	ppu := &PPU{
		vram:         newVRAM(mirroring),
		mapper:       mapper,
		mirroring:    mirroring,
		primaryOAM:   oam(po),
		secondaryOAM: oam(so),

		Cycle: -1,

		renderer: renderer,
		cpu:      c,
		trace:    trace,
	}
	return ppu
}

func (ppu *PPU) ReadController() byte {
	// note: $2000 write only
	return ppu.openbus
}

// $2000: PPUCTRL
func (ppu *PPU) WriteController(val byte) {
	ppu.openbus = val
	beforeGeneratedVBlankNMI := ppu.ctrl.GenerateVBlankNMI()
	ppu.ctrl = ControlRegister(val)
	if beforeGeneratedVBlankNMI && !ppu.ctrl.GenerateVBlankNMI() {
		if ppu.scanLine == 241 && (ppu.Cycle == 1 || ppu.Cycle == 2) {
			ppu.cpu.SetNMI(false)
			ppu.nmiDelay = 0
		}
	} else if 241 <= ppu.scanLine && ppu.scanLine <= 260 && !beforeGeneratedVBlankNMI && ppu.ctrl.GenerateVBlankNMI() && ppu.status.VBlankStarted() {
		// https://www.nesdev.org/wiki/PPU_registers#Controller_($2000)_%3E_write
		// > If the PPU is currently in vertical blank, and the PPUSTATUS ($2002) vblank flag is still set (1),
		// > changing the NMI flag in bit 7 of $2000 from 0 to 1 will immediately generate an NMI.
		// vblank flagがoffになるのは scanline=261,cycle=1 のときだけど、(261,0)でもnmiは発生させないようにする for ppu_vbl_nmi/07-nmi_on_timing.nes
		ppu.cpu.SetNMI(true)
		// https://archive.nes.science/nesdev-forums/f3/t10006.xhtml#p111038
		ppu.cpu.SetDelayNMI()
	}
	// t: ...GH.. ........ <- d: ......GH
	// <used elsewhere>    <- d: ABCDEF..
	ppu.t = (ppu.t & 0xF3FF) | (uint16(val)&0x03)<<10
}

func (ppu *PPU) ReadMask() byte {
	// note: $2001 write only
	return ppu.openbus
}

// $2001: PPUMASK
func (ppu *PPU) WriteMask(val byte) {
	ppu.openbus = val
	ppu.mask = MaskRegister(val)
}

// $2002: PPUSTATUS
func (ppu *PPU) ReadStatus() byte {
	st := ppu.status.Get()
	ppu.status.SetVBlankStarted(false)

	// https://www.nesdev.org/wiki/NMI#Race_condition
	// https://www.nesdev.org/wiki/PPU_frame_timing#VBL_Flag_Timing
	if ppu.scanLine == 241 {
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
	// ????
	return st | (ppu.openbus & 0x1F)
}
func (ppu *PPU) WriteStatus(val byte) {
	// note: $2002 read only
	ppu.openbus = val
}

func (ppu *PPU) ReadOAMAddr() byte {
	// note: $2003 write only
	return ppu.openbus
}

// $2003: OAMADDR
func (ppu *PPU) WriteOAMAddr(val byte) {
	ppu.openbus = val
	ppu.oamAddr = val
}

// $2004: OAMDATA read
func (ppu *PPU) ReadOAMData() byte {
	res := ppu.primaryOAM.GetByte(ppu.oamAddr)
	ppu.openbus = res
	return res
}

// $2004: OAMDATA write
func (ppu *PPU) WriteOAMData(val byte) {
	ppu.openbus = val

	// https://www.nesdev.org/wiki/PPU_registers#OAM_data_($2004)_%3C%3E_read/write
	// > Writes to OAMDATA during rendering (on the pre-render line and the visible lines 0-239, provided either sprite or background rendering is enabled) do not modify values in OAM
	if !((ppu.scanLine < 240 || ppu.scanLine == 261) && (ppu.mask.ShowBackground() || ppu.mask.ShowSprites())) {
		// https://www.nesdev.org/wiki/PPU_OAM#Byte_2
		// > The three unimplemented bits of each sprite's byte 2 do not exist in the PPU and always read back as 0 on PPU revisions that allow reading PPU OAM through OAMDATA ($2004).
		// > This can be emulated by ANDing byte 2 with $E3 either when writing to or when reading from OAM.
		if (ppu.oamAddr & 0x03) == 0x02 {
			val &= 0xE3
		}
		ppu.primaryOAM.SetByte(ppu.oamAddr, val)
		ppu.oamAddr++
	} else {
		// https://www.nesdev.org/wiki/PPU_registers#OAM_data_($2004)_%3C%3E_read/write
		// > but do perform a glitchy increment of OAMADDR, bumping only the high 6 bits (i.e., it bumps the [n] value in PPU sprite evaluation
		// https://forums.nesdev.org/viewtopic.php?t=14140
		// 今はeval spriteの処理をまとめてやってるのでこれは影響ないはず
		ppu.oamAddr += 4
	}
}

func (ppu *PPU) ReadScroll() byte {
	// note: $2005 write only
	return ppu.openbus
}

// $2005: PPUSCROLL
func (ppu *PPU) WriteScroll(val byte) {
	ppu.openbus = val
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

func (ppu *PPU) ReadPPUAddr() byte {
	return ppu.openbus
}

// $2006: PPUADDR
func (ppu *PPU) WritePPUAddr(val byte) {
	ppu.openbus = val
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
	res := ppu.readPPUData(ppu.v)
	ppu.openbus = res

	if (ppu.scanLine < 240 || ppu.scanLine == 261) && (ppu.mask.ShowBackground() || ppu.mask.ShowSprites()) {
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
		res := ppu.paletteTable.Read(parsePaletteAddr(byte(addr - 0x3F00)))
		ppu.buf = ppu.vram.Read(addr - 0x1000)
		return res
	case 0x3F20 <= addr && addr <= 0x3FFF:
		// Mirrors of $3F00-$3F1F
		return ppu.readPPUData(0x3F00 + addr%0x20)
	default:
		panic(fmt.Sprintf("readPPUData invalid addr = 0x%04x", addr))
	}
}

// $2007: PPUDATA write
func (ppu *PPU) WritePPUData(val byte) {
	ppu.openbus = val
	ppu.writePPUData(ppu.v, val)

	if (ppu.scanLine < 240 || ppu.scanLine == 261) && (ppu.mask.ShowBackground() || ppu.mask.ShowSprites()) {
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
		ppu.paletteTable.Write(parsePaletteAddr(byte(addr-0x3F00)), val)
	case 0x3F20 <= addr && addr <= 0x3FFF:
		// Mirrors of $3F00-$3F1F
		ppu.writePPUData(0x3F00+addr%0x20, val)
	default:
		panic("uaaaaaaaaaaaaaaa")
	}
}

// $4014: OAMDMA write
func (ppu *PPU) WriteOAMDMA(data []byte) {
	for i := 0; i < len(data); i++ {
		ppu.primaryOAM.SetByte(ppu.oamAddr, data[i])
		ppu.oamAddr++
	}
	ppu.cpu.DMASuspend()
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
	s := ppu.secondaryOAM.GetSpriteByIndex(byte(sidx))
	lo, hi, ok := func() (byte, byte, bool) {
		if ppu.ctrl.SpriteSize() == 8 {
			y := uint16(ppu.scanLine) - uint16(s.y)
			if y > 7 {
				// eval時点で範囲内しか見ない&0xFFで初期化されるが、初期化途中でsprite&bgともにdisableされて前回の状態が残ることがある
				// eval時点だけでなくこのfetchタイミングでも範囲内か確認して少なくとも場外のspriteは表示させないようにしておく
				return 0, 0, false
			}
			if s.attributes.FlipSpriteVertically() {
				y = 7 - y
			}
			addr := ppu.ctrl.SpritePatternAddr() | (uint16(s.tileIndex) << 4) | y
			lo := ppu.mapper.Read(addr)
			hi := ppu.mapper.Read(addr + 8)
			return lo, hi, true
		} else {
			// 8x16
			// https://www.nesdev.org/wiki/PPU_OAM#Byte_1
			// > For 8x16 sprites, the PPU ignores the pattern table selection and selects a pattern table from bit 0 of this number.
			bankTile := (uint16(s.tileIndex) & 0b1) * 0x1000
			tileIndex := uint16(s.tileIndex) & 0b11111110
			y := uint16(ppu.scanLine) - uint16(s.y)
			if y > 15 {
				return 0, 0, false
			}
			if s.attributes.FlipSpriteVertically() {
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
		x:          s.x,
		attributes: s.attributes,
		lo:         lo,
		hi:         hi,
	}
	ppu.spriteFounds++
}

func (ppu *PPU) evalSpriteForNextScanline() {
	if ppu.scanLine >= 240 {
		panic("uaaaaaaaaaaaaxxxxxaaaa")
	}

	// debug
	// before secondaryOAM
	// for i := 0; i < 8; i++ {
	// 	s := ppu.secondaryOAM.GetSpriteByIndex(byte(i))
	// 	fmt.Printf("evalSpriteForNextScanline: before secondaryOAM[%d] = %v\n", i, *s)
	// }

	sidx := byte(0)
	for i := byte(0); i < 64; i++ {
		s := ppu.primaryOAM.GetSpriteByIndex(i)
		// in y range
		d := ppu.ctrl.SpriteSize()

		if uint(s.y) <= uint(ppu.scanLine) && uint(ppu.scanLine) < uint(s.y)+uint(d) {
			// debug
			//fmt.Printf("in y range: set secondaryOAM[%d] ppu.scanline:%d, s.y:%d\n", sidx, ppu.scanLine, s.y)

			if sidx < 8 {
				ppu.secondaryOAM.SetSpriteByIndex(sidx, *s)
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
	return ppu.scanLine == 261 || ppu.scanLine < 240
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

func (ppu *PPU) backgroundPixelPaletteAddr() *paletteAddr {
	if !ppu.mask.ShowBackground() {
		return newPaletteAddr(false, 0, 0)
	}
	x := ppu.Cycle - 1
	if x < 8 && !ppu.mask.ShowBackgroundLeftMost8pxlScreen() {
		return newPaletteAddr(false, 0, 0)
	}

	talo := byte((ppu.patternAttributeLowBit >> (15 - ppu.x)) & 0b1)
	tahi := byte((ppu.patternAttributeHighBit >> (15 - ppu.x)) & 0b1)

	ttlo := byte((ppu.patternTableLowBit >> (15 - ppu.x)) & 0b1)
	tthi := byte((ppu.patternTableHighBit >> (15 - ppu.x)) & 0b1)

	return newPaletteAddr(false, (tahi<<1)|talo, (tthi<<1)|ttlo)
}

func (ppu *PPU) spritePixelPaletteAddrAndSlotIndex() (*paletteAddr, int) {
	if !ppu.mask.ShowSprites() {
		return newPaletteAddr(true, 0, 0), -1
	}
	x := ppu.Cycle - 1
	if x < 8 && !ppu.mask.ShowSpritesLeftMost8pxlScreen() {
		return newPaletteAddr(true, 0, 0), -1
	}
	for i := 0; i < ppu.spriteFounds; i++ {
		s := ppu.spriteSlots[i]
		if int(s.x) <= x && x < int(s.x)+8 {
			// in range
			dx := x - int(s.x)
			if s.attributes.FlipSpriteHorizontally() {
				dx = 7 - dx
			}
			hb := (s.hi & (1 << (7 - dx))) >> (7 - dx)
			lb := (s.lo & (1 << (7 - dx))) >> (7 - dx)
			p := (hb << 1) | lb
			if p == 0 {
				// transparent pixel
				// https://www.nesdev.org/wiki/PPU_rendering#Preface
				// > The current pixel for each "active" sprite is checked (from highest to lowest priority),
				// > and the first non-transparent pixel moves on to a multiplexer, where it joins the BG pixel.
				continue
			}
			return newPaletteAddr(true, s.attributes.Palette(), p), i
		}
	}
	return newPaletteAddr(true, 0, 0), -1
}

func (ppu *PPU) multiplexPixelPaletteAddr() *paletteAddr {
	bp := ppu.backgroundPixelPaletteAddr()
	sp, slotIdx := ppu.spritePixelPaletteAddrAndSlotIndex()

	if bp.pixel == 0 && sp.pixel == 0 {
		// 0x3F00, universal background color
		return &paletteAddr{}
	} else if bp.pixel == 0 && sp.pixel != 0 {
		return sp
	} else if bp.pixel != 0 && sp.pixel == 0 {
		return bp
	} else {
		// bp != 0 && sp != 0
		s := ppu.spriteSlots[slotIdx]
		x := ppu.Cycle - 1

		if x < 255 && slotIdx == 0 {
			ppu.status.SetSprite0Hit(true)
		}
		if s.attributes.BehindBackground() {
			return bp
		} else {
			return sp
		}
	}
}

func (ppu *PPU) renderPixel() {
	x := ppu.Cycle - 1 // visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	y := ppu.scanLine

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
		if ppu.f == 1 && ppu.scanLine == 261 && ppu.Cycle == 339 {
			// skip 1 cycle
			ppu.Cycle = 340
		}
	}

	ppu.Cycle++
	if ppu.Cycle > 340 {
		ppu.Cycle = 0
		ppu.scanLine++
		if ppu.scanLine > 261 {
			ppu.scanLine = 0
			ppu.suppressVBlankFlag = false
			ppu.f ^= 1
		}
	}

	rendering := ppu.mask.ShowBackground() || ppu.mask.ShowSprites()
	preLine := ppu.scanLine == 261
	renderLine := preLine || ppu.visibleScanLine()
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
		if ppu.visibleScanLine() && visibleCycle {
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
		if 1 <= ppu.Cycle && ppu.Cycle <= 64 && ppu.visibleScanLine() {
			if ppu.Cycle%2 == 1 {
				addr := ppu.Cycle / 2
				ppu.secondaryOAM.SetByte(byte(addr), 0xFF)
			}
		}

		// sprite eval for next scanline
		// 65 <= ppu.Cycle <= 256
		if ppu.Cycle == 256 && ppu.visibleScanLine() {
			ppu.evalSpriteForNextScanline()
		}
		// sprite fetch
		if 257 <= ppu.Cycle && ppu.Cycle <= 320 && (preLine || ppu.visibleScanLine()) {
			// https://www.nesdev.org/wiki/PPU_registers#OAM_address_($2003)_%3E_write
			// > Values during rendering
			// > OAMADDR is set to 0 during each of ticks 257-320 (the sprite tile loading interval) of the pre-render and visible scanlines.
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
				if ppu.visibleScanLine() {
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
	if ppu.scanLine == 241 && ppu.Cycle == 1 {
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
