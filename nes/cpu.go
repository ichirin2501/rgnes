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

type CPU struct {
	A  byte   // Accumulator
	X  byte   // Index
	Y  byte   // Index
	PC uint16 // Program Counter
	S  byte   // Stack Pointer
	P  byte   // Status Register

	interrupt *Interrupt
	memory    Memory
	noCopy    noCopy
}

func NewCPU(mem Memory, interrupt *Interrupt) *CPU {
	// ref. http://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_note-1
	return &CPU{
		A: 0x00,
		X: 0x00,
		Y: 0x00,
		//		PC: TODO: 0xFFFC
		S: 0xFD,
		P: reservedFlagMask | breakFlagMask | interruptDisableFlagMask,

		interrupt: interrupt,
		memory:    mem,
	}
}

func (cpu *CPU) carryFlag() bool {
	return (cpu.P & carryFlagMask) == carryFlagMask
}
func (cpu *CPU) setCarryFlag() {
	cpu.P |= carryFlagMask
}
func (cpu *CPU) unsetCarryFlag() {
	cpu.P &= ^carryFlagMask
}
func (cpu *CPU) zeroFlag() bool {
	return (cpu.P & zeroFlagMask) == zeroFlagMask
}
func (cpu *CPU) setZeroFlag() {
	cpu.P |= zeroFlagMask
}
func (cpu *CPU) unsetZeroFlag() {
	cpu.P &= ^zeroFlagMask
}
func (cpu *CPU) interruptDisableFlag() bool {
	return (cpu.P & interruptDisableFlagMask) == interruptDisableFlagMask
}
func (cpu *CPU) setInterruptDisableFlag() {
	cpu.P |= interruptDisableFlagMask
}
func (cpu *CPU) unsetInterruptDisableFlag() {
	cpu.P &= ^interruptDisableFlagMask
}
func (cpu *CPU) decimalFlag() bool {
	return (cpu.P & decimalFlagMask) == decimalFlagMask
}
func (cpu *CPU) setDecimalFlag() {
	cpu.P |= decimalFlagMask
}
func (cpu *CPU) unsetDecimalFlag() {
	cpu.P &= ^decimalFlagMask
}
func (cpu *CPU) breakFlag() bool {
	return (cpu.P & breakFlagMask) == breakFlagMask
}
func (cpu *CPU) setBreakFlag() {
	cpu.P |= breakFlagMask
}
func (cpu *CPU) unsetBreakFlag() {
	cpu.P &= ^breakFlagMask
}
func (cpu *CPU) overflowFlag() bool {
	return (cpu.P & overflowFlagMask) == overflowFlagMask
}
func (cpu *CPU) setOverflowFlag() {
	cpu.P |= overflowFlagMask
}
func (cpu *CPU) unsetOverflowFlag() {
	cpu.P &= ^overflowFlagMask
}
func (cpu *CPU) negativeFlag() bool {
	return (cpu.P & negativeFlagMask) == negativeFlagMask
}
func (cpu *CPU) setNegativeFlag() {
	cpu.P |= negativeFlagMask
}
func (cpu *CPU) unsetNegativeFlag() {
	cpu.P &= ^negativeFlagMask
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

func (cpu *CPU) reset() {
	cpu.PC = cpu.read16(0xFFFC)
	cpu.P = reservedFlagMask | breakFlagMask | interruptDisableFlagMask
}

func (cpu *CPU) nmi() {
	cpu.unsetBreakFlag()
	cpu.push16(cpu.PC)
	cpu.push(cpu.P)
	cpu.setInterruptDisableFlag()
	cpu.PC = cpu.read16(0xFFFA)
}

func (cpu *CPU) irq() {
	if cpu.interruptDisableFlag() {
		return
	}
	cpu.unsetBreakFlag()
	cpu.push16(cpu.PC)
	cpu.push(cpu.P)
	cpu.setInterruptDisableFlag()
	cpu.PC = cpu.read16(0xFFFE)
}
