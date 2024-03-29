package nes

import (
	"fmt"
	"io"
)

type Mapper interface {
	Read(addr uint16) byte
	Write(addr uint16, val byte)
	MirroingType() MirroringType
	String() string
	Reset()
}

func NewMapper(r io.Reader) (Mapper, error) {
	c, err := NewCassette(r)
	if err != nil {
		return nil, err
	}
	return NewMapperFromCassette(c), nil
}

func NewMapperFromCassette(c *Cassette) Mapper {
	switch c.Mapper {
	case 0:
		return newMapper0(c)
	case 2:
		return newMapper2(c)
	case 3:
		return newMapper3(c)
	}
	panic(fmt.Sprintf("Unsupported mapper: %0x", c.Mapper))
}
