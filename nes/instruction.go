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
	ldx
	ldy
	sta
	stx
	sty
	tax
	tay
	tsx
	txa
	txs
	tya
	adc
	and
	asl
	bit
	cmp
	cpx
	cpy
	dec
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
	0xA2: opcode{Name: ldx, Mode: immediate, Size: 2, Cycle: 2},
	0xA6: opcode{Name: ldx, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB6: opcode{Name: ldx, Mode: zeroPageY, Size: 2, Cycle: 4},
	0xAE: opcode{Name: ldx, Mode: absolute, Size: 3, Cycle: 4},
	0xBE: opcode{Name: ldx, Mode: absoluteY, Size: 3, Cycle: 4},
	0xA0: opcode{Name: ldy, Mode: immediate, Size: 2, Cycle: 2},
	0xA4: opcode{Name: ldy, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB4: opcode{Name: ldy, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xAC: opcode{Name: ldy, Mode: absolute, Size: 3, Cycle: 4},
	0xBC: opcode{Name: ldy, Mode: absoluteX, Size: 3, Cycle: 4},
	0x85: opcode{Name: sta, Mode: zeroPage, Size: 2, Cycle: 3},
	0x95: opcode{Name: sta, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x8D: opcode{Name: sta, Mode: absolute, Size: 3, Cycle: 4},
	0x9D: opcode{Name: sta, Mode: absoluteX, Size: 3, Cycle: 5},
	0x99: opcode{Name: sta, Mode: absoluteY, Size: 3, Cycle: 5},
	0x81: opcode{Name: sta, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x91: opcode{Name: sta, Mode: indirectIndexed, Size: 2, Cycle: 6},
	0x86: opcode{Name: stx, Mode: zeroPage, Size: 2, Cycle: 3},
	0x96: opcode{Name: stx, Mode: zeroPageY, Size: 2, Cycle: 4},
	0x8E: opcode{Name: stx, Mode: absolute, Size: 3, Cycle: 4},
	0x84: opcode{Name: sty, Mode: zeroPage, Size: 2, Cycle: 3},
	0x94: opcode{Name: sty, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x8C: opcode{Name: sty, Mode: absolute, Size: 3, Cycle: 4},
	0xAA: opcode{Name: tax, Mode: implied, Size: 1, Cycle: 2},
	0xA8: opcode{Name: tay, Mode: implied, Size: 1, Cycle: 2},
	0xBA: opcode{Name: tsx, Mode: implied, Size: 1, Cycle: 2},
	0x8A: opcode{Name: txa, Mode: implied, Size: 1, Cycle: 2},
	0x9A: opcode{Name: txs, Mode: implied, Size: 1, Cycle: 2},
	0x98: opcode{Name: tya, Mode: implied, Size: 1, Cycle: 2},
	0x69: opcode{Name: adc, Mode: immediate, Size: 2, Cycle: 2},
	0x65: opcode{Name: adc, Mode: zeroPage, Size: 2, Cycle: 3},
	0x75: opcode{Name: adc, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x6D: opcode{Name: adc, Mode: absolute, Size: 3, Cycle: 4},
	0x7D: opcode{Name: adc, Mode: absoluteX, Size: 3, Cycle: 4},
	0x79: opcode{Name: adc, Mode: absoluteY, Size: 3, Cycle: 4},
	0x61: opcode{Name: adc, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x71: opcode{Name: adc, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x29: opcode{Name: and, Mode: immediate, Size: 2, Cycle: 2},
	0x25: opcode{Name: and, Mode: zeroPage, Size: 2, Cycle: 3},
	0x35: opcode{Name: and, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x2D: opcode{Name: and, Mode: absolute, Size: 3, Cycle: 4},
	0x3D: opcode{Name: and, Mode: absoluteX, Size: 3, Cycle: 4},
	0x39: opcode{Name: and, Mode: absoluteY, Size: 3, Cycle: 4},
	0x21: opcode{Name: and, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x31: opcode{Name: and, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x0A: opcode{Name: asl, Mode: accumulator, Size: 1, Cycle: 2},
	0x06: opcode{Name: asl, Mode: zeroPage, Size: 2, Cycle: 5},
	0x16: opcode{Name: asl, Mode: zeroPageX, Size: 2, Cycle: 6},
	0x0E: opcode{Name: asl, Mode: absolute, Size: 3, Cycle: 6},
	0x1E: opcode{Name: asl, Mode: absoluteX, Size: 3, Cycle: 7},
	0x24: opcode{Name: bit, Mode: zeroPage, Size: 2, Cycle: 3},
	0x2C: opcode{Name: bit, Mode: absolute, Size: 3, Cycle: 4},
	0xC9: opcode{Name: cmp, Mode: immediate, Size: 2, Cycle: 2},
	0xC5: opcode{Name: cmp, Mode: zeroPage, Size: 2, Cycle: 3},
	0xD5: opcode{Name: cmp, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xCD: opcode{Name: cmp, Mode: absolute, Size: 3, Cycle: 4},
	0xDD: opcode{Name: cmp, Mode: absoluteX, Size: 3, Cycle: 4},
	0xD9: opcode{Name: cmp, Mode: absoluteY, Size: 3, Cycle: 4},
	0xC1: opcode{Name: cmp, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0xD1: opcode{Name: cmp, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0xE0: opcode{Name: cpx, Mode: immediate, Size: 2, Cycle: 2},
	0xE4: opcode{Name: cpx, Mode: zeroPage, Size: 2, Cycle: 3},
	0xEC: opcode{Name: cpx, Mode: absolute, Size: 3, Cycle: 4},
	0xC0: opcode{Name: cpy, Mode: immediate, Size: 2, Cycle: 2},
	0xC4: opcode{Name: cpy, Mode: zeroPage, Size: 2, Cycle: 3},
	0xCC: opcode{Name: cpy, Mode: absolute, Size: 3, Cycle: 4},
}
