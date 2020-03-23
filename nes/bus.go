package nes

type CPUBus struct {
	cycle  *CPUCycle
	ram    Memory
	prg    MemoryReader
	noCopy noCopy
}

func NewCPUBus(cycle *CPUCycle, ram Memory, prg MemoryReader) *CPUBus {
	return &CPUBus{
		cycle: cycle,
		ram:   ram,
		prg:   prg,
	}
}

func (bus *CPUBus) Read(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return bus.ram.Read(addr % 0x800)
	case addr >= 0x6000:
		return bus.prg.Read(addr - 0xC000)
	default:
		panic("unimplemented")
	}
}

func (bus *CPUBus) Write(addr uint16, val byte) {
	switch {
	case addr < 0x2000:
		bus.ram.Write(addr%0x800, val)
	default:
		panic("unimplemented")
	}
}
