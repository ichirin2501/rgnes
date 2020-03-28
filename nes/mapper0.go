package nes

type Mapper0 struct {
	*Cassette
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
	case 0x8000 <= addr && addr < 0xC000:
		return m.PRG[addr-0x8000]
	case 0xC000 <= addr && addr <= 0xFFFF:
		index := (len(m.PRG)/0x4000-1)*0x4000 + int(addr-0xC000)
		return m.PRG[index]
	}
	panic("Unable to reach here?")
}

func (m *Mapper0) Write(addr uint16, val byte) {
	switch {
	case 0x0000 <= addr && addr < 0x2000:
		m.CHR[addr] = val
	}
	panic("Unable to reach here?")
}
