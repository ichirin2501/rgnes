package nes

import "fmt"

type Mapper interface {
	Memory
	Reset()
}

func NewMapper(c *Cassette) Mapper {
	switch c.Mapper {
	case 0:
		return &Mapper0{c}
	}
	panic(fmt.Sprintf("Unsupported mapper: %0x", c.Mapper))
}
