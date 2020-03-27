package nes

import "fmt"

type Mapper interface {
	Memory
}

type Mapper0 struct {
	*Cassette
}

func (m *Mapper0) Read(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return m.CHR[addr]
	case addr >= 0xC000:
		index := (len(m.PRG)/0x4000-1)*0x4000 + int(addr-0xC000)
		return m.PRG[index]
	case addr >= 0x8000:
		return m.PRG[addr-0x8000]
	case addr >= 0x6000:
		// TODO:
	}
	panic("Unable to reach here?")
}
func (m *Mapper0) Write(addr uint16, val byte) {
	// TODO
}

func NewMapper(c *Cassette) Mapper {
	switch c.Mapper {
	case 0:
		return &Mapper0{c}
	}
	panic(fmt.Sprintf("Unsupported mapper: %0x", c.Mapper))
}
