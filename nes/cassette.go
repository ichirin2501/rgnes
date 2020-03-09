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

type iNESHeader [NESHeaderSize]byte

func (h iNESHeader) ProgramSize() int {
	return int(h[4]) * ProgramROMUnit
}
func (h iNESHeader) CharacterSize() int {
	return int(h[5]) * CharacterROMUnit
}

type Cassette struct {
	Header       iNESHeader
	ProgramROM   []byte
	CharacterROM []byte
}

func NewCassette(path string) (*Cassette, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var header iNESHeader
	hn, err := f.ReadAt(header[:], 0)
	if err != nil {
		return nil, err
	}
	if hn != NESHeaderSize {
		return nil, errors.New("fail read header")
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
