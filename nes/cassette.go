package nes

import (
	"encoding/binary"
	"errors"
	"io"
)

type MirroringType int

const (
	iNESMagicNumber  = 0x1a53454e
	programROMUnit   = 0x4000
	characterROMUnit = 0x2000

	MirroringVertical MirroringType = iota
	MirroringHorizontal
	MirroringFourScreen
)

func (m *MirroringType) IsVertical() bool {
	return *m == MirroringVertical
}
func (m *MirroringType) IsHorizontal() bool {
	return *m == MirroringHorizontal
}
func (m *MirroringType) IsFourScreen() bool {
	return *m == MirroringFourScreen
}

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
	Mirror MirroringType
}

func NewCassette(r io.Reader) (*Cassette, error) {
	header := &iNESHeader{}
	if err := binary.Read(r, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	if header.Magic != iNESMagicNumber {
		return nil, errors.New("invalid ines file")
	}
	if ((header.Flags7 >> 2) & 0x02) == 2 {
		return nil, errors.New("NES2.0 format is not supported yet")
	}

	var mirroringType MirroringType
	mirrorFourScreenFlag := false
	mirrorVerticalFlag := false
	if (header.Flags6 & 0x08) != 0 {
		mirrorFourScreenFlag = true
	}
	if (header.Flags6 & 0x01) != 0 {
		mirrorVerticalFlag = true
	}
	if mirrorFourScreenFlag {
		mirroringType = MirroringFourScreen
	} else if mirrorVerticalFlag {
		mirroringType = MirroringVertical
	} else {
		mirroringType = MirroringHorizontal
	}

	mapper := (header.Flags7 & 0xF0) | (header.Flags6&0xF0)>>4

	prgRom := make([]byte, int(header.PRGSize)*programROMUnit)
	if _, err := io.ReadFull(r, prgRom); err != nil {
		return nil, err
	}

	chrRom := make([]byte, int(header.CHRSize)*characterROMUnit)
	if _, err := io.ReadFull(r, chrRom); err != nil {
		return nil, err
	}
	if header.CHRSize == 0 {
		chrRom = make([]byte, 8192)
	}

	return &Cassette{
		PRG:    prgRom,
		CHR:    chrRom,
		Mapper: mapper,
		Mirror: mirroringType,
	}, nil
}

func (c *Cassette) MirroingType() MirroringType {
	return c.Mirror
}
