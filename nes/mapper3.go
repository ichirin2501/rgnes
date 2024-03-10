package nes

import "fmt"

type mapper3 struct {
	*Cassette
	chrBank byte
}

func newMapper3(c *Cassette) *mapper3 {
	return &mapper3{
		Cassette: c,
		chrBank:  0,
	}
}

func (m *mapper3) String() string {
	return "Mapper 3"
}

func (m *mapper3) Reset() {
	// nothing
}

func (m *mapper3) Read(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		// > PPU $0000-$1FFF: 8 KB switchable CHR ROM bank
		index := int(m.chrBank)*0x2000 + int(addr)
		return m.CHR[index]
	case 0x6000 <= addr && addr < 0x8000:
		// mapper 3 don't have PRG RAM
		return 0
	case 0x8000 <= addr && addr <= 0xFFFF:
		// > PRG ROM size: 16 KiB or 32 KiB
		// > PRG ROM bank size: Not bankswitched
		index := int(addr-0x8000) % len(m.PRG)
		return m.PRG[index]
	default:
		panic(fmt.Sprintf("Unable to reach %s Read(0x%04x)", m, addr))
	}
}

func (m *mapper3) Write(addr uint16, val byte) {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		// read only (for ppu_read_buffer test)
	case 0x6000 <= addr && addr < 0x8000:
		// mapper 3 don't have PRG RAM
	case 0x8000 <= addr && addr <= 0xFFFF:
		// https://www.nesdev.org/wiki/INES_Mapper_003#Bank_select_($8000-$FFFF)
		// 7  bit  0
		// ---- ----
		// cccc ccCC
		// |||| ||||
		// ++++-++++- Select 8 KB CHR ROM bank for PPU $0000-$1FFF
		// > CNROM only implements the lowest 2 bits, capping it at 32 KiB CHR. Other boards may implement 4 or more bits for larger CHR.
		m.chrBank = val
	default:
		panic(fmt.Sprintf("Unable to reach %s Write(0x%04x) = 0x%02x", m, addr, val))
	}
}
