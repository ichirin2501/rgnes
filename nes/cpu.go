package nes

const (
	carryFlagMask byte = (1 << iota)
	zeroFlagMask
	interruptDisableFlagMask
	decimalFlagMask
	breakFlagMask
	reservedFlagMask
	overflowFlagMask
	negativeFlagMask
)

type addressingMode int

const (
	absolute addressingMode = iota + 1
	absoluteX
	absoluteY
	accumulator
	immediate
	implied
	indexedIndirect
	indirect
	indirectIndexed
	relative
	zeroPage
	zeroPageX
	zeroPageY
)

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
	A  byte   // Accumulator
	X  byte   // Index
	Y  byte   // Index
	PC uint16 // Program Counter
	S  byte   // Stack Pointer
	P  byte   // Status Register

	cycle     *CPUCycle
	interrupt *Interrupt
	memory    Memory
	noCopy    noCopy
}

func NewCPU(mem Memory, cycle *CPUCycle, interrupt *Interrupt) *CPU {
	// ref. http://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_note-1
	return &CPU{
		A:  0x00,
		X:  0x00,
		Y:  0x00,
		PC: 0xFFFC,
		S:  0xFD,
		P:  reservedFlagMask | breakFlagMask | interruptDisableFlagMask,

		cycle:     cycle,
		interrupt: interrupt,
		memory:    mem,
	}
}

func (cpu *CPU) carryFlag() bool {
	return (cpu.P & carryFlagMask) == carryFlagMask
}
func (cpu *CPU) setCarryFlag(cond bool) {
	if cond {
		cpu.P |= carryFlagMask
	} else {
		cpu.P &= ^carryFlagMask
	}
}
func (cpu *CPU) zeroFlag() bool {
	return (cpu.P & zeroFlagMask) == zeroFlagMask
}
func (cpu *CPU) setZeroFlag(cond bool) {
	if cond {
		cpu.P |= zeroFlagMask
	} else {
		cpu.P &= ^zeroFlagMask
	}
}
func (cpu *CPU) interruptDisableFlag() bool {
	return (cpu.P & interruptDisableFlagMask) == interruptDisableFlagMask
}
func (cpu *CPU) setInterruptDisableFlag(cond bool) {
	if cond {
		cpu.P |= interruptDisableFlagMask
	} else {
		cpu.P &= ^interruptDisableFlagMask
	}
}
func (cpu *CPU) decimalFlag() bool {
	return (cpu.P & decimalFlagMask) == decimalFlagMask
}
func (cpu *CPU) setDecimalFlag(cond bool) {
	if cond {
		cpu.P |= decimalFlagMask
	} else {
		cpu.P &= ^decimalFlagMask
	}
}
func (cpu *CPU) breakFlag() bool {
	return (cpu.P & breakFlagMask) == breakFlagMask
}
func (cpu *CPU) setBreakFlag(cond bool) {
	if cond {
		cpu.P |= breakFlagMask
	} else {
		cpu.P &= ^breakFlagMask
	}
}
func (cpu *CPU) overflowFlag() bool {
	return (cpu.P & overflowFlagMask) == overflowFlagMask
}
func (cpu *CPU) setOverflowFlag(cond bool) {
	if cond {
		cpu.P |= overflowFlagMask
	} else {
		cpu.P &= ^overflowFlagMask
	}
}
func (cpu *CPU) negativeFlag() bool {
	return (cpu.P & negativeFlagMask) == negativeFlagMask
}
func (cpu *CPU) setNegativeFlag(cond bool) {
	if cond {
		cpu.P |= negativeFlagMask
	} else {
		cpu.P &= ^negativeFlagMask
	}
}
func (cpu *CPU) read16(addr uint16) uint16 {
	l := cpu.memory.Read(addr)
	h := cpu.memory.Read(addr + 1)
	return (uint16(h) << 8) | uint16(l)
}

func (cpu *CPU) push(val byte) {
	cpu.S--
	cpu.memory.Write(0x100|uint16(cpu.S), val)
}

