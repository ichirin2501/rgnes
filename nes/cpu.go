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
	r         *cpuRegister
	cycle     *CPUCycle
	interrupt *Interrupt
	memory    Memory
	noCopy    noCopy
}

func NewCPU(mem Memory, cycle *CPUCycle, interrupt *Interrupt) *CPU {
	// ref. http://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_note-1
	return &CPU{
		r:         newCPURegister(),
		cycle:     cycle,
		interrupt: interrupt,
		memory:    mem,
	}
}

// TODO: after reset
func (cpu *CPU) reset() {
	cpu.r.PC = read16(cpu.memory, 0xFFFC)
	cpu.r.P = reservedFlagMask | breakFlagMask | interruptDisableFlagMask
}

func (cpu *CPU) Step() int {
	if cpu.cycle.Stall() > 0 {
		cpu.cycle.AddStall(-1)
		return 1
	}

	if cpu.interrupt.IsNMI() {
		nmi(cpu.r, cpu.memory)
		cpu.interrupt.DeassertNMI()
	} else if cpu.interrupt.IsIRQ() {
		irq(cpu.r, cpu.memory)
		cpu.interrupt.DeassertIRQ()
	}

	// 色々TODO

	// opcodeByte := fetch(cpu.r, cpu.memory)
	// opcode := opcodeMap[opcodeByte]

	// addr, pageCrossed := fetchOperand(cpu.r, cpu.memory, opcode.Mode)

	// execute(opcode)

	return 0
}

func fetch(r *cpuRegister, m MemoryReader) byte {
	v := m.Read(r.PC)
	r.PC++
	return v
}

func fetch16(r *cpuRegister, m MemoryReader) uint16 {
	l := fetch(r, m)
	h := fetch(r, m)
	return uint16(h)<<8 | uint16(l)
}

func fetchOperand(r *cpuRegister, m MemoryReader, mode addressingMode) (addr uint16, pageCrossed bool) {
	pageCrossed = false

	switch mode {
	case absolute:
		addr = fetch16(r, m)
	case absoluteX:
		addr = fetch16(r, m) + uint16(r.X)
		pageCrossed = pagesCross(addr, addr-uint16(r.X))
	case absoluteY:
		addr = fetch16(r, m) + uint16(r.Y)
		pageCrossed = pagesCross(addr, addr-uint16(r.Y))
	case accumulator:
		addr = 0
	case immediate:
		addr = r.PC
		r.PC++
	case implied:
		addr = 0
	case indexedIndirect:
		baseAddr := uint16((fetch(r, m) + r.X) & 0xFF)
		addr = uint16(m.Read((baseAddr+1)&0xFF))<<8 | uint16(m.Read(baseAddr))
	case indirect:
		baseAddr := fetch16(r, m)
		addr = uint16(m.Read((baseAddr+1)&0xFF))<<8 | uint16(m.Read(baseAddr))
	case indirectIndexed:
		baseAddr := uint16(fetch(r, m))
		baseAddr2 := uint16(m.Read((baseAddr+1)&0xFF))<<8 | uint16(m.Read(baseAddr))
		addr = baseAddr2 + uint16(r.Y)
		pageCrossed = pagesCross(addr, addr-uint16(r.Y))
	case relative:
		offset := uint16(fetch(r, m))
		if offset < 0x80 {
			addr = r.PC + 2 + offset
		} else {
			addr = r.PC + 2 + offset - 0x100
		}
	case zeroPage:
		addr = uint16(fetch(r, m))
	case zeroPageX:
		addr = uint16(fetch(r, m)+r.X) & 0xFF
	case zeroPageY:
		addr = uint16(fetch(r, m)+r.Y) & 0xFF
	}

	return addr, pageCrossed
}
