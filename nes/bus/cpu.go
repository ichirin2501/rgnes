package bus

import "github.com/ichirin2501/rgnes/nes/memory"

type CPUBus struct {
	ram    memory.Memory
	ppu    memory.Memory
	apu    memory.Memory
	mapper memory.Memory
}

func NewCPUBus(ram memory.Memory, ppu memory.Memory, apu memory.Memory, mapper memory.Memory) *CPUBus {
	return &CPUBus{
		ram:    ram,
		ppu:    ppu,
		apu:    apu,
		mapper: mapper,
	}
}

func (bus *CPUBus) Read(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return bus.ram.Read(addr % 0x800)
	case addr < 0x4000:
		return bus.ppu.Read(addr % 0x08)
	case addr == 0x4016: // TODO: keypad
	case addr == 0x4017: // TODO: 2p
	case addr < 0x4020:
		return bus.apu.Read(addr % 0x20)
	case addr >= 0x4020:
		return bus.mapper.Read(addr)
	}
	panic("Unable to reach here")
}

func (bus *CPUBus) Write(addr uint16, val byte) {
	switch {
	case addr < 0x2000:
		bus.ram.Write(addr%0x800, val)
	case addr < 0x4000:
		bus.ppu.Write(addr%0x08, val)
	case addr == 0x4014: // TODO: DMA
	case addr == 0x4016: // TODO: keypad
	case addr < 0x4020:
		bus.apu.Write(addr%0x20, val)
	case addr >= 0x4020:
		bus.mapper.Write(addr, val)
	}
}
