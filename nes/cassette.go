package nes

import (
	"errors"
	"os"
)

const (
	NESHeaderSize    = 0x0010
	ProgramROMUnit   = 0x4000
	CharacterROMUnit = 0x2000
)

type iNESHeader struct {
	prgSize int
	chrSize int
	// TODO: flag
}

func (h *iNESHeader) ProgramSize() int {
	return h.prgSize * ProgramROMUnit
}
func (h *iNESHeader) CharacterSize() int {
	return h.chrSize * CharacterROMUnit
}

type ProgramROM []byte
type CharacterROM []byte

func (r *ProgramROM) Read(address uint16) byte {
	// TODO
	return 0
}

func (c *CharacterROM) Read(address uint16) byte {
	// TODO
	return 0
}

type Cassette struct {
	Header       *iNESHeader
	ProgramROM   ProgramROM
	CharacterROM CharacterROM
}

func NewCassette(path string) (*Cassette, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, NESHeaderSize)
	if _, err := f.ReadAt(buf, 0); err != nil {
		return nil, err
	}

	header := &iNESHeader{
		prgSize: int(buf[4]),
		chrSize: int(buf[5]),
	}

	// read prg-rom
	prgRom := make([]byte, header.ProgramSize())
	n, err := f.ReadAt(prgRom, NESHeaderSize)
	if err != nil {
		return nil, err
	}
	if n != header.ProgramSize() {
		return nil, errors.New("fail read prg-rom")
	}

	// read chr-rom
	chrRom := make([]byte, header.CharacterSize())
	m, err := f.ReadAt(chrRom, NESHeaderSize+int64(header.ProgramSize()))
	if err != nil {
		return nil, err
	}
	if m != header.CharacterSize() {
		return nil, errors.New("fail read chr-rom")
	}

	return &Cassette{
		Header:       header,
		ProgramROM:   prgRom,
		CharacterROM: chrRom,
	}, nil
}
