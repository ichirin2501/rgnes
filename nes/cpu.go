package nes

type CPU struct {
	A  byte   // Accumulator
	X  byte   // Index
	Y  byte   // Index
	PC uint16 // Program Counter
	S  byte   // Stack Pointer
	P  byte   // Status Register

	noCopy noCopy
}

func NewCPU() *CPU {
	// ref. http://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_note-1
	return &CPU{
		A: 0x00,
		X: 0x00,
		Y: 0x00,
		//		PC: TODO: 0xFFFC
		S: 0xFD,
		P: 0x34,
	}
}
