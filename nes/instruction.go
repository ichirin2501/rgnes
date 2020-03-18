package nes

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

type instruction int

const (
	UNKNOWN instruction = iota
	LDA
)

type opcode struct {
	Name  instruction
	Mode  addressingMode
	Size  int
	Cycle int
}

var m = map[byte]opcode{
	0xA9: opcode{Name: LDA, Mode: immediate, Size: 2, Cycle: 2},
	0xA5: opcode{Name: LDA, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB5: opcode{Name: LDA, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xAD: opcode{Name: LDA, Mode: absolute, Size: 3, Cycle: 4},
	0xBD: opcode{Name: LDA, Mode: absoluteX, Size: 3, Cycle: 4},
	0xB9: opcode{Name: LDA, Mode: absoluteY, Size: 3, Cycle: 4},
	0xA1: opcode{Name: LDA, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0xB1: opcode{Name: LDA, Mode: indirectIndexed, Size: 2, Cycle: 5},
}
