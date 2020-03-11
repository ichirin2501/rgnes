package nes

const (
	carryFlagMask byte = (1 << iota)
	zeroFlagMask
	interruptDisableFlagMask
	decimalFlagMask
	breakFlagMask
	_
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
		P: 0x34,

		memory: mem,
	}
}

func (cpu *CPU) SetFlags(f byte) {
	cpu.P |= f
}
