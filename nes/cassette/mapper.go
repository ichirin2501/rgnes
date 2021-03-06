package cassette

import (
	"fmt"
)

type Mapper interface {
	Read(addr uint16) byte
	Write(addr uint16, val byte)
	String() string
	Reset()
}

func NewMapper(c *Cassette) Mapper {
	switch c.Mapper {
	case 0:
		return &Mapper0{c}
	}
	panic(fmt.Sprintf("Unsupported mapper: %0x", c.Mapper))
}
