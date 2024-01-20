package nes

import (
	"fmt"
)

type Bus struct {
	ram    []byte
	ppu    *PPU
	apu    *APU
	mapper Mapper
	joypad *Joypad

	// This clock is used to adjust the clock difference for each instruction,
	// so keep the state separate from the $4014 dma stall.
	clock int
	stall int
}

func NewBus(ppu *PPU, apu *APU, mapper Mapper, joypad *Joypad) *Bus {
	return &Bus{
		ram:    make([]byte, 2048),
		ppu:    ppu,
		apu:    apu,
		mapper: mapper,
		joypad: joypad,
	}
}

func (bus *Bus) Read(addr uint16) byte {
	bus.tick(1)
	return bus._read(addr)
}
func (bus *Bus) _read(addr uint16) byte {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		// 2KB internal RAM
		return bus.ram[addr%0x800]
	case 0x2000 <= addr && addr <= 0x2007:
		// NES PPU registers
		switch {
		case addr == 0x2000:
			return bus.ppu.readController()
		case addr == 0x2001:
			return bus.ppu.readMask()
		case addr == 0x2002:
			return bus.ppu.readStatus()
		case addr == 0x2003:
			return bus.ppu.readOAMAddr()
		case addr == 0x2004:
			return bus.ppu.readOAMData()
		case addr == 0x2005:
			return bus.ppu.readScroll()
		case addr == 0x2006:
			return bus.ppu.readPPUAddr()
		case addr == 0x2007:
			return bus.ppu.readPPUData()
		default:
			panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Read", addr))
		}
	case 0x2008 <= addr && addr <= 0x3FFF:
		// Mirrors of $2000-2007 (repeats every 8 bytes)
		return bus._read(0x2000 + addr%0x08)
	case 0x4000 <= addr && addr <= 0x4017:
		// NES APU and I/O registers
		switch {
		case addr == 0x4015:
			return bus.apu.readStatus()
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
		return bus.mapper.Read(addr)
	default:
		panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Read", addr))
	}
}

func (bus *Bus) Write(addr uint16, val byte) {
	bus.tick(1)
	bus._write(addr, val)
}
func (bus *Bus) _write(addr uint16, val byte) {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		// 2KB internal RAM
		bus.ram[addr%0x800] = val
	case 0x2000 <= addr && addr <= 0x2007:
		// NES PPU registers
		switch {
		case addr == 0x2000:
			bus.ppu.writeController(val)
		case addr == 0x2001:
			bus.ppu.writeMask(val)
		case addr == 0x2002:
			bus.ppu.writeStatus(val)
		case addr == 0x2003:
			bus.ppu.writeOAMAddr(val)
		case addr == 0x2004:
			bus.ppu.writeOAMData(val)
		case addr == 0x2005:
			bus.ppu.writeScroll(val)
		case addr == 0x2006:
			bus.ppu.writePPUAddr(val)
		case addr == 0x2007:
			bus.ppu.writePPUData(val)
		default:
			panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Write", addr))
		}
	case 0x2008 <= addr && addr <= 0x3FFF:
		// Mirrors of $2000-2007 (repeats every 8 bytes)
		bus._write(0x2000+addr%0x08, val)
	case 0x4000 <= addr && addr <= 0x4017:
		// NES APU and I/O registers
		switch {
		case addr == 0x4000:
			bus.apu.writePulse1Controller(val)
		case addr == 0x4001:
			bus.apu.writePulse1Sweep(val)
		case addr == 0x4002:
			bus.apu.writePulse1TimerLow(val)
		case addr == 0x4003:
			bus.apu.writePulse1LengthAndTimerHigh(val)
		case addr == 0x4004:
			bus.apu.writePulse2Controller(val)
		case addr == 0x4005:
			bus.apu.writePulse2Sweep(val)
		case addr == 0x4006:
			bus.apu.writePulse2TimerLow(val)
		case addr == 0x4007:
			bus.apu.writePulse2LengthAndTimerHigh(val)
		case addr == 0x4008:
			bus.apu.writeTriangleController(val)
		case addr == 0x4009:
			// unused
		case addr == 0x400A:
			bus.apu.writeTriangleTimerLow(val)
		case addr == 0x400B:
			bus.apu.writeTriangleLengthAndTimerHigh(val)
		case addr == 0x400C:
			bus.apu.writeNoiseController(val)
		case addr == 0x400D:
			// unused
		case addr == 0x400E:
			bus.apu.writeNoiseLoopAndPeriod(val)
		case addr == 0x400F:
			bus.apu.writeNoiseLength(val)
		case addr == 0x4010:
			bus.apu.writeDMCController(val)
		case addr == 0x4011:
			bus.apu.writeDMCLoadCounter(val)
		case addr == 0x4012:
			bus.apu.writeDMCSampleAddr(val)
		case addr == 0x4013:
			bus.apu.writeDMCSampleLength(val)
		case addr == 0x4014:
			a := uint16(val) << 8
			for i := uint16(0); i < 256; i++ {
				bus.ppu.writeOAMDMAByte(bus._read(a + i))
			}
			bus.tickStall(513 + bus.realClock()%2)
		case addr == 0x4015:
			bus.apu.writeStatus(val)
		case addr == 0x4016:
			bus.joypad.Write(val)
		case addr == 0x4017:
			bus.apu.writeFrameCounter(val)
		default:
			// basically, ignore
		}
	case 0x4018 <= addr && addr <= 0x401F:
		// APU and I/O functionality that is normally disabled.
	case 0x4020 <= addr && addr <= 0xFFFF:
		// Cartridge space
		bus.mapper.Write(addr, val)
	default:
		panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Write", addr))
	}
}

func (bus *Bus) realClock() int {
	return bus.clock + bus.stall
}

func (bus *Bus) tick(cpuCycle int) {
	bus.clock += cpuCycle
	for i := 0; i < cpuCycle; i++ {
		bus.apu.Step()

		bus.ppu.Step()
		bus.ppu.Step()
		bus.ppu.Step()
	}
}

func (bus *Bus) tickStall(cpuCycle int) {
	bus.stall += cpuCycle
	for i := 0; i < cpuCycle; i++ {
		bus.apu.Step()

		bus.ppu.Step()
		bus.ppu.Step()
		bus.ppu.Step()
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
		// NES APU and I/O registers
		// todo
		switch {
		case addr == 0x4015:
			return bus.apu.PeekStatus()
		case addr == 0x4016:
			return bus.joypad.Peek()
		default:
			return 0
		}
	case 0x4018 <= addr && addr <= 0x401F:
		// todo
		return 0
	case 0x4020 <= addr && addr <= 0xFFFF:
		return bus.mapper.Read(addr)
	default:
		//return 0
		panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Peek", addr))
	}

}
