package cpu

import (
	"fmt"
)

type APU interface {
	// todo
}

type Joypad interface {
	Read() byte
	Write(byte)
}

type Mapper interface {
	Read(uint16) byte
	Write(uint16, byte)
}

type PPU interface {
	Step()
	ReadController() byte
	ReadMask() byte
	ReadStatus() byte
	ReadOAMAddr() byte
	ReadOAMData() byte
	ReadScroll() byte
	ReadPPUAddr() byte
	ReadPPUData() byte

	PeekController() byte
	PeekMask() byte
	PeekStatus() byte
	PeekOAMAddr() byte
	PeekOAMData() byte
	PeekScroll() byte
	PeekPPUAddr() byte
	PeekPPUData() byte

	WriteController(byte)
	WriteMask(byte)
	WriteStatus(byte)
	WriteOAMAddr(byte)
	WriteOAMData(byte)
	WriteScroll(byte)
	WritePPUAddr(byte)
	WritePPUData(byte)
	WriteOAMDMA([]byte)
}

type Bus struct {
	ram    []byte
	ppu    PPU
	apu    APU
	Mapper Mapper
	joypad Joypad

	clock int
}

func NewBus(ppu PPU, apu APU, mapper Mapper, joypad Joypad) *Bus {
	return &Bus{
		ram:    make([]byte, 2048),
		ppu:    ppu,
		apu:    apu,
		Mapper: mapper,
		joypad: joypad,
	}
}

func (bus *Bus) Read(addr uint16) byte {
	bus.clock++
	bus.ppu.Step()
	bus.ppu.Step()
	bus.ppu.Step()
	return bus.read(addr)
}
func (bus *Bus) read(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		// 2KB internal RAM
		return bus.ram[addr%0x800]
	case 0x2000 <= addr && addr <= 0x2007:
		// NES PPU registers
		switch {
		case addr == 0x2000:
			return bus.ppu.ReadController()
		case addr == 0x2001:
			return bus.ppu.ReadMask()
		case addr == 0x2002:
			return bus.ppu.ReadStatus()
		case addr == 0x2003:
			return bus.ppu.ReadOAMAddr()
		case addr == 0x2004:
			return bus.ppu.ReadOAMData()
		case addr == 0x2005:
			return bus.ppu.ReadScroll()
		case addr == 0x2006:
			return bus.ppu.ReadPPUAddr()
		case addr == 0x2007:
			return bus.ppu.ReadPPUData()
		default:
			panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Read", addr))
		}
	case 0x2008 <= addr && addr <= 0x3FFF:
		// Mirrors of $2000-2007 (repeats every 8 bytes)
		return bus.read(0x2000 + addr%0x08)
	case 0x4000 <= addr && addr <= 0x4017:
		// NES APU and I/O registers
		switch {
		case addr == 0x4016:
			return bus.joypad.Read()
		case addr == 0x4017: // TODO: 2p
			//panic("unimplemented Bus.Read 0x4017(2p keypad)")
			return 0
		default:
			// basically, ignore
			return 0
		}
	case 0x4018 <= addr && addr <= 0x401F:
		// APU and I/O functionality that is normally disabled.
		return 0
	case 0x4020 <= addr && addr <= 0xFFFF:
		// Cartridge space
		return bus.Mapper.Read(addr)
	default:
		panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Read", addr))
	}
}

func (bus *Bus) Write(addr uint16, val byte) {
	bus.clock++
	bus.ppu.Step()
	bus.ppu.Step()
	bus.ppu.Step()
	bus.write(addr, val)
}
func (bus *Bus) write(addr uint16, val byte) {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		// 2KB internal RAM
		bus.ram[addr%0x800] = val
	case 0x2000 <= addr && addr <= 0x2007:
		// NES PPU registers
		switch {
		case addr == 0x2000:
			bus.ppu.WriteController(val)
		case addr == 0x2001:
			bus.ppu.WriteMask(val)
		case addr == 0x2002:
			bus.ppu.WriteStatus(val)
		case addr == 0x2003:
			bus.ppu.WriteOAMAddr(val)
		case addr == 0x2004:
			bus.ppu.WriteOAMData(val)
		case addr == 0x2005:
			bus.ppu.WriteScroll(val)
		case addr == 0x2006:
			bus.ppu.WritePPUAddr(val)
		case addr == 0x2007:
			bus.ppu.WritePPUData(val)
		default:
			panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Write", addr))
		}
	case 0x2008 <= addr && addr <= 0x3FFF:
		// Mirrors of $2000-2007 (repeats every 8 bytes)
		bus.write(0x2000+addr%0x08, val)
	case 0x4000 <= addr && addr <= 0x4017:
		// NES APU and I/O registers
		switch {
		case addr == 0x4014:
			buf := make([]byte, 256)
			a := uint16(val) << 8
			for i := uint16(0); i < 256; i++ {
				buf[i] = bus.read(a + i)
			}
			bus.ppu.WriteOAMDMA(buf)
		case addr == 0x4016:
			bus.joypad.Write(val)
		default:
			// basically, ignore
		}
	case 0x4018 <= addr && addr <= 0x401F:
		// APU and I/O functionality that is normally disabled.
	case 0x4020 <= addr && addr <= 0xFFFF:
		// Cartridge space
		bus.Mapper.Write(addr, val)
	default:
		panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Write", addr))
	}
}

// Peek is used for debugging
func (bus *Bus) Peek(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		return bus.ram[addr%0x800]
	case 0x2000 <= addr && addr <= 0x2007:
		// NES PPU registers
		switch {
		case addr == 0x2000:
			return bus.ppu.PeekController()
		case addr == 0x2001:
			return bus.ppu.PeekMask()
		case addr == 0x2002:
			return bus.ppu.PeekStatus()
		case addr == 0x2003:
			return bus.ppu.PeekOAMAddr()
		case addr == 0x2004:
			return bus.ppu.PeekOAMData()
		case addr == 0x2005:
			return bus.ppu.PeekScroll()
		case addr == 0x2006:
			return bus.ppu.PeekPPUAddr()
		case addr == 0x2007:
			return bus.ppu.PeekPPUData()
		default:
			panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Peek", addr))
		}

	case 0x2008 <= addr && addr <= 0x3FFF:
		// Mirrors of $2000-2007 (repeats every 8 bytes)
		return bus.Peek(0x2000 + addr%0x08)
	case 0x4000 <= addr && addr <= 0x4017:
		// todo
		return 0
	case 0x4018 <= addr && addr <= 0x401F:
		// todo
		return 0
	case 0x4020 <= addr && addr <= 0xFFFF:
		return bus.Mapper.Read(addr)
	default:
		//return 0
		panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Peek", addr))
	}

}
