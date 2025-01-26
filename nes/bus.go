package nes

import (
	"fmt"
)

type cpuBus struct {
	ram    []byte
	ppu    *ppu
	apu    *apu
	mapper Mapper
	joypad *joypad
	dma    *dma

	// This clock is used to adjust the clock difference for each instruction,
	// so keep the state separate from the $4014 dma stall.
	clock int
	stall int
}

func newCPUBus(ppu *ppu, apu *apu, mapper Mapper, joypad *joypad, dma *dma) *cpuBus {
	return &cpuBus{
		ram:    make([]byte, 2048),
		ppu:    ppu,
		apu:    apu,
		mapper: mapper,
		joypad: joypad,
		dma:    dma,
	}
}

func (bus *cpuBus) read(addr uint16) byte {
	switch {
	// 2KB internal RAM
	case 0x0000 <= addr && addr <= 0x1FFF:
		return bus.ram[addr%0x800]

	// NES PPU registers
	case 0x2000 <= addr && addr <= 0x3FFF:
		return bus.ppu.readRegister(addr)

	// NES APU and I/O registers
	case 0x4000 <= addr && addr <= 0x4017:
		switch {
		case addr == 0x4015:
			return bus.apu.readStatus()
		case addr == 0x4016:
			return bus.joypad.read()
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

func (bus *cpuBus) write(addr uint16, val byte) {
	switch {
	// 2KB internal RAM
	case 0x0000 <= addr && addr <= 0x1FFF:
		bus.ram[addr%0x800] = val

	// NES PPU registers
	case 0x2000 <= addr && addr <= 0x3FFF:
		bus.ppu.writeRegister(addr, val)

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
		bus.dma.triggerOnOAM(a)
	case addr == 0x4015:
		bus.apu.writeStatus(val)
	case addr == 0x4016:
		bus.joypad.write(val)
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

func (bus *cpuBus) realClock() int {
	return bus.clock + bus.stall
}

func (bus *cpuBus) tick(cpuCycle int) {
	bus.clock += cpuCycle
	for i := 0; i < cpuCycle; i++ {
		bus.apu.step()

		bus.ppu.step()
		bus.ppu.step()
		bus.ppu.step()
	}
}

func (bus *cpuBus) tickStall(cpuCycle int) {
	bus.stall += cpuCycle
	for i := 0; i < cpuCycle; i++ {
		bus.apu.step()

		bus.ppu.step()
		bus.ppu.step()
		bus.ppu.step()
	}
}

// peek is used for debugging
func (bus *cpuBus) peek(addr uint16) byte {
	switch {
	// 2KB internal RAM
	case 0x0000 <= addr && addr <= 0x1FFF:
		return bus.ram[addr%0x800]

	// NES PPU registers
	case 0x2000 <= addr && addr <= 0x3FFF:
		return bus.ppu.peekRegister(addr)

	// NES APU and I/O registers
	case 0x4000 <= addr && addr <= 0x4017:
		// todo
		switch {
		case addr == 0x4015:
			return bus.apu.peekStatus()
		case addr == 0x4016:
			return bus.joypad.peek()
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

func (bus *cpuBus) runDMAIfOccurred(readCycle bool) {
	for {
		if bus.dma.dmcDelay > 0 {
			bus.dma.dmcDelay--
			if bus.dma.dmcDelay == 0 {
				if bus.dma.dmcState == dmcDMANoneState {
					bus.dma.dmcState = dmcDMAHaltState
				}
			}
		}
		if readCycle == false {
			return
		}
		if bus.dma.oamState == oamDMANoneState && bus.dma.dmcState == dmcDMANoneState {
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
		case oamDMAHaltState:
			// init
			bus.dma.oamCount = 0

			if bus.realClock()%2 == 0 { // get cycle
				bus.dma.oamState = oamDMAAlignmentState
				bus.dma.oamSaveState = oamDMAReadState
			} else {
				bus.dma.oamState = oamDMAReadState
			}
		case oamDMAAlignmentState:
			bus.dma.oamState = bus.dma.oamSaveState
		case oamDMAReadState:
			if bus.dma.dmcState == dmcDMARunState {
				bus.dma.oamSaveState = oamDMAReadState
				bus.dma.oamState = oamDMAAlignmentState
			} else {
				bus.dma.oamTempByte = bus.read(bus.dma.oamTargetAddr + bus.dma.oamCount)
				bus.dma.oamState = oamDMAWriteState
			}
		case oamDMAWriteState:
			if bus.dma.dmcState == dmcDMARunState {
				bus.dma.oamSaveState = oamDMAWriteState
				bus.dma.oamState = oamDMAAlignmentState
			} else {
				bus.ppu.writeOAMData(bus.dma.oamTempByte)
				bus.dma.oamCount++
				if bus.dma.oamCount < 256 {
					bus.dma.oamState = oamDMAReadState
				} else {
					bus.dma.oamState = oamDMANoneState
				}
			}
		}
		// DMC
		switch bus.dma.dmcState {
		case dmcDMAHaltState:
			bus.dma.dmcState = dmcDMADummyState
		case dmcDMADummyState:
			if bus.realClock()%2 == 0 { // get cycle
				bus.dma.dmcState = dmcDMAAlignmentState
			} else {
				bus.dma.dmcState = dmcDMARunState
			}
		case dmcDMAAlignmentState:
			bus.dma.dmcState = dmcDMARunState
		case dmcDMARunState:
			val := bus.read(bus.dma.dmcTargetAddr)
			bus.apu.dmc.setSampleBuffer(val)
			bus.dma.dmcState = dmcDMANoneState
		}
	}
}
