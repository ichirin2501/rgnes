package nes

import (
	"errors"
	"os"
)

const (
	nesHeaderSize    = 0x0010
	programROMUnit   = 0x4000
	characterROMUnit = 0x2000
)

type iNESHeader [nesHeaderSize]byte

func (h iNESHeader) ProgramSize() int {
	return int(h[4]) * programROMUnit
}
func (h iNESHeader) CharacterSize() int {
	return int(h[5]) * characterROMUnit
}

type rom []byte

func (r rom) Read(offset uint16) byte {
	return r[offset]
}

type Cassette struct {
	Header       iNESHeader
	ProgramROM   MemoryReader
	CharacterROM MemoryReader
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
	if hn != nesHeaderSize {
		return nil, errors.New("fail read header")
	}

	// read prg-rom
	prgRom := make([]byte, header.ProgramSize())
	n, err := f.ReadAt(prgRom, nesHeaderSize)
	if err != nil {
		return nil, err
	}
	if n != header.ProgramSize() {
		return nil, errors.New("fail read prg-rom")
	}

	// read chr-rom
	chrRom := make([]byte, header.CharacterSize())
	m, err := f.ReadAt(chrRom, nesHeaderSize+int64(header.ProgramSize()))
	if err != nil {
		return nil, err
	}
	if m != header.CharacterSize() {
		return nil, errors.New("fail read chr-rom")
	}

	return &Cassette{
		Header:       header,
		ProgramROM:   rom(prgRom),
		CharacterROM: rom(chrRom),
	}, nil
}