func (cpu *CPU) push16(val uint16) {
	l := byte(val & 0xFF)
	h := byte(val >> 8)
	cpu.push(h)
	cpu.push(l)
}

func (cpu *CPU) pop() byte {
	cpu.S++
	return cpu.memory.Read(0x100 | uint16(cpu.S))
}

func (cpu *CPU) fetch() byte {
	return cpu.memory.Read(cpu.PC)
}

// TODO: after reset
func (cpu *CPU) reset() {
	cpu.PC = cpu.read16(0xFFFC)
	cpu.P = reservedFlagMask | breakFlagMask | interruptDisableFlagMask
}

func (cpu *CPU) nmi() {
	cpu.setBreakFlag(false)
	cpu.push16(cpu.PC)
	cpu.push(cpu.P)
	cpu.setInterruptDisableFlag(true)
	cpu.PC = cpu.read16(0xFFFA)
}

func (cpu *CPU) irq() {
	if cpu.interruptDisableFlag() {
		return
	}
	cpu.setBreakFlag(false)
	cpu.push16(cpu.PC)
	cpu.push(cpu.P)
	cpu.setInterruptDisableFlag(true)
	cpu.PC = cpu.read16(0xFFFE)
}

func (cpu *CPU) calcAddressing(mode addressingMode) (addr uint16, pageCrossed bool) {
	pageCrossed = false

	switch mode {
	case absolute:
		addr = cpu.read16(cpu.PC + 1)
	case absoluteX:
		addr = cpu.read16(cpu.PC+1) + uint16(cpu.X)
		pageCrossed = pagesCross(addr, addr-uint16(cpu.X))
	case absoluteY:
		addr = cpu.read16(cpu.PC+1) + uint16(cpu.Y)
		pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
	case accumulator:
		addr = 0
	case immediate:
		addr = cpu.PC + 1
	case implied:
		addr = 0
	case indexedIndirect:
		baseAddr := uint16((cpu.memory.Read(cpu.PC+1) + cpu.X) & 0xFF)
		addr = uint16(cpu.memory.Read((baseAddr+1)&0xFF))<<8 | uint16(cpu.memory.Read(baseAddr))
	case indirect:
		baseAddr := cpu.read16(cpu.PC + 1)
		addr = uint16(cpu.memory.Read((baseAddr+1)&0xFF))<<8 | uint16(cpu.memory.Read(baseAddr))
	case indirectIndexed:
		baseAddr := uint16(cpu.memory.Read(cpu.PC + 1))
		baseAddr2 := uint16(cpu.memory.Read((baseAddr+1)&0xFF))<<8 | uint16(cpu.memory.Read(baseAddr))
		addr = baseAddr2 + uint16(cpu.Y)
		pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
	case relative:
		offset := uint16(cpu.memory.Read(cpu.PC + 1))
		if offset < 0x80 {
			addr = cpu.PC + 2 + offset
		} else {
			addr = cpu.PC + 2 + offset - 0x100
		}
	case zeroPage:
		addr = uint16(cpu.memory.Read(cpu.PC + 1))
	case zeroPageX:
		addr = uint16(cpu.memory.Read(cpu.PC+1)+cpu.X) & 0xFF
	case zeroPageY:
		addr = uint16(cpu.memory.Read(cpu.PC+1)+cpu.Y) & 0xFF
	}

	return addr, pageCrossed
}

func (cpu *CPU) Step() int {
	if cpu.cycle.Stall() > 0 {
		cpu.cycle.AddStall(-1)
		return 1
	}

	if cpu.interrupt.IsNMI() {
		cpu.nmi()
		cpu.interrupt.DeassertNMI()
	} else if cpu.interrupt.IsIRQ() {
		cpu.irq()
		cpu.interrupt.DeassertIRQ()
	}

	// 色々TODO

	// opcode := cpu.fetch()

	return 0
}

func pagesCross(a uint16, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}
