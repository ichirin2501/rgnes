package nes

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

const (
	iNESMagicNumber  = 0x1a53454e
	programROMUnit   = 0x4000
	characterROMUnit = 0x2000
)

type iNESHeader struct {
	Magic   uint32
	PRGSize byte
	CHRSize byte
	Flags6  byte // Mapper, mirroring, battery, trainer
	Flags7  byte // Mapper, VS/Playchoice, NES 2.0
	Flags8  byte // PRG-RAM size (rarely used extension)
	Flags9  byte // TV system (rarely used extension)
	Flags10 byte // TV system, PRG-RAM presence (unofficial, rarely used extension)
	_       [5]byte
}

type Cassette struct {
	PRG    []byte
	CHR    []byte
	Mapper byte
}

func NewCassette(path string) (*Cassette, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	header := &iNESHeader{}
	if err := binary.Read(f, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	if header.Magic != iNESMagicNumber {
		return nil, errors.New("invalid file")
	}

	mapper := (header.Flags7 & 0xF0) | (header.Flags6&0xF0)>>4

	prgRom := make([]byte, int(header.PRGSize)*programROMUnit)
	if _, err := io.ReadFull(f, prgRom); err != nil {
		return nil, err
	}

	chrRom := make([]byte, int(header.CHRSize)*characterROMUnit)
	if _, err := io.ReadFull(f, chrRom); err != nil {
		return nil, err
	}

	return &Cassette{
		PRG:    prgRom,
		CHR:    chrRom,
		Mapper: mapper,
	}, nil
}
