package bus

import "github.com/ichirin2501/rgnes/nes/memory"

type PPUBus struct {
	ram    memory.Memory
	mapper memory.Memory
}

func NewPPUBus(ram memory.Memory, mapper memory.Memory) *PPUBus {
	return &PPUBus{
		ram:    ram,
		mapper: mapper,
	}
}

func (bus *PPUBus) Read(addr uint16) byte {
	switch {
	case addr < 0x2000: // Pattern table
		return bus.mapper.Read(addr)
	case addr < 0x3000: // Nametable
		return bus.ram.Read(addr - 0x2000)
	case addr < 0x3EFF: // Mirrors of $2000-$2EFF
		return bus.ram.Read(addr - 0x1000 - 0x2000)
	case addr < 0x3F20: // Palette RAM indexes
		return bus.ram.Read(addr - 0x2000)
	case addr < 0x4000: // Mirrors of $3F00-$3F1F
		return bus.ram.Read(addr - 0x20 - 0x2000)
	}
	panic("Unable to reach here")
}

func (bus *PPUBus) Write(addr uint16, val byte) {
	switch {
	case addr < 0x2000:
		bus.mapper.Write(addr, val)
	case addr < 0x3000: // Nametable
		bus.ram.Write(addr-0x2000, val)
	case addr < 0x3EFF: // Mirrors of $2000-$2EFF
		bus.ram.Write(addr-0x1000-0x2000, val)
	case addr < 0x3F20: // Palette RAM indexes
		bus.ram.Write(addr-0x2000, val)
	case addr < 0x4000: // Mirrors of $3F00-$3F1F
		bus.ram.Write(addr-0x20-0x2000, val)
	}
}
