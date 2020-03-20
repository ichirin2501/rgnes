package nes

type CPUCycle struct {
	stall  int
	cycles int

	noCopy noCopy
}

func NewCPUCycle() *CPUCycle {
	return &CPUCycle{}
}
func (c *CPUCycle) Stall() int {
	return c.stall
}
func (c *CPUCycle) Cycles() int {
	return c.cycles
}
func (c *CPUCycle) AddStall(x int) int {
	c.stall += x
	return c.stall
}
func (c *CPUCycle) AddCycles(x int) int {
	c.cycles += x
	return c.cycles
}

type CPU struct {
	reg       cpuRegisterer
	cycle     *CPUCycle
	interrupt *Interrupt
	memory    Memory
	noCopy    noCopy
}

func NewCPU(mem Memory, cycle *CPUCycle, interrupt *Interrupt) *CPU {
	// ref. http://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_note-1
	return &CPU{
		reg:       newCPURegister(),
		cycle:     cycle,
		interrupt: interrupt,
		memory:    mem,
	}
}

func (cpu *CPU) read16(addr uint16) uint16 {
	return read16(cpu.memory, addr)
}

func (cpu *CPU) push(val byte) {
	push(cpu.reg, cpu.memory, val)
}

func (cpu *CPU) push16(val uint16) {
	l := byte(val & 0xFF)
	h := byte(val >> 8)
	cpu.push(h)
	cpu.push(l)
}

func (cpu *CPU) pop() byte {
	return pop(cpu.reg, cpu.memory)
}

// TODO: after reset
func (cpu *CPU) reset() {
	cpu.reg.SetPC(cpu.read16(0xFFFC))
	cpu.reg.SetP(reservedFlagMask | breakFlagMask | interruptDisableFlagMask)
}

func (cpu *CPU) nmi() {
	nmi(cpu.reg, cpu.memory)
}

// func (cpu *CPU) irq() {
// 	if cpu.interruptDisableFlag() {
// 		return
// 	}
// 	cpu.setBreakFlag(false)
// 	cpu.push16(cpu.PC)
// 	cpu.push(cpu.P)
// 	cpu.setInterruptDisableFlag(true)
// 	cpu.PC = cpu.read16(0xFFFE)
// }

// func (cpu *CPU) calcAddressing(mode addressingMode) (addr uint16, pageCrossed bool) {
// 	pageCrossed = false

// 	switch mode {
// 	case absolute:
// 		addr = cpu.read16(cpu.PC + 1)
// 	case absoluteX:
// 		addr = cpu.read16(cpu.PC+1) + uint16(cpu.X)
// 		pageCrossed = pagesCross(addr, addr-uint16(cpu.X))
// 	case absoluteY:
// 		addr = cpu.read16(cpu.PC+1) + uint16(cpu.Y)
// 		pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
// 	case accumulator:
// 		addr = 0
// 	case immediate:
// 		addr = cpu.PC + 1
// 	case implied:
// 		addr = 0
// 	case indexedIndirect:
// 		baseAddr := uint16((cpu.memory.Read(cpu.PC+1) + cpu.X) & 0xFF)
// 		addr = uint16(cpu.memory.Read((baseAddr+1)&0xFF))<<8 | uint16(cpu.memory.Read(baseAddr))
// 	case indirect:
// 		baseAddr := cpu.read16(cpu.PC + 1)
// 		addr = uint16(cpu.memory.Read((baseAddr+1)&0xFF))<<8 | uint16(cpu.memory.Read(baseAddr))
// 	case indirectIndexed:
// 		baseAddr := uint16(cpu.memory.Read(cpu.PC + 1))
// 		baseAddr2 := uint16(cpu.memory.Read((baseAddr+1)&0xFF))<<8 | uint16(cpu.memory.Read(baseAddr))
// 		addr = baseAddr2 + uint16(cpu.Y)
// 		pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
// 	case relative:
// 		offset := uint16(cpu.memory.Read(cpu.PC + 1))
// 		if offset < 0x80 {
// 			addr = cpu.PC + 2 + offset
// 		} else {
// 			addr = cpu.PC + 2 + offset - 0x100
// 		}
// 	case zeroPage:
// 		addr = uint16(cpu.memory.Read(cpu.PC + 1))
// 	case zeroPageX:
// 		addr = uint16(cpu.memory.Read(cpu.PC+1)+cpu.X) & 0xFF
// 	case zeroPageY:
// 		addr = uint16(cpu.memory.Read(cpu.PC+1)+cpu.Y) & 0xFF
// 	}

// 	return addr, pageCrossed
// }

// func (cpu *CPU) Step() int {
// 	if cpu.cycle.Stall() > 0 {
// 		cpu.cycle.AddStall(-1)
// 		return 1
// 	}

// 	if cpu.interrupt.IsNMI() {
// 		cpu.nmi()
// 		cpu.interrupt.DeassertNMI()
// 	} else if cpu.interrupt.IsIRQ() {
// 		cpu.irq()
// 		cpu.interrupt.DeassertIRQ()
// 	}

// 	// 色々TODO

// 	// opcode := cpu.fetch()

// 	return 0
// }

func pagesCross(a uint16, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

func nmi(reg cpuRegisterer, memory Memory) {
	reg.SetBreakFlag(false)
	push16(reg, memory, reg.PC())
	push(reg, memory, reg.P())
	reg.SetInterruptDisableFlag(true)
	reg.SetPC(read16(memory, 0xFFFA))
}

func push(reg registerer, memory MemoryWriter, val byte) {
	memory.Write(0x100|uint16(reg.S()), val)
	reg.SetS(reg.S() - 1)
}

func push16(reg registerer, memory MemoryWriter, val uint16) {
	l := byte(val & 0xFF)
	h := byte(val >> 8)
	push(reg, memory, h)
	push(reg, memory, l)
}

func pop(reg registerer, memory MemoryReader) byte {
	reg.SetS(reg.S() + 1)
	return memory.Read(0x100 | uint16(reg.S()))
}

func read16(memory MemoryReader, addr uint16) uint16 {
	l := memory.Read(addr)
	h := memory.Read(addr + 1)
	return (uint16(h) << 8) | uint16(l)
}
