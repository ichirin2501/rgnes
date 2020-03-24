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

func (a addressingMode) String() string {
	switch a {
	case absoluteX:
		return "absoluteX"
	case absoluteY:
		return "absoluteY"
	case accumulator:
		return "accumulator"
	case immediate:
		return "immediate"
	case implied:
		return "implied"
	case indexedIndirect:
		return "indexedIndirect"
	case indirect:
		return "indirect"
	case indirectIndexed:
		return "indirectIndexed"
	case relative:
		return "relative"
	case zeroPage:
		return "zeroPage"
	case zeroPageX:
		return "zeroPageX"
	case zeroPageY:
		return "zeroPageY"
	default:
		panic("Unable to reach here")
	}
}

type Mnemonic int

const (
	UnknownMnemonic Mnemonic = iota
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

func (m Mnemonic) String() string {
	switch m {
	case UnknownMnemonic:
		return "UnknownMnemonic"
	case LDA:
		return "LDA"
	case LDX:
		return "LDX"
	case LDY:
		return "LDY"
	case STA:
		return "STA"
	case STX:
		return "STX"
	case STY:
		return "STY"
	case TAX:
		return "TAX"
	case TAY:
		return "TAY"
	case TSX:
		return "TSX"
	case TXA:
		return "TXA"
	case TXS:
		return "TXS"
	case TYA:
		return "TYA"
	case ADC:
		return "ADC"
	case AND:
		return "AND"
	case ASL:
		return "ASL"
	case BIT:
		return "BIT"
	case CMP:
		return "CMP"
	case CPX:
		return "CPX"
	case CPY:
		return "CPY"
	case DEC:
		return "DEC"
	case DEX:
		return "DEX"
	case DEY:
		return "DEY"
	case EOR:
		return "EOR"
	case INC:
		return "INC"
	case INX:
		return "INX"
	case INY:
		return "INY"
	case LSR:
		return "LSR"
	case ORA:
		return "ORA"
	case ROL:
		return "ROL"
	case ROR:
		return "ROR"
	case SBC:
		return "SBC"
	case PHA:
		return "PHA"
	case PHP:
		return "PHP"
	case PLA:
		return "PLA"
	case PLP:
		return "PLP"
	case JMP:
		return "JMP"
	case JSR:
		return "JSR"
	case RTS:
		return "RTS"
	case RTI:
		return "RTI"
	case BCC:
		return "BCC"
	case BCS:
		return "BCS"
	case BEQ:
		return "BEQ"
	case BMI:
		return "BMI"
	case BNE:
		return "BNE"
	case BPL:
		return "BPL"
	case BVC:
		return "BVC"
	case BVS:
		return "BVS"
	case CLC:
		return "CLC"
	case CLD:
		return "CLD"
	case CLI:
		return "CLI"
	case CLV:
		return "CLV"
	case SEC:
		return "SEC"
	case SED:
		return "SED"
	case SEI:
		return "SEI"
	case BRK:
		return "BRK"
	case NOP:
		return "NOP"
	default:
		panic("Unable to reach here")
	}
}

type opcode struct {
	Name      Mnemonic
	Mode      addressingMode
	Cycle     int
	PageCycle int
}

var opcodeMap = []*opcode{
	/* 0x00 */ &opcode{Name: BRK, Mode: implied, Cycle: 7, PageCycle: 0},
	/* 0x01 */ &opcode{Name: ORA, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x02 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x03 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x04 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x05 */ &opcode{Name: ORA, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x06 */ &opcode{Name: ASL, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x07 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x08 */ &opcode{Name: PHP, Mode: implied, Cycle: 3, PageCycle: 0},
	/* 0x09 */ &opcode{Name: ORA, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x0a */ &opcode{Name: ASL, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x0b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x0c */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x0d */ &opcode{Name: ORA, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x0e */ &opcode{Name: ASL, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x0f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x10 */ &opcode{Name: BPL, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x11 */ &opcode{Name: ORA, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x12 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x13 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x14 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x15 */ &opcode{Name: ORA, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x16 */ &opcode{Name: ASL, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x17 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x18 */ &opcode{Name: CLC, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x19 */ &opcode{Name: ORA, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x1a */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x1b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x1c */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x1d */ &opcode{Name: ORA, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x1e */ &opcode{Name: ASL, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0x1f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x20 */ &opcode{Name: JSR, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x21 */ &opcode{Name: AND, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x22 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x23 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x24 */ &opcode{Name: BIT, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x25 */ &opcode{Name: AND, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x26 */ &opcode{Name: ROL, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x27 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x28 */ &opcode{Name: PLP, Mode: implied, Cycle: 4, PageCycle: 0},
	/* 0x29 */ &opcode{Name: AND, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x2a */ &opcode{Name: ROL, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x2b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x2c */ &opcode{Name: BIT, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x2d */ &opcode{Name: AND, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x2e */ &opcode{Name: ROL, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x2f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x30 */ &opcode{Name: BMI, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x31 */ &opcode{Name: AND, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x32 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x33 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x34 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x35 */ &opcode{Name: AND, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x36 */ &opcode{Name: ROL, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x37 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x38 */ &opcode{Name: SEC, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x39 */ &opcode{Name: AND, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x3a */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x3b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x3c */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x3d */ &opcode{Name: AND, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x3e */ &opcode{Name: ROL, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0x3f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x40 */ &opcode{Name: RTI, Mode: implied, Cycle: 6, PageCycle: 0},
	/* 0x41 */ &opcode{Name: EOR, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x42 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x43 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x44 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x45 */ &opcode{Name: EOR, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x46 */ &opcode{Name: LSR, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x47 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x48 */ &opcode{Name: PHA, Mode: implied, Cycle: 3, PageCycle: 0},
	/* 0x49 */ &opcode{Name: EOR, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x4a */ &opcode{Name: LSR, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x4b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x4c */ &opcode{Name: JMP, Mode: absolute, Cycle: 3, PageCycle: 0},
	/* 0x4d */ &opcode{Name: EOR, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x4e */ &opcode{Name: LSR, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x4f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x50 */ &opcode{Name: BVC, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x51 */ &opcode{Name: EOR, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x52 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x53 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x54 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x55 */ &opcode{Name: EOR, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x56 */ &opcode{Name: LSR, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x57 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x58 */ &opcode{Name: CLI, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x59 */ &opcode{Name: EOR, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x5a */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x5b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x5c */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x5d */ &opcode{Name: EOR, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x5e */ &opcode{Name: LSR, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0x5f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x60 */ &opcode{Name: RTS, Mode: implied, Cycle: 6, PageCycle: 0},
	/* 0x61 */ &opcode{Name: ADC, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x62 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x63 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x64 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x65 */ &opcode{Name: ADC, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x66 */ &opcode{Name: ROR, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x67 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x68 */ &opcode{Name: PLA, Mode: implied, Cycle: 4, PageCycle: 0},
	/* 0x69 */ &opcode{Name: ADC, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x6a */ &opcode{Name: ROR, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x6b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x6c */ &opcode{Name: JMP, Mode: indirect, Cycle: 5, PageCycle: 0},
	/* 0x6d */ &opcode{Name: ADC, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x6e */ &opcode{Name: ROR, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x6f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x70 */ &opcode{Name: BVS, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x71 */ &opcode{Name: ADC, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x72 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x73 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x74 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x75 */ &opcode{Name: ADC, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x76 */ &opcode{Name: ROR, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x77 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x78 */ &opcode{Name: SEI, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x79 */ &opcode{Name: ADC, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x7a */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x7b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x7c */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x7d */ &opcode{Name: ADC, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x7e */ &opcode{Name: ROR, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0x7f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x80 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x81 */ &opcode{Name: STA, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x82 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x83 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x84 */ &opcode{Name: STY, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x85 */ &opcode{Name: STA, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x86 */ &opcode{Name: STX, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x87 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x88 */ &opcode{Name: DEY, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x89 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x8a */ &opcode{Name: TXA, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x8b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x8c */ &opcode{Name: STY, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x8d */ &opcode{Name: STA, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x8e */ &opcode{Name: STX, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x8f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x90 */ &opcode{Name: BCC, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x91 */ &opcode{Name: STA, Mode: indirectIndexed, Cycle: 6, PageCycle: 0},
	/* 0x92 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x93 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x94 */ &opcode{Name: STY, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x95 */ &opcode{Name: STA, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x96 */ &opcode{Name: STX, Mode: zeroPageY, Cycle: 4, PageCycle: 0},
	/* 0x97 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x98 */ &opcode{Name: TYA, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x99 */ &opcode{Name: STA, Mode: absoluteY, Cycle: 5, PageCycle: 0},
	/* 0x9a */ &opcode{Name: TXS, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x9b */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x9c */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x9d */ &opcode{Name: STA, Mode: absoluteX, Cycle: 5, PageCycle: 0},
	/* 0x9e */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0x9f */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xa0 */ &opcode{Name: LDY, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xa1 */ &opcode{Name: LDA, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0xa2 */ &opcode{Name: LDX, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xa3 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xa4 */ &opcode{Name: LDY, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xa5 */ &opcode{Name: LDA, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xa6 */ &opcode{Name: LDX, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xa7 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xa8 */ &opcode{Name: TAY, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xa9 */ &opcode{Name: LDA, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xaa */ &opcode{Name: TAX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xab */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xac */ &opcode{Name: LDY, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xad */ &opcode{Name: LDA, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xae */ &opcode{Name: LDX, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xaf */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xb0 */ &opcode{Name: BCS, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0xb1 */ &opcode{Name: LDA, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0xb2 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xb3 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xb4 */ &opcode{Name: LDY, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xb5 */ &opcode{Name: LDA, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xb6 */ &opcode{Name: LDX, Mode: zeroPageY, Cycle: 4, PageCycle: 0},
	/* 0xb7 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xb8 */ &opcode{Name: CLV, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xb9 */ &opcode{Name: LDA, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xba */ &opcode{Name: TSX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xbb */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xbc */ &opcode{Name: LDY, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xbd */ &opcode{Name: LDA, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xbe */ &opcode{Name: LDX, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xbf */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xc0 */ &opcode{Name: CPY, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xc1 */ &opcode{Name: CMP, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0xc2 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xc3 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xc4 */ &opcode{Name: CPY, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xc5 */ &opcode{Name: CMP, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xc6 */ &opcode{Name: DEC, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0xc7 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xc8 */ &opcode{Name: INY, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xc9 */ &opcode{Name: CMP, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xca */ &opcode{Name: DEX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xcb */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xcc */ &opcode{Name: CPY, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xcd */ &opcode{Name: CMP, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xce */ &opcode{Name: DEC, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0xcf */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xd0 */ &opcode{Name: BNE, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0xd1 */ &opcode{Name: CMP, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0xd2 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xd3 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xd4 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xd5 */ &opcode{Name: CMP, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xd6 */ &opcode{Name: DEC, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0xd7 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xd8 */ &opcode{Name: CLD, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xd9 */ &opcode{Name: CMP, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xda */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xdb */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xdc */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xdd */ &opcode{Name: CMP, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xde */ &opcode{Name: DEC, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0xdf */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xe0 */ &opcode{Name: CPX, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xe1 */ &opcode{Name: SBC, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0xe2 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xe3 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xe4 */ &opcode{Name: CPX, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xe5 */ &opcode{Name: SBC, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xe6 */ &opcode{Name: INC, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0xe7 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xe8 */ &opcode{Name: INX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xe9 */ &opcode{Name: SBC, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xea */ &opcode{Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xeb */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xec */ &opcode{Name: CPX, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xed */ &opcode{Name: SBC, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xee */ &opcode{Name: INC, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0xef */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xf0 */ &opcode{Name: BEQ, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0xf1 */ &opcode{Name: SBC, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0xf2 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xf3 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xf4 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xf5 */ &opcode{Name: SBC, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xf6 */ &opcode{Name: INC, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0xf7 */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xf8 */ &opcode{Name: SED, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xf9 */ &opcode{Name: SBC, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xfa */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xfb */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xfc */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
	/* 0xfd */ &opcode{Name: SBC, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xfe */ &opcode{Name: INC, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0xff */ &opcode{Name: UnknownMnemonic, Mode: implied, Cycle: 0, PageCycle: 0}, /* not yet implemented */
}
