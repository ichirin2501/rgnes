package cpu

import (
	"fmt"

	"github.com/ichirin2501/rgnes/nes"
	"github.com/ichirin2501/rgnes/nes/memory"
	"github.com/ichirin2501/rgnes/nes/ppu"
)

type Bus struct {
	CPURAM []byte
	ppu    *ppu.PPU
	apu    memory.Memory
	Mapper memory.Memory
	joypad *nes.Joypad
}

func NewBus(ppu *ppu.PPU, apu memory.Memory, mapper memory.Memory, joypad *nes.Joypad) *Bus {
	return &Bus{
		CPURAM: make([]byte, 2048),
		ppu:    ppu,
		apu:    apu,
		Mapper: mapper,
		joypad: joypad,
	}
}

func (bus *Bus) Read(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		// 2KB internal RAM
		return bus.CPURAM[addr%0x800]
	case 0x2000 <= addr && addr <= 0x2007:
		// NES PPU registers
		switch {
		case addr == 0x2000:
			// ignore ppu read ctrl
			return 0
		case addr == 0x2001:
			// ignore ppu read mask
			return 0
		case addr == 0x2002:
			return bus.ppu.ReadStatus()
		case addr == 0x2003:
			// ignore ppu
			return 0
		case addr == 0x2004:
			return bus.ppu.ReadOAMData()
		case addr == 0x2005:
			return 0
		case addr == 0x2006:
			return 0
		case addr == 0x2007:
			return bus.ppu.ReadPPUData()
		default:
			panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Read", addr))
		}
	case 0x2008 <= addr && addr <= 0x3FFF:
		// Mirrors of $2000-2007 (repeats every 8 bytes)
		return bus.Read(0x2000 + ((addr - 0x2008) % 0x08))
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
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		// 2KB internal RAM
		bus.CPURAM[addr%0x800] = val
	case 0x2000 <= addr && addr <= 0x2007:
		// NES PPU registers
		switch {
		case addr == 0x2000:
			bus.ppu.WriteController(val)
		case addr == 0x2001:
			bus.ppu.WriteMask(val)
		case addr == 0x2002:
			// ignore
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
		bus.Write(0x2000+((addr-0x2008)%0x08), val)
	case 0x4000 <= addr && addr <= 0x4017:
		// NES APU and I/O registers
		switch {
		case addr == 0x4014:
			buf := make([]byte, 256)
			a := uint16(val) << 8
			for i := 0; i < 256; i++ {
				buf[i] = bus.CPURAM[(a+uint16(i))%0x800]
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

func (bus *Bus) ReadForTest(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return bus.CPURAM[addr%0x800]
	case 0x2000 <= addr && addr <= 0x2007:
		//fmt.Printf("[warn] read ppu data addr = 0x%04x\n", addr)
		return 0
	case 0x2008 <= addr && addr <= 0x3FFF:
		// Mirrors of $2000-2007 (repeats every 8 bytes)
		return bus.ReadForTest(0x2000 + ((addr - 0x2008) % 0x08))
	case addr == 0x4016: // TODO: keypad
		//panic("unimplemented Bus.Read 0x4016(keypad)")
		return 0
	case addr == 0x4017: // TODO: 2p
		//panic("unimplemented Bus.Read 0x4017(2p keypad)")
		return 0
	case addr < 0x6000:
		return 0
	case addr >= 0x6000:
		return bus.Mapper.Read(addr)
	case addr >= 0x4020:
		return 0
	default:
		//return 0
		panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.ReadForTest", addr))
	}

}
