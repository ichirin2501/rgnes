package nes

import (
	"fmt"
)

type CPUBus struct {
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

func NewCPUBus(ppu *PPU, apu *APU, mapper Mapper, joypad *Joypad, dma *DMA) *CPUBus {
	return &CPUBus{
		ram:    make([]byte, 2048),
		ppu:    ppu,
		apu:    apu,
		mapper: mapper,
		joypad: joypad,
		dma:    dma,
	}
}

func (bus *CPUBus) read(addr uint16) byte {
	switch {
	// 2KB internal RAM
	case 0x0000 <= addr && addr <= 0x1FFF:
		return bus.ram[addr%0x800]

	// NES PPU registers
	case 0x2000 <= addr && addr <= 0x3FFF:
		return bus.ppu.ReadRegister(addr)

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

func (bus *CPUBus) write(addr uint16, val byte) {
	switch {
	// 2KB internal RAM
	case 0x0000 <= addr && addr <= 0x1FFF:
		bus.ram[addr%0x800] = val

	// NES PPU registers
	case 0x2000 <= addr && addr <= 0x3FFF:
		bus.ppu.WriteRegister(addr, val)

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

func (bus *CPUBus) realClock() int {
	return bus.clock + bus.stall
}

func (bus *CPUBus) tick(cpuCycle int) {
	bus.clock += cpuCycle
	for i := 0; i < cpuCycle; i++ {
		bus.apu.Step()

		bus.ppu.Step()
		bus.ppu.Step()
		bus.ppu.Step()
	}
}

func (bus *CPUBus) tickStall(cpuCycle int) {
	bus.stall += cpuCycle
	for i := 0; i < cpuCycle; i++ {
		bus.apu.Step()

		bus.ppu.Step()
		bus.ppu.Step()
		bus.ppu.Step()
	}
}

// Peek is used for debugging
func (bus *CPUBus) Peek(addr uint16) byte {
	switch {
	// 2KB internal RAM
	case 0x0000 <= addr && addr <= 0x1FFF:
		return bus.ram[addr%0x800]

	// NES PPU registers
	case 0x2000 <= addr && addr <= 0x3FFF:
		return bus.ppu.PeekRegister(addr)

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

func (bus *CPUBus) RunDMAIfOccurred(readCycle bool) {
	for {
		if bus.dma.dmcDelay > 0 {
			bus.dma.dmcDelay--
			if bus.dma.dmcDelay == 0 {
				if bus.dma.dmcState == DMCDMANoneState {
					bus.dma.dmcState = DMCDMAHaltState
				}
			}
		}
		if readCycle == false {
			return
		}
		if bus.dma.oamState == OAMDMANoneState && bus.dma.dmcState == DMCDMANoneState {
			break
		}
		bus.tickStall(1)

		// https://www.nesdev.org/wiki/DMA#Behavior
		// > Get and put cycles are aligned to the first and second halves of the APU clock, respectively (called apu_clk1 and apu_clk2 in Visual2A03).
		// > While these cycles are sometimes described as even and odd CPU cycles, this is not accurate because the CPU and APU randomly power into either of 2 alignments relative to each other.
		// > Therefore, get and put may occur on different CPU cycle parities across different power cycles.
		// Since get/put cycles are randomly determined when the power is turned on, this emulator implementation fixes it as follows.
		// get = CPU even cycles
		// put = CPU odd cycles
		// And I don't know why, but it passed the cpu_interrupts_v2/4-irq_and_dma.nes test...

		// OAM
		switch bus.dma.oamState {
		case OAMDMAHaltState:
			// init
			bus.dma.oamCount = 0

			if bus.realClock()%2 == 0 { // get cycle
				bus.dma.oamState = OAMDMAAlignmentState
				bus.dma.oamSaveState = OAMDMAReadState
			} else {
				bus.dma.oamState = OAMDMAReadState
			}
		case OAMDMAAlignmentState:
			bus.dma.oamState = bus.dma.oamSaveState
		case OAMDMAReadState:
			if bus.dma.dmcState == DMCDMARunState {
				bus.dma.oamSaveState = OAMDMAReadState
				bus.dma.oamState = OAMDMAAlignmentState
			} else {
				bus.dma.oamTempByte = bus.read(bus.dma.oamTargetAddr + bus.dma.oamCount)
				bus.dma.oamState = OAMDMAWriteState
			}
		case OAMDMAWriteState:
			if bus.dma.dmcState == DMCDMARunState {
				bus.dma.oamSaveState = OAMDMAWriteState
				bus.dma.oamState = OAMDMAAlignmentState
			} else {
				bus.ppu.writeOAMData(bus.dma.oamTempByte)
				bus.dma.oamCount++
				if bus.dma.oamCount < 256 {
					bus.dma.oamState = OAMDMAReadState
				} else {
					bus.dma.oamState = OAMDMANoneState
				}
			}
		}
		// DMC
		switch bus.dma.dmcState {
		case DMCDMAHaltState:
			bus.dma.dmcState = DMCDMADummyState
		case DMCDMADummyState:
			if bus.realClock()%2 == 0 { // get cycle
				bus.dma.dmcState = DMCDMAAlignmentState
			} else {
				bus.dma.dmcState = DMCDMARunState
			}
		case DMCDMAAlignmentState:
			bus.dma.dmcState = DMCDMARunState
		case DMCDMARunState:
			val := bus.read(bus.dma.dmcTargetAddr)
			bus.apu.dmc.setSampleBuffer(val)
			bus.dma.dmcState = DMCDMANoneState
		}
	}
}
