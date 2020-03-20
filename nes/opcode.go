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

const (
	unknownInstruction = iota
	LDA
	LDX
	LDY
	STA
	STX
	STY
	TAX
	TAY
	TSX
	TXA
	TXS
	TYA
	ADC
	AND
	ASL
	BIT
	CMP
	CPX
	CPY
	DEC
	DEX
	DEY
	EOR
	INC
	INX
	INY
	LSR
	ORA
	ROL
	ROR
	SBC
	PHA
	PHP
	PLA
	PLP
	JMP
	JSR
	RTS
	RTI
	BCC
	BCS
	BEQ
	BMI
	BNE
	BPL
	BVC
	BVS
	CLC
	CLD
	CLI
	CLV
	SEC
	SED
	SEI
	BRK
	NOP
)

type opcode struct {
	Label int
	Mode  addressingMode
	Size  int
	Cycle int
}

var opcodeMap = map[byte]opcode{
	0xA9: opcode{Label: LDA, Mode: immediate, Size: 2, Cycle: 2},
	0xA5: opcode{Label: LDA, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB5: opcode{Label: LDA, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xAD: opcode{Label: LDA, Mode: absolute, Size: 3, Cycle: 4},
	0xBD: opcode{Label: LDA, Mode: absoluteX, Size: 3, Cycle: 4},
	0xB9: opcode{Label: LDA, Mode: absoluteY, Size: 3, Cycle: 4},
	0xA1: opcode{Label: LDA, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0xB1: opcode{Label: LDA, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0xA2: opcode{Label: LDX, Mode: immediate, Size: 2, Cycle: 2},
	0xA6: opcode{Label: LDX, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB6: opcode{Label: LDX, Mode: zeroPageY, Size: 2, Cycle: 4},
	0xAE: opcode{Label: LDX, Mode: absolute, Size: 3, Cycle: 4},
	0xBE: opcode{Label: LDX, Mode: absoluteY, Size: 3, Cycle: 4},
	0xA0: opcode{Label: LDY, Mode: immediate, Size: 2, Cycle: 2},
	0xA4: opcode{Label: LDY, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB4: opcode{Label: LDY, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xAC: opcode{Label: LDY, Mode: absolute, Size: 3, Cycle: 4},
	0xBC: opcode{Label: LDY, Mode: absoluteX, Size: 3, Cycle: 4},
	0x85: opcode{Label: STA, Mode: zeroPage, Size: 2, Cycle: 3},
	0x95: opcode{Label: STA, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x8D: opcode{Label: STA, Mode: absolute, Size: 3, Cycle: 4},
	0x9D: opcode{Label: STA, Mode: absoluteX, Size: 3, Cycle: 5},
	0x99: opcode{Label: STA, Mode: absoluteY, Size: 3, Cycle: 5},
	0x81: opcode{Label: STA, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x91: opcode{Label: STA, Mode: indirectIndexed, Size: 2, Cycle: 6},
	0x86: opcode{Label: STX, Mode: zeroPage, Size: 2, Cycle: 3},
	0x96: opcode{Label: STX, Mode: zeroPageY, Size: 2, Cycle: 4},
	0x8E: opcode{Label: STX, Mode: absolute, Size: 3, Cycle: 4},
	0x84: opcode{Label: STY, Mode: zeroPage, Size: 2, Cycle: 3},
	0x94: opcode{Label: STY, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x8C: opcode{Label: STY, Mode: absolute, Size: 3, Cycle: 4},
	0xAA: opcode{Label: TAX, Mode: implied, Size: 1, Cycle: 2},
	0xA8: opcode{Label: TAY, Mode: implied, Size: 1, Cycle: 2},
	0xBA: opcode{Label: TSX, Mode: implied, Size: 1, Cycle: 2},
	0x8A: opcode{Label: TXA, Mode: implied, Size: 1, Cycle: 2},
	0x9A: opcode{Label: TXS, Mode: implied, Size: 1, Cycle: 2},
	0x98: opcode{Label: TYA, Mode: implied, Size: 1, Cycle: 2},
	0x69: opcode{Label: ADC, Mode: immediate, Size: 2, Cycle: 2},
	0x65: opcode{Label: ADC, Mode: zeroPage, Size: 2, Cycle: 3},
	0x75: opcode{Label: ADC, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x6D: opcode{Label: ADC, Mode: absolute, Size: 3, Cycle: 4},
	0x7D: opcode{Label: ADC, Mode: absoluteX, Size: 3, Cycle: 4},
	0x79: opcode{Label: ADC, Mode: absoluteY, Size: 3, Cycle: 4},
	0x61: opcode{Label: ADC, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x71: opcode{Label: ADC, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x29: opcode{Label: AND, Mode: immediate, Size: 2, Cycle: 2},
	0x25: opcode{Label: AND, Mode: zeroPage, Size: 2, Cycle: 3},
	0x35: opcode{Label: AND, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x2D: opcode{Label: AND, Mode: absolute, Size: 3, Cycle: 4},
	0x3D: opcode{Label: AND, Mode: absoluteX, Size: 3, Cycle: 4},
	0x39: opcode{Label: AND, Mode: absoluteY, Size: 3, Cycle: 4},
	0x21: opcode{Label: AND, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x31: opcode{Label: AND, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x0A: opcode{Label: ASL, Mode: accumulator, Size: 1, Cycle: 2},
	0x06: opcode{Label: ASL, Mode: zeroPage, Size: 2, Cycle: 5},
	0x16: opcode{Label: ASL, Mode: zeroPageX, Size: 2, Cycle: 6},
	0x0E: opcode{Label: ASL, Mode: absolute, Size: 3, Cycle: 6},
	0x1E: opcode{Label: ASL, Mode: absoluteX, Size: 3, Cycle: 7},
	0x24: opcode{Label: BIT, Mode: zeroPage, Size: 2, Cycle: 3},
	0x2C: opcode{Label: BIT, Mode: absolute, Size: 3, Cycle: 4},
	0xC9: opcode{Label: CMP, Mode: immediate, Size: 2, Cycle: 2},
	0xC5: opcode{Label: CMP, Mode: zeroPage, Size: 2, Cycle: 3},
	0xD5: opcode{Label: CMP, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xCD: opcode{Label: CMP, Mode: absolute, Size: 3, Cycle: 4},
	0xDD: opcode{Label: CMP, Mode: absoluteX, Size: 3, Cycle: 4},
	0xD9: opcode{Label: CMP, Mode: absoluteY, Size: 3, Cycle: 4},
	0xC1: opcode{Label: CMP, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0xD1: opcode{Label: CMP, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0xE0: opcode{Label: CPX, Mode: immediate, Size: 2, Cycle: 2},
	0xE4: opcode{Label: CPX, Mode: zeroPage, Size: 2, Cycle: 3},
	0xEC: opcode{Label: CPX, Mode: absolute, Size: 3, Cycle: 4},
	0xC0: opcode{Label: CPY, Mode: immediate, Size: 2, Cycle: 2},
	0xC4: opcode{Label: CPY, Mode: zeroPage, Size: 2, Cycle: 3},
	0xCC: opcode{Label: CPY, Mode: absolute, Size: 3, Cycle: 4},
	0xC6: opcode{Label: DEC, Mode: zeroPage, Size: 2, Cycle: 5},
	0xD6: opcode{Label: DEC, Mode: zeroPageX, Size: 2, Cycle: 6},
	0xCE: opcode{Label: DEC, Mode: absolute, Size: 3, Cycle: 6},
	0xDE: opcode{Label: DEC, Mode: absoluteX, Size: 3, Cycle: 7},
	0xCA: opcode{Label: DEX, Mode: implied, Size: 1, Cycle: 2},
	0x88: opcode{Label: DEY, Mode: implied, Size: 1, Cycle: 2},
	0x49: opcode{Label: EOR, Mode: immediate, Size: 2, Cycle: 2},
	0x45: opcode{Label: EOR, Mode: zeroPage, Size: 2, Cycle: 3},
	0x55: opcode{Label: EOR, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x4D: opcode{Label: EOR, Mode: absolute, Size: 3, Cycle: 4},
	0x5D: opcode{Label: EOR, Mode: absoluteX, Size: 3, Cycle: 4},
	0x59: opcode{Label: EOR, Mode: absoluteY, Size: 3, Cycle: 4},
	0x41: opcode{Label: EOR, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x51: opcode{Label: EOR, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0xE6: opcode{Label: INC, Mode: zeroPage, Size: 2, Cycle: 5},
	0xF6: opcode{Label: INC, Mode: zeroPageX, Size: 2, Cycle: 6},
	0xEE: opcode{Label: INC, Mode: absolute, Size: 3, Cycle: 6},
	0xFE: opcode{Label: INC, Mode: absoluteX, Size: 3, Cycle: 7},
	0xE8: opcode{Label: INX, Mode: implied, Size: 1, Cycle: 2},
	0xC8: opcode{Label: INY, Mode: implied, Size: 1, Cycle: 2},
	0x4A: opcode{Label: LSR, Mode: accumulator, Size: 1, Cycle: 2},
	0x46: opcode{Label: LSR, Mode: zeroPage, Size: 2, Cycle: 5},
	0x56: opcode{Label: LSR, Mode: zeroPageX, Size: 2, Cycle: 6},
	0x4E: opcode{Label: LSR, Mode: absolute, Size: 3, Cycle: 6},
	0x5E: opcode{Label: LSR, Mode: absoluteX, Size: 3, Cycle: 7},
	0x09: opcode{Label: ORA, Mode: immediate, Size: 2, Cycle: 2},
	0x05: opcode{Label: ORA, Mode: zeroPage, Size: 2, Cycle: 3},
	0x15: opcode{Label: ORA, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x0D: opcode{Label: ORA, Mode: absolute, Size: 3, Cycle: 4},
	0x1D: opcode{Label: ORA, Mode: absoluteX, Size: 3, Cycle: 4},
	0x19: opcode{Label: ORA, Mode: absoluteY, Size: 3, Cycle: 4},
	0x01: opcode{Label: ORA, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x11: opcode{Label: ORA, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x2A: opcode{Label: ROL, Mode: accumulator, Size: 1, Cycle: 2},
	0x26: opcode{Label: ROL, Mode: zeroPage, Size: 2, Cycle: 5},
	0x36: opcode{Label: ROL, Mode: zeroPageX, Size: 2, Cycle: 6},
	0x2E: opcode{Label: ROL, Mode: absolute, Size: 3, Cycle: 6},
	0x3E: opcode{Label: ROL, Mode: absoluteX, Size: 3, Cycle: 7},
	0x6A: opcode{Label: ROR, Mode: accumulator, Size: 1, Cycle: 2},
	0x66: opcode{Label: ROR, Mode: zeroPage, Size: 2, Cycle: 5},
	0x76: opcode{Label: ROR, Mode: zeroPageX, Size: 2, Cycle: 6},
	0x6E: opcode{Label: ROR, Mode: absolute, Size: 3, Cycle: 6},
	0x7E: opcode{Label: ROR, Mode: absoluteX, Size: 3, Cycle: 7},
	0xE9: opcode{Label: SBC, Mode: immediate, Size: 2, Cycle: 2},
	0xE5: opcode{Label: SBC, Mode: zeroPage, Size: 2, Cycle: 3},
	0xF5: opcode{Label: SBC, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xED: opcode{Label: SBC, Mode: absolute, Size: 3, Cycle: 4},
	0xFD: opcode{Label: SBC, Mode: absoluteX, Size: 3, Cycle: 4},
	0xF9: opcode{Label: SBC, Mode: absoluteY, Size: 3, Cycle: 4},
	0xE1: opcode{Label: SBC, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0xF1: opcode{Label: SBC, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x48: opcode{Label: PHA, Mode: implied, Size: 1, Cycle: 3},
	0x08: opcode{Label: PHP, Mode: implied, Size: 1, Cycle: 3},
	0x68: opcode{Label: PLA, Mode: implied, Size: 1, Cycle: 4},
	0x28: opcode{Label: PLP, Mode: implied, Size: 1, Cycle: 4},
	0x4C: opcode{Label: JMP, Mode: absolute, Size: 3, Cycle: 3},
	0x6C: opcode{Label: JMP, Mode: indirect, Size: 3, Cycle: 5},
	0x20: opcode{Label: JSR, Mode: absolute, Size: 3, Cycle: 6},
	0x60: opcode{Label: RTS, Mode: implied, Size: 1, Cycle: 6},
	0x40: opcode{Label: RTI, Mode: implied, Size: 1, Cycle: 6},
	0x90: opcode{Label: BCC, Mode: relative, Size: 2, Cycle: 2},
	0xB0: opcode{Label: BCS, Mode: relative, Size: 2, Cycle: 2},
	0xF0: opcode{Label: BEQ, Mode: relative, Size: 2, Cycle: 2},
	0x30: opcode{Label: BMI, Mode: relative, Size: 2, Cycle: 2},
	0xD0: opcode{Label: BNE, Mode: relative, Size: 2, Cycle: 2},
	0x10: opcode{Label: BPL, Mode: relative, Size: 2, Cycle: 2},
	0x50: opcode{Label: BVC, Mode: relative, Size: 2, Cycle: 2},
	0x70: opcode{Label: BVS, Mode: relative, Size: 2, Cycle: 2},
	0x18: opcode{Label: CLC, Mode: implied, Size: 1, Cycle: 2},
	0xD8: opcode{Label: CLD, Mode: implied, Size: 1, Cycle: 2},
	0x58: opcode{Label: CLI, Mode: implied, Size: 1, Cycle: 2},
	0xB8: opcode{Label: CLV, Mode: implied, Size: 1, Cycle: 2},
	0x38: opcode{Label: SEC, Mode: implied, Size: 1, Cycle: 2},
	0xF8: opcode{Label: SED, Mode: implied, Size: 1, Cycle: 2},
	0x78: opcode{Label: SEI, Mode: implied, Size: 1, Cycle: 2},
	0x00: opcode{Label: BRK, Mode: implied, Size: 1, Cycle: 7},
	0xEA: opcode{Label: NOP, Mode: implied, Size: 1, Cycle: 2},
}
