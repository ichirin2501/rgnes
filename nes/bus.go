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
	dma    *DMA

	// This clock is used to adjust the clock difference for each instruction,
	// so keep the state separate from the $4014 dma stall.
	clock int
	stall int
}

func NewBus(ppu *PPU, apu *APU, mapper Mapper, joypad *Joypad, dma *DMA) *Bus {
	return &Bus{
		ram:    make([]byte, 2048),
		ppu:    ppu,
		apu:    apu,
		mapper: mapper,
		joypad: joypad,
		dma:    dma,
	}
}

func (bus *Bus) read(addr uint16) byte {
	switch {
	// 2KB internal RAM
	case 0x0000 <= addr && addr <= 0x1FFF:
		return bus.ram[addr%0x800]

	// NES PPU registers
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

	// Mirrors of $2000-2007 (repeats every 8 bytes)
	case 0x2008 <= addr && addr <= 0x3FFF:
		return bus.read(0x2000 + addr%0x08)

	// NES APU and I/O registers
	case 0x4000 <= addr && addr <= 0x4017:
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

	// APU and I/O functionality that is normally disabled.
	case 0x4018 <= addr && addr <= 0x401F:
		return 0

	// Cartridge space
	case 0x4020 <= addr && addr <= 0xFFFF:
		return bus.mapper.Read(addr)
	default:
		panic(fmt.Sprintf("Unable to reach addr:0x%0x in Bus.Read", addr))
	}
}

func (bus *Bus) write(addr uint16, val byte) {
	switch {
	// 2KB internal RAM
	case 0x0000 <= addr && addr <= 0x1FFF:
		bus.ram[addr%0x800] = val

	// NES PPU registers
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

	// Mirrors of $2000-2007 (repeats every 8 bytes)
	case 0x2008 <= addr && addr <= 0x3FFF:
		bus.write(0x2000+addr%0x08, val)

	// NES APU and I/O registers
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
		bus.dma.TriggerOnOAM(a)
	case addr == 0x4015:
		bus.apu.writeStatus(val)
	case addr == 0x4016:
		bus.joypad.Write(val)
	case addr == 0x4017:
		bus.apu.writeFrameCounter(val)

	// APU and I/O functionality that is normally disabled.
	case 0x4018 <= addr && addr <= 0x401F:

	// Cartridge space
	case 0x4020 <= addr && addr <= 0xFFFF:
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
	// 2KB internal RAM
	case 0x0000 <= addr && addr <= 0x1FFF:
		return bus.ram[addr%0x800]

	// NES PPU registers
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

	// Mirrors of $2000-2007 (repeats every 8 bytes)
	case 0x2008 <= addr && addr <= 0x3FFF:
		return bus.Peek(0x2000 + addr%0x08)

	// NES APU and I/O registers
	case 0x4000 <= addr && addr <= 0x4017:
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

func (bus *Bus) RunDMAIfOccurred(readCycle bool) {
	if bus.dma.dcmDelay > 0 {
		bus.dma.dcmDelay--
		if bus.dma.dcmDelay == 0 {
			bus.dma.dcmDMAOccurred = true
		}
	}
	if !(bus.dma.dcmDMAOccurred || bus.dma.oamDMAOccurred) {
		return
	}

	// ref: https://www.nesdev.org/wiki/DMA#Behavior
	// > When DMA is scheduled, the associated DMA unit attempts to halt the CPU. The CPU only allows this on read cycles.
	// > If the CPU is writing, it ignores the halt and the DMA unit waits until the next cycle to try again, repeating until successful.
	if !readCycle {
		// DMA attempts to halt
		return
	}

	if bus.dma.dcmDMAOccurred || bus.dma.oamDMAOccurred {
		bus.tickStall(1) // DMA halt cycle
	}

	// TODO: implement DMC DMA during OAM DMA

	if bus.dma.dcmDMAOccurred {
		bus.dma.dcmDMAOccurred = false
		bus.tickStall(1) // DMA dummy cycle
		if bus.realClock()%2 != 0 {
			bus.tickStall(1) // DMA alignment cycle
		}
		val := bus.read(bus.dma.dcmTargetAddr)
		bus.apu.dmc.setSampleBuffer(val)
	}

	if bus.dma.oamDMAOccurred {
		bus.dma.oamDMAOccurred = false
		if bus.realClock()%2 != 0 {
			bus.tickStall(1) // DMA alignment cycle
		}
		for i := uint16(0); i < 256; i++ {
			bus.tickStall(1)
			t := bus.read(bus.dma.oamTargetAddr + i)
			bus.tickStall(1)
			bus.ppu.writeOAMDMAByte(t)
		}
	}
}
