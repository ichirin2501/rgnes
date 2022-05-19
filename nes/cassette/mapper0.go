package cassette

import "fmt"

type Mapper0 struct {
	*Cassette
	prgBanks int
	prgBank1 int
	prgBank2 int
}

func NewMapper0(c *Cassette) *Mapper0 {
	prgBanks := len(c.PRG) / 0x4000
	prgBank1 := 0
	prgBank2 := prgBanks - 1
	return &Mapper0{
		Cassette: c,
		prgBanks: prgBanks,
		prgBank1: prgBank1,
		prgBank2: prgBank2,
	}
}

func (m *Mapper0) String() string {
	return "Mapper 0"
}

func (m *Mapper0) Reset() {
	// nothing
}

func (m *Mapper0) Read(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		return m.CHR[addr]
	case 0x6000 <= addr && addr < 0x8000:
		return m.SRAM[addr-0x6000]
	case 0x8000 <= addr && addr < 0xC000:
		index := m.prgBank1*0x4000 + int(addr-0x8000)
		return m.PRG[index]
	case 0xC000 <= addr && addr <= 0xFFFF:
		index := m.prgBank2*0x4000 + int(addr-0xC000)
		return m.PRG[index]
	default:
		panic(fmt.Sprintf("Unable to reach Mapper0.Read(0x%04x)", addr))
	}
}

func (m *Mapper0) Write(addr uint16, val byte) {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		// https://www.nesdev.org/wiki/NROM
		// > CHR capacity: 8 KiB ROM (DIP-28 standard pinout) but most emulators support RAM
		m.CHR[addr] = val
	case 0x6000 <= addr && addr < 0x8000:
		m.SRAM[addr-0x6000] = val
	case 0x8000 <= addr:
		m.prgBank1 = int(val) % m.prgBanks
	default:
		panic(fmt.Sprintf("Unable to reach Mapper0.Write(0x%04x) = 0x%02x", addr, val))
	}
}
