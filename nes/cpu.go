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

	memory Memory
	noCopy noCopy
}

func NewCPU(mem Memory) *CPU {
	// ref. http://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_note-1
	return &CPU{
		A: 0x00,
		X: 0x00,
		Y: 0x00,
		//		PC: TODO: 0xFFFC
		S: 0xFD,
		P: reservedFlagMask | breakFlagMask | interruptDisableFlagMask,

		memory: mem,
	}
}

func (cpu *CPU) CarryFlag() bool {
	return (cpu.P & carryFlagMask) == carryFlagMask
}

func (cpu *CPU) ZeroFlag() bool {
	return (cpu.P & zeroFlagMask) == zeroFlagMask
}

func (cpu *CPU) InterruptDisableFlag() bool {
	return (cpu.P & interruptDisableFlagMask) == interruptDisableFlagMask
}

func (cpu *CPU) DecimalFlag() bool {
	return (cpu.P & decimalFlagMask) == decimalFlagMask
}

func (cpu *CPU) BreakFlag() bool {
	return (cpu.P & breakFlagMask) == breakFlagMask
}

func (cpu *CPU) OverflowFlag() bool {
	return (cpu.P & overflowFlagMask) == overflowFlagMask
}

func (cpu *CPU) NegativeFlag() bool {
	return (cpu.P & negativeFlagMask) == negativeFlagMask
}

func (cpu *CPU) reset() {
	l := cpu.memory.Read(0xFFFC)
	h := cpu.memory.Read(0xFFFD)
	cpu.PC = (uint16(h) << 8) | uint16(l)
	cpu.P = reservedFlagMask | breakFlagMask | interruptDisableFlagMask
}
