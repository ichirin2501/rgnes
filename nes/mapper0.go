package nes

import "fmt"

type mapper0 struct {
	*Cassette
	SRAM []byte
}

func newMapper0(c *Cassette) *mapper0 {
	return &mapper0{
		Cassette: c,
		// > PRG RAM: 2 or 4 KiB, not bankswitched, only in Family Basic (but most emulators provide 8)
		SRAM: make([]byte, 0x2000), // 8KiB
	}
}

func (m *mapper0) String() string {
	return "Mapper 0"
}

func (m *mapper0) Reset() {
	// nothing
}

func (m *mapper0) Read(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		return m.readCHR(addr)
	case 0x6000 <= addr && addr < 0x8000:
		return m.SRAM[addr-0x6000]
	case 0x8000 <= addr && addr < 0xC000:
		index := int(addr - 0x8000)
		return m.PRG[index]
	case 0xC000 <= addr && addr <= 0xFFFF:
		// > CPU $8000-$BFFF: First 16 KB of ROM.
		// > CPU $C000-$FFFF: Last 16 KB of ROM (NROM-256) or mirror of $8000-$BFFF (NROM-128).
		// len(m.PRG) = 16K or 32K
		index := 0x4000*(len(m.PRG)/0x4000-1) + int(addr-0xC000)
		return m.PRG[index]
	default:
		panic(fmt.Sprintf("Unable to reach %s Read(0x%04x)", m, addr))
	}
}

func (m *mapper0) Write(addr uint16, val byte) {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		// https://www.nesdev.org/wiki/NROM
		// > CHR capacity: 8 KiB ROM (DIP-28 standard pinout) but most emulators support RAM
		m.writeCHR(addr, val)
	case 0x6000 <= addr && addr < 0x8000:
		m.SRAM[addr-0x6000] = val
	case 0x8000 <= addr && addr <= 0xFFFF:
		// nothing
	default:
		panic(fmt.Sprintf("Unable to reach %s Write(0x%04x) = 0x%02x", m, addr, val))
	}
}
