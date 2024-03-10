package nes

import "fmt"

type mapper2 struct {
	*Cassette
	prgBank byte
}

func newMapper2(c *Cassette) *mapper2 {
	return &mapper2{
		Cassette: c,
		prgBank:  0,
	}
}

func (m *mapper2) String() string {
	return "Mapper 2"
}

func (m *mapper2) Reset() {
	// nothing
}

func (m *mapper2) Read(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		return m.CHR[addr]
	case 0x6000 <= addr && addr < 0x8000:
		// mapper2 dont'h have PRG RAM
		return 0
	case 0x8000 <= addr && addr < 0xC000:
		// > CPU $8000-$BFFF: 16 KB switchable PRG ROM bank
		index := 0x4000*int(m.prgBank) + int(addr-0x8000)
		return m.PRG[index]
	case 0xC000 <= addr && addr <= 0xFFFF:
		// > CPU $C000-$FFFF: 16 KB PRG ROM bank, fixed to the last bank
		index := 0x4000*(len(m.PRG)/0x4000-1) + int(addr-0xC000)
		return m.PRG[index]
	default:
		panic(fmt.Sprintf("Unable to reach %s Read(0x%04x)", m, addr))
	}
}

func (m *mapper2) Write(addr uint16, val byte) {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		m.CHR[addr] = val
	case 0x6000 <= addr && addr < 0x8000:
		// mapper2 dont'h have PRG RAM
	case 0x8000 <= addr && addr <= 0xFFFF:
		// 7  bit  0
		// ---- ----
		// xxxx pPPP
		//      ||||
		//      ++++- Select 16 KB PRG ROM bank for CPU $8000-$BFFF
		//            (UNROM uses bits 2-0; UOROM uses bits 3-0)
		m.prgBank = val & 0x0F
	default:
		panic(fmt.Sprintf("Unable to reach %s Write(0x%04x) = 0x%02x", m, addr, val))
	}
}
