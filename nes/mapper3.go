package nes

import "fmt"

type Mapper3 struct {
	*Cassette
	chrBank  int
	prgBank1 int
	prgBank2 int
}

func NewMapper3(c *Cassette) *Mapper3 {
	prgBanks := len(c.PRG) / 0x4000
	prgBank2 := prgBanks - 1
	return &Mapper3{
		Cassette: c,
		chrBank:  0,
		prgBank1: 0,
		prgBank2: prgBank2,
	}
}

func (m *Mapper3) String() string {
	return "Mapper 3"
}

func (m *Mapper3) Reset() {
	// nothing
}

func (m *Mapper3) Read(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		index := m.chrBank*0x2000 + int(addr)
		return m.CHR[index]
	case 0x6000 <= addr && addr < 0x8000:
		return m.SRAM[addr-0x6000]
	case 0x8000 <= addr && addr < 0xC000:
		index := m.prgBank1*0x4000 + int(addr-0x8000)
		return m.PRG[index]
	case 0xC000 <= addr && addr <= 0xFFFF:
		index := m.prgBank2*0x4000 + int(addr-0xC000)
		return m.PRG[index]
	default:
		panic(fmt.Sprintf("Unable to reach Mapper3.Read(0x%04x)", addr))
	}
}

func (m *Mapper3) Write(addr uint16, val byte) {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		// read only (for ppu_read_buffer test)
	case 0x6000 <= addr && addr < 0x8000:
		m.SRAM[addr-0x6000] = val
	case 0x8000 <= addr:
		// https://www.nesdev.org/wiki/INES_Mapper_003#Bank_select_($8000-$FFFF)
		m.chrBank = int(val & 0x3)
	default:
		panic(fmt.Sprintf("Unable to reach Mapper3.Write(0x%04x) = 0x%02x", addr, val))
	}
}
