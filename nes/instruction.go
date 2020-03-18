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
	unknownInstruction instruction = iota
	lda
)

type opcode struct {
	Name  instruction
	Mode  addressingMode
	Size  int
	Cycle int
}

var m = map[byte]opcode{
	0xA9: opcode{Name: lda, Mode: immediate, Size: 2, Cycle: 2},
	0xA5: opcode{Name: lda, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB5: opcode{Name: lda, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xAD: opcode{Name: lda, Mode: absolute, Size: 3, Cycle: 4},
	0xBD: opcode{Name: lda, Mode: absoluteX, Size: 3, Cycle: 4},
	0xB9: opcode{Name: lda, Mode: absoluteY, Size: 3, Cycle: 4},
	0xA1: opcode{Name: lda, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0xB1: opcode{Name: lda, Mode: indirectIndexed, Size: 2, Cycle: 5},
}
