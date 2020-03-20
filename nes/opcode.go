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
	Name  Mnemonic
	Mode  addressingMode
	Size  int
	Cycle int
}

var opcodeMap = map[byte]opcode{
	0xA9: opcode{Name: LDA, Mode: immediate, Size: 2, Cycle: 2},
	0xA5: opcode{Name: LDA, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB5: opcode{Name: LDA, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xAD: opcode{Name: LDA, Mode: absolute, Size: 3, Cycle: 4},
	0xBD: opcode{Name: LDA, Mode: absoluteX, Size: 3, Cycle: 4},
	0xB9: opcode{Name: LDA, Mode: absoluteY, Size: 3, Cycle: 4},
	0xA1: opcode{Name: LDA, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0xB1: opcode{Name: LDA, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0xA2: opcode{Name: LDX, Mode: immediate, Size: 2, Cycle: 2},
	0xA6: opcode{Name: LDX, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB6: opcode{Name: LDX, Mode: zeroPageY, Size: 2, Cycle: 4},
	0xAE: opcode{Name: LDX, Mode: absolute, Size: 3, Cycle: 4},
	0xBE: opcode{Name: LDX, Mode: absoluteY, Size: 3, Cycle: 4},
	0xA0: opcode{Name: LDY, Mode: immediate, Size: 2, Cycle: 2},
	0xA4: opcode{Name: LDY, Mode: zeroPage, Size: 2, Cycle: 3},
	0xB4: opcode{Name: LDY, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xAC: opcode{Name: LDY, Mode: absolute, Size: 3, Cycle: 4},
	0xBC: opcode{Name: LDY, Mode: absoluteX, Size: 3, Cycle: 4},
	0x85: opcode{Name: STA, Mode: zeroPage, Size: 2, Cycle: 3},
	0x95: opcode{Name: STA, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x8D: opcode{Name: STA, Mode: absolute, Size: 3, Cycle: 4},
	0x9D: opcode{Name: STA, Mode: absoluteX, Size: 3, Cycle: 5},
	0x99: opcode{Name: STA, Mode: absoluteY, Size: 3, Cycle: 5},
	0x81: opcode{Name: STA, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x91: opcode{Name: STA, Mode: indirectIndexed, Size: 2, Cycle: 6},
	0x86: opcode{Name: STX, Mode: zeroPage, Size: 2, Cycle: 3},
	0x96: opcode{Name: STX, Mode: zeroPageY, Size: 2, Cycle: 4},
	0x8E: opcode{Name: STX, Mode: absolute, Size: 3, Cycle: 4},
	0x84: opcode{Name: STY, Mode: zeroPage, Size: 2, Cycle: 3},
	0x94: opcode{Name: STY, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x8C: opcode{Name: STY, Mode: absolute, Size: 3, Cycle: 4},
	0xAA: opcode{Name: TAX, Mode: implied, Size: 1, Cycle: 2},
	0xA8: opcode{Name: TAY, Mode: implied, Size: 1, Cycle: 2},
	0xBA: opcode{Name: TSX, Mode: implied, Size: 1, Cycle: 2},
	0x8A: opcode{Name: TXA, Mode: implied, Size: 1, Cycle: 2},
	0x9A: opcode{Name: TXS, Mode: implied, Size: 1, Cycle: 2},
	0x98: opcode{Name: TYA, Mode: implied, Size: 1, Cycle: 2},
	0x69: opcode{Name: ADC, Mode: immediate, Size: 2, Cycle: 2},
	0x65: opcode{Name: ADC, Mode: zeroPage, Size: 2, Cycle: 3},
	0x75: opcode{Name: ADC, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x6D: opcode{Name: ADC, Mode: absolute, Size: 3, Cycle: 4},
	0x7D: opcode{Name: ADC, Mode: absoluteX, Size: 3, Cycle: 4},
	0x79: opcode{Name: ADC, Mode: absoluteY, Size: 3, Cycle: 4},
	0x61: opcode{Name: ADC, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x71: opcode{Name: ADC, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x29: opcode{Name: AND, Mode: immediate, Size: 2, Cycle: 2},
	0x25: opcode{Name: AND, Mode: zeroPage, Size: 2, Cycle: 3},
	0x35: opcode{Name: AND, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x2D: opcode{Name: AND, Mode: absolute, Size: 3, Cycle: 4},
	0x3D: opcode{Name: AND, Mode: absoluteX, Size: 3, Cycle: 4},
	0x39: opcode{Name: AND, Mode: absoluteY, Size: 3, Cycle: 4},
	0x21: opcode{Name: AND, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x31: opcode{Name: AND, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x0A: opcode{Name: ASL, Mode: accumulator, Size: 1, Cycle: 2},
	0x06: opcode{Name: ASL, Mode: zeroPage, Size: 2, Cycle: 5},
	0x16: opcode{Name: ASL, Mode: zeroPageX, Size: 2, Cycle: 6},
	0x0E: opcode{Name: ASL, Mode: absolute, Size: 3, Cycle: 6},
	0x1E: opcode{Name: ASL, Mode: absoluteX, Size: 3, Cycle: 7},
	0x24: opcode{Name: BIT, Mode: zeroPage, Size: 2, Cycle: 3},
	0x2C: opcode{Name: BIT, Mode: absolute, Size: 3, Cycle: 4},
	0xC9: opcode{Name: CMP, Mode: immediate, Size: 2, Cycle: 2},
	0xC5: opcode{Name: CMP, Mode: zeroPage, Size: 2, Cycle: 3},
	0xD5: opcode{Name: CMP, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xCD: opcode{Name: CMP, Mode: absolute, Size: 3, Cycle: 4},
	0xDD: opcode{Name: CMP, Mode: absoluteX, Size: 3, Cycle: 4},
	0xD9: opcode{Name: CMP, Mode: absoluteY, Size: 3, Cycle: 4},
	0xC1: opcode{Name: CMP, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0xD1: opcode{Name: CMP, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0xE0: opcode{Name: CPX, Mode: immediate, Size: 2, Cycle: 2},
	0xE4: opcode{Name: CPX, Mode: zeroPage, Size: 2, Cycle: 3},
	0xEC: opcode{Name: CPX, Mode: absolute, Size: 3, Cycle: 4},
	0xC0: opcode{Name: CPY, Mode: immediate, Size: 2, Cycle: 2},
	0xC4: opcode{Name: CPY, Mode: zeroPage, Size: 2, Cycle: 3},
	0xCC: opcode{Name: CPY, Mode: absolute, Size: 3, Cycle: 4},
	0xC6: opcode{Name: DEC, Mode: zeroPage, Size: 2, Cycle: 5},
	0xD6: opcode{Name: DEC, Mode: zeroPageX, Size: 2, Cycle: 6},
	0xCE: opcode{Name: DEC, Mode: absolute, Size: 3, Cycle: 6},
	0xDE: opcode{Name: DEC, Mode: absoluteX, Size: 3, Cycle: 7},
	0xCA: opcode{Name: DEX, Mode: implied, Size: 1, Cycle: 2},
	0x88: opcode{Name: DEY, Mode: implied, Size: 1, Cycle: 2},
	0x49: opcode{Name: EOR, Mode: immediate, Size: 2, Cycle: 2},
	0x45: opcode{Name: EOR, Mode: zeroPage, Size: 2, Cycle: 3},
	0x55: opcode{Name: EOR, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x4D: opcode{Name: EOR, Mode: absolute, Size: 3, Cycle: 4},
	0x5D: opcode{Name: EOR, Mode: absoluteX, Size: 3, Cycle: 4},
	0x59: opcode{Name: EOR, Mode: absoluteY, Size: 3, Cycle: 4},
	0x41: opcode{Name: EOR, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x51: opcode{Name: EOR, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0xE6: opcode{Name: INC, Mode: zeroPage, Size: 2, Cycle: 5},
	0xF6: opcode{Name: INC, Mode: zeroPageX, Size: 2, Cycle: 6},
	0xEE: opcode{Name: INC, Mode: absolute, Size: 3, Cycle: 6},
	0xFE: opcode{Name: INC, Mode: absoluteX, Size: 3, Cycle: 7},
	0xE8: opcode{Name: INX, Mode: implied, Size: 1, Cycle: 2},
	0xC8: opcode{Name: INY, Mode: implied, Size: 1, Cycle: 2},
	0x4A: opcode{Name: LSR, Mode: accumulator, Size: 1, Cycle: 2},
	0x46: opcode{Name: LSR, Mode: zeroPage, Size: 2, Cycle: 5},
	0x56: opcode{Name: LSR, Mode: zeroPageX, Size: 2, Cycle: 6},
	0x4E: opcode{Name: LSR, Mode: absolute, Size: 3, Cycle: 6},
	0x5E: opcode{Name: LSR, Mode: absoluteX, Size: 3, Cycle: 7},
	0x09: opcode{Name: ORA, Mode: immediate, Size: 2, Cycle: 2},
	0x05: opcode{Name: ORA, Mode: zeroPage, Size: 2, Cycle: 3},
	0x15: opcode{Name: ORA, Mode: zeroPageX, Size: 2, Cycle: 4},
	0x0D: opcode{Name: ORA, Mode: absolute, Size: 3, Cycle: 4},
	0x1D: opcode{Name: ORA, Mode: absoluteX, Size: 3, Cycle: 4},
	0x19: opcode{Name: ORA, Mode: absoluteY, Size: 3, Cycle: 4},
	0x01: opcode{Name: ORA, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0x11: opcode{Name: ORA, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x2A: opcode{Name: ROL, Mode: accumulator, Size: 1, Cycle: 2},
	0x26: opcode{Name: ROL, Mode: zeroPage, Size: 2, Cycle: 5},
	0x36: opcode{Name: ROL, Mode: zeroPageX, Size: 2, Cycle: 6},
	0x2E: opcode{Name: ROL, Mode: absolute, Size: 3, Cycle: 6},
	0x3E: opcode{Name: ROL, Mode: absoluteX, Size: 3, Cycle: 7},
	0x6A: opcode{Name: ROR, Mode: accumulator, Size: 1, Cycle: 2},
	0x66: opcode{Name: ROR, Mode: zeroPage, Size: 2, Cycle: 5},
	0x76: opcode{Name: ROR, Mode: zeroPageX, Size: 2, Cycle: 6},
	0x6E: opcode{Name: ROR, Mode: absolute, Size: 3, Cycle: 6},
	0x7E: opcode{Name: ROR, Mode: absoluteX, Size: 3, Cycle: 7},
	0xE9: opcode{Name: SBC, Mode: immediate, Size: 2, Cycle: 2},
	0xE5: opcode{Name: SBC, Mode: zeroPage, Size: 2, Cycle: 3},
	0xF5: opcode{Name: SBC, Mode: zeroPageX, Size: 2, Cycle: 4},
	0xED: opcode{Name: SBC, Mode: absolute, Size: 3, Cycle: 4},
	0xFD: opcode{Name: SBC, Mode: absoluteX, Size: 3, Cycle: 4},
	0xF9: opcode{Name: SBC, Mode: absoluteY, Size: 3, Cycle: 4},
	0xE1: opcode{Name: SBC, Mode: indexedIndirect, Size: 2, Cycle: 6},
	0xF1: opcode{Name: SBC, Mode: indirectIndexed, Size: 2, Cycle: 5},
	0x48: opcode{Name: PHA, Mode: implied, Size: 1, Cycle: 3},
	0x08: opcode{Name: PHP, Mode: implied, Size: 1, Cycle: 3},
	0x68: opcode{Name: PLA, Mode: implied, Size: 1, Cycle: 4},
	0x28: opcode{Name: PLP, Mode: implied, Size: 1, Cycle: 4},
	0x4C: opcode{Name: JMP, Mode: absolute, Size: 3, Cycle: 3},
	0x6C: opcode{Name: JMP, Mode: indirect, Size: 3, Cycle: 5},
	0x20: opcode{Name: JSR, Mode: absolute, Size: 3, Cycle: 6},
	0x60: opcode{Name: RTS, Mode: implied, Size: 1, Cycle: 6},
	0x40: opcode{Name: RTI, Mode: implied, Size: 1, Cycle: 6},
	0x90: opcode{Name: BCC, Mode: relative, Size: 2, Cycle: 2},
	0xB0: opcode{Name: BCS, Mode: relative, Size: 2, Cycle: 2},
	0xF0: opcode{Name: BEQ, Mode: relative, Size: 2, Cycle: 2},
	0x30: opcode{Name: BMI, Mode: relative, Size: 2, Cycle: 2},
	0xD0: opcode{Name: BNE, Mode: relative, Size: 2, Cycle: 2},
	0x10: opcode{Name: BPL, Mode: relative, Size: 2, Cycle: 2},
	0x50: opcode{Name: BVC, Mode: relative, Size: 2, Cycle: 2},
	0x70: opcode{Name: BVS, Mode: relative, Size: 2, Cycle: 2},
	0x18: opcode{Name: CLC, Mode: implied, Size: 1, Cycle: 2},
	0xD8: opcode{Name: CLD, Mode: implied, Size: 1, Cycle: 2},
	0x58: opcode{Name: CLI, Mode: implied, Size: 1, Cycle: 2},
	0xB8: opcode{Name: CLV, Mode: implied, Size: 1, Cycle: 2},
	0x38: opcode{Name: SEC, Mode: implied, Size: 1, Cycle: 2},
	0xF8: opcode{Name: SED, Mode: implied, Size: 1, Cycle: 2},
	0x78: opcode{Name: SEI, Mode: implied, Size: 1, Cycle: 2},
	0x00: opcode{Name: BRK, Mode: implied, Size: 1, Cycle: 7},
	0xEA: opcode{Name: NOP, Mode: implied, Size: 1, Cycle: 2},
}
