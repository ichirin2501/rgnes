package nes

type addressingMode int

const (
	absolute addressingMode = iota + 1
	absoluteX
	absoluteX_D
	absoluteY
	absoluteY_D
	accumulator
	immediate
	implied
	indexedIndirect
	indirect
	indirectIndexed
	indirectIndexed_D
	relative
	zeroPage
	zeroPageX
	zeroPageY
)

func (a addressingMode) String() string {
	switch a {
	case absolute:
		return "absolute"
	case absoluteX:
		return "absoluteX"
	case absoluteX_D:
		return "absoluteX"
	case absoluteY:
		return "absoluteY"
	case absoluteY_D:
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
	case indirectIndexed_D:
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
	KIL
	SLO
	ANC
	RLA
	SRE
	ALR
	RRA
	ARR
	SAX
	XAA
	AHX
	TAS
	SHY
	SHX
	LAX
	LAS
	DCP
	AXS
	ISB
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
	case KIL:
		return "KIL"
	case SLO:
		return "SLO"
	case ANC:
		return "ANC"
	case RLA:
		return "RLA"
	case SRE:
		return "SRE"
	case ALR:
		return "ALR"
	case RRA:
		return "RRA"
	case ARR:
		return "ARR"
	case SAX:
		return "SAX"
	case XAA:
		return "XAA"
	case AHX:
		return "AHX"
	case TAS:
		return "TAS"
	case SHY:
		return "SHY"
	case SHX:
		return "SHX"
	case LAX:
		return "LAX"
	case LAS:
		return "LAS"
	case DCP:
		return "DCP"
	case AXS:
		return "AXS"
	case ISB:
		return "ISB"
	default:
		panic("Unable to reach here")
	}
}

type opcode struct {
	name       Mnemonic
	mode       addressingMode
	cycle      int
	pageCycle  int
	unofficial bool
}

var opcodeMap = []*opcode{
	/* 0x00 */ {name: BRK, mode: implied, cycle: 7, pageCycle: 0},
	/* 0x01 */ {name: ORA, mode: indexedIndirect, cycle: 6, pageCycle: 0},
	/* 0x02 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x03 */ {name: SLO, mode: indexedIndirect, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x04 */ {name: NOP, mode: zeroPage, cycle: 3, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x05 */ {name: ORA, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0x06 */ {name: ASL, mode: zeroPage, cycle: 5, pageCycle: 0},
	/* 0x07 */ {name: SLO, mode: zeroPage, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x08 */ {name: PHP, mode: implied, cycle: 3, pageCycle: 0},
	/* 0x09 */ {name: ORA, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0x0A */ {name: ASL, mode: accumulator, cycle: 2, pageCycle: 0},
	/* 0x0B */ {name: ANC, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x0C */ {name: NOP, mode: absolute, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x0D */ {name: ORA, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0x0E */ {name: ASL, mode: absolute, cycle: 6, pageCycle: 0},
	/* 0x0F */ {name: SLO, mode: absolute, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x10 */ {name: BPL, mode: relative, cycle: 2, pageCycle: 1},
	/* 0x11 */ {name: ORA, mode: indirectIndexed, cycle: 5, pageCycle: 1},
	/* 0x12 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x13 */ {name: SLO, mode: indirectIndexed_D, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x14 */ {name: NOP, mode: zeroPageX, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x15 */ {name: ORA, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0x16 */ {name: ASL, mode: zeroPageX, cycle: 6, pageCycle: 0},
	/* 0x17 */ {name: SLO, mode: zeroPageX, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x18 */ {name: CLC, mode: implied, cycle: 2, pageCycle: 0},
	/* 0x19 */ {name: ORA, mode: absoluteY, cycle: 4, pageCycle: 1},
	/* 0x1A */ {name: NOP, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x1B */ {name: SLO, mode: absoluteY_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x1C */ {name: NOP, mode: absoluteX, cycle: 4, pageCycle: 1, unofficial: true}, /* undocumented opcode */
	/* 0x1D */ {name: ORA, mode: absoluteX, cycle: 4, pageCycle: 1},
	/* 0x1E */ {name: ASL, mode: absoluteX_D, cycle: 7, pageCycle: 0},
	/* 0x1F */ {name: SLO, mode: absoluteX_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x20 */ {name: JSR, mode: absolute, cycle: 6, pageCycle: 0},
	/* 0x21 */ {name: AND, mode: indexedIndirect, cycle: 6, pageCycle: 0},
	/* 0x22 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x23 */ {name: RLA, mode: indexedIndirect, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x24 */ {name: BIT, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0x25 */ {name: AND, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0x26 */ {name: ROL, mode: zeroPage, cycle: 5, pageCycle: 0},
	/* 0x27 */ {name: RLA, mode: zeroPage, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x28 */ {name: PLP, mode: implied, cycle: 4, pageCycle: 0},
	/* 0x29 */ {name: AND, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0x2A */ {name: ROL, mode: accumulator, cycle: 2, pageCycle: 0},
	/* 0x2B */ {name: ANC, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x2C */ {name: BIT, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0x2D */ {name: AND, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0x2E */ {name: ROL, mode: absolute, cycle: 6, pageCycle: 0},
	/* 0x2F */ {name: RLA, mode: absolute, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x30 */ {name: BMI, mode: relative, cycle: 2, pageCycle: 1},
	/* 0x31 */ {name: AND, mode: indirectIndexed, cycle: 5, pageCycle: 1},
	/* 0x32 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x33 */ {name: RLA, mode: indirectIndexed_D, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x34 */ {name: NOP, mode: zeroPageX, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x35 */ {name: AND, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0x36 */ {name: ROL, mode: zeroPageX, cycle: 6, pageCycle: 0},
	/* 0x37 */ {name: RLA, mode: zeroPageX, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x38 */ {name: SEC, mode: implied, cycle: 2, pageCycle: 0},
	/* 0x39 */ {name: AND, mode: absoluteY, cycle: 4, pageCycle: 1},
	/* 0x3A */ {name: NOP, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x3B */ {name: RLA, mode: absoluteY_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x3C */ {name: NOP, mode: absoluteX, cycle: 4, pageCycle: 1, unofficial: true}, /* undocumented opcode */
	/* 0x3D */ {name: AND, mode: absoluteX, cycle: 4, pageCycle: 1},
	/* 0x3E */ {name: ROL, mode: absoluteX_D, cycle: 7, pageCycle: 0},
	/* 0x3F */ {name: RLA, mode: absoluteX_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x40 */ {name: RTI, mode: implied, cycle: 6, pageCycle: 0},
	/* 0x41 */ {name: EOR, mode: indexedIndirect, cycle: 6, pageCycle: 0},
	/* 0x42 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x43 */ {name: SRE, mode: indexedIndirect, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x44 */ {name: NOP, mode: zeroPage, cycle: 3, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x45 */ {name: EOR, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0x46 */ {name: LSR, mode: zeroPage, cycle: 5, pageCycle: 0},
	/* 0x47 */ {name: SRE, mode: zeroPage, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x48 */ {name: PHA, mode: implied, cycle: 3, pageCycle: 0},
	/* 0x49 */ {name: EOR, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0x4A */ {name: LSR, mode: accumulator, cycle: 2, pageCycle: 0},
	/* 0x4B */ {name: ALR, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x4C */ {name: JMP, mode: absolute, cycle: 3, pageCycle: 0},
	/* 0x4D */ {name: EOR, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0x4E */ {name: LSR, mode: absolute, cycle: 6, pageCycle: 0},
	/* 0x4F */ {name: SRE, mode: absolute, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x50 */ {name: BVC, mode: relative, cycle: 2, pageCycle: 1},
	/* 0x51 */ {name: EOR, mode: indirectIndexed, cycle: 5, pageCycle: 1},
	/* 0x52 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x53 */ {name: SRE, mode: indirectIndexed_D, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x54 */ {name: NOP, mode: zeroPageX, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x55 */ {name: EOR, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0x56 */ {name: LSR, mode: zeroPageX, cycle: 6, pageCycle: 0},
	/* 0x57 */ {name: SRE, mode: zeroPageX, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x58 */ {name: CLI, mode: implied, cycle: 2, pageCycle: 0},
	/* 0x59 */ {name: EOR, mode: absoluteY, cycle: 4, pageCycle: 1},
	/* 0x5A */ {name: NOP, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x5B */ {name: SRE, mode: absoluteY_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x5C */ {name: NOP, mode: absoluteX, cycle: 4, pageCycle: 1, unofficial: true}, /* undocumented opcode */
	/* 0x5D */ {name: EOR, mode: absoluteX, cycle: 4, pageCycle: 1},
	/* 0x5E */ {name: LSR, mode: absoluteX_D, cycle: 7, pageCycle: 0},
	/* 0x5F */ {name: SRE, mode: absoluteX_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x60 */ {name: RTS, mode: implied, cycle: 6, pageCycle: 0},
	/* 0x61 */ {name: ADC, mode: indexedIndirect, cycle: 6, pageCycle: 0},
	/* 0x62 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x63 */ {name: RRA, mode: indexedIndirect, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x64 */ {name: NOP, mode: zeroPage, cycle: 3, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x65 */ {name: ADC, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0x66 */ {name: ROR, mode: zeroPage, cycle: 5, pageCycle: 0},
	/* 0x67 */ {name: RRA, mode: zeroPage, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x68 */ {name: PLA, mode: implied, cycle: 4, pageCycle: 0},
	/* 0x69 */ {name: ADC, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0x6A */ {name: ROR, mode: accumulator, cycle: 2, pageCycle: 0},
	/* 0x6B */ {name: ARR, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x6C */ {name: JMP, mode: indirect, cycle: 5, pageCycle: 0},
	/* 0x6D */ {name: ADC, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0x6E */ {name: ROR, mode: absolute, cycle: 6, pageCycle: 0},
	/* 0x6F */ {name: RRA, mode: absolute, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x70 */ {name: BVS, mode: relative, cycle: 2, pageCycle: 1},
	/* 0x71 */ {name: ADC, mode: indirectIndexed, cycle: 5, pageCycle: 1},
	/* 0x72 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x73 */ {name: RRA, mode: indirectIndexed_D, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x74 */ {name: NOP, mode: zeroPageX, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x75 */ {name: ADC, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0x76 */ {name: ROR, mode: zeroPageX, cycle: 6, pageCycle: 0},
	/* 0x77 */ {name: RRA, mode: zeroPageX, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x78 */ {name: SEI, mode: implied, cycle: 2, pageCycle: 0},
	/* 0x79 */ {name: ADC, mode: absoluteY, cycle: 4, pageCycle: 1},
	/* 0x7A */ {name: NOP, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x7B */ {name: RRA, mode: absoluteY_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x7C */ {name: NOP, mode: absoluteX, cycle: 4, pageCycle: 1, unofficial: true}, /* undocumented opcode */
	/* 0x7D */ {name: ADC, mode: absoluteX, cycle: 4, pageCycle: 1},
	/* 0x7E */ {name: ROR, mode: absoluteX_D, cycle: 7, pageCycle: 0},
	/* 0x7F */ {name: RRA, mode: absoluteX_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x80 */ {name: NOP, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x81 */ {name: STA, mode: indexedIndirect, cycle: 6, pageCycle: 0},
	/* 0x82 */ {name: NOP, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x83 */ {name: SAX, mode: indexedIndirect, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x84 */ {name: STY, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0x85 */ {name: STA, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0x86 */ {name: STX, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0x87 */ {name: SAX, mode: zeroPage, cycle: 3, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x88 */ {name: DEY, mode: implied, cycle: 2, pageCycle: 0},
	/* 0x89 */ {name: NOP, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x8A */ {name: TXA, mode: implied, cycle: 2, pageCycle: 0},
	/* 0x8B */ {name: XAA, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x8C */ {name: STY, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0x8D */ {name: STA, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0x8E */ {name: STX, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0x8F */ {name: SAX, mode: absolute, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x90 */ {name: BCC, mode: relative, cycle: 2, pageCycle: 1},
	/* 0x91 */ {name: STA, mode: indirectIndexed_D, cycle: 6, pageCycle: 0},
	/* 0x92 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x93 */ {name: AHX, mode: indirectIndexed_D, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x94 */ {name: STY, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0x95 */ {name: STA, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0x96 */ {name: STX, mode: zeroPageY, cycle: 4, pageCycle: 0},
	/* 0x97 */ {name: SAX, mode: zeroPageY, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x98 */ {name: TYA, mode: implied, cycle: 2, pageCycle: 0},
	/* 0x99 */ {name: STA, mode: absoluteY_D, cycle: 5, pageCycle: 0},
	/* 0x9A */ {name: TXS, mode: implied, cycle: 2, pageCycle: 0},
	/* 0x9B */ {name: TAS, mode: absoluteY_D, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x9C */ {name: SHY, mode: absoluteX_D, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x9D */ {name: STA, mode: absoluteX_D, cycle: 5, pageCycle: 0},
	/* 0x9E */ {name: SHX, mode: absoluteY_D, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0x9F */ {name: AHX, mode: absoluteY_D, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xA0 */ {name: LDY, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0xA1 */ {name: LDA, mode: indexedIndirect, cycle: 6, pageCycle: 0},
	/* 0xA2 */ {name: LDX, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0xA3 */ {name: LAX, mode: indexedIndirect, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xA4 */ {name: LDY, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0xA5 */ {name: LDA, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0xA6 */ {name: LDX, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0xA7 */ {name: LAX, mode: zeroPage, cycle: 3, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xA8 */ {name: TAY, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xA9 */ {name: LDA, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0xAA */ {name: TAX, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xAB */ {name: LAX, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xAC */ {name: LDY, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0xAD */ {name: LDA, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0xAE */ {name: LDX, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0xAF */ {name: LAX, mode: absolute, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xB0 */ {name: BCS, mode: relative, cycle: 2, pageCycle: 1},
	/* 0xB1 */ {name: LDA, mode: indirectIndexed, cycle: 5, pageCycle: 1},
	/* 0xB2 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xB3 */ {name: LAX, mode: indirectIndexed, cycle: 5, pageCycle: 1, unofficial: true}, /* undocumented opcode */
	/* 0xB4 */ {name: LDY, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0xB5 */ {name: LDA, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0xB6 */ {name: LDX, mode: zeroPageY, cycle: 4, pageCycle: 0},
	/* 0xB7 */ {name: LAX, mode: zeroPageY, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xB8 */ {name: CLV, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xB9 */ {name: LDA, mode: absoluteY, cycle: 4, pageCycle: 1},
	/* 0xBA */ {name: TSX, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xBB */ {name: LAS, mode: absoluteY, cycle: 4, pageCycle: 1, unofficial: true}, /* undocumented opcode */
	/* 0xBC */ {name: LDY, mode: absoluteX, cycle: 4, pageCycle: 1},
	/* 0xBD */ {name: LDA, mode: absoluteX, cycle: 4, pageCycle: 1},
	/* 0xBE */ {name: LDX, mode: absoluteY, cycle: 4, pageCycle: 1},
	/* 0xBF */ {name: LAX, mode: absoluteY, cycle: 4, pageCycle: 1, unofficial: true}, /* undocumented opcode */
	/* 0xC0 */ {name: CPY, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0xC1 */ {name: CMP, mode: indexedIndirect, cycle: 6, pageCycle: 0},
	/* 0xC2 */ {name: NOP, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xC3 */ {name: DCP, mode: indexedIndirect, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xC4 */ {name: CPY, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0xC5 */ {name: CMP, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0xC6 */ {name: DEC, mode: zeroPage, cycle: 5, pageCycle: 0},
	/* 0xC7 */ {name: DCP, mode: zeroPage, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xC8 */ {name: INY, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xC9 */ {name: CMP, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0xCA */ {name: DEX, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xCB */ {name: AXS, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xCC */ {name: CPY, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0xCD */ {name: CMP, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0xCE */ {name: DEC, mode: absolute, cycle: 6, pageCycle: 0},
	/* 0xCF */ {name: DCP, mode: absolute, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xD0 */ {name: BNE, mode: relative, cycle: 2, pageCycle: 1},
	/* 0xD1 */ {name: CMP, mode: indirectIndexed, cycle: 5, pageCycle: 1},
	/* 0xD2 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xD3 */ {name: DCP, mode: indirectIndexed_D, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xD4 */ {name: NOP, mode: zeroPageX, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xD5 */ {name: CMP, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0xD6 */ {name: DEC, mode: zeroPageX, cycle: 6, pageCycle: 0},
	/* 0xD7 */ {name: DCP, mode: zeroPageX, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xD8 */ {name: CLD, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xD9 */ {name: CMP, mode: absoluteY, cycle: 4, pageCycle: 1},
	/* 0xDA */ {name: NOP, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xDB */ {name: DCP, mode: absoluteY_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xDC */ {name: NOP, mode: absoluteX, cycle: 4, pageCycle: 1, unofficial: true}, /* undocumented opcode */
	/* 0xDD */ {name: CMP, mode: absoluteX, cycle: 4, pageCycle: 1},
	/* 0xDE */ {name: DEC, mode: absoluteX_D, cycle: 7, pageCycle: 0},
	/* 0xDF */ {name: DCP, mode: absoluteX_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xE0 */ {name: CPX, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0xE1 */ {name: SBC, mode: indexedIndirect, cycle: 6, pageCycle: 0},
	/* 0xE2 */ {name: NOP, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xE3 */ {name: ISB, mode: indexedIndirect, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xE4 */ {name: CPX, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0xE5 */ {name: SBC, mode: zeroPage, cycle: 3, pageCycle: 0},
	/* 0xE6 */ {name: INC, mode: zeroPage, cycle: 5, pageCycle: 0},
	/* 0xE7 */ {name: ISB, mode: zeroPage, cycle: 5, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xE8 */ {name: INX, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xE9 */ {name: SBC, mode: immediate, cycle: 2, pageCycle: 0},
	/* 0xEA */ {name: NOP, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xEB */ {name: SBC, mode: immediate, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xEC */ {name: CPX, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0xED */ {name: SBC, mode: absolute, cycle: 4, pageCycle: 0},
	/* 0xEE */ {name: INC, mode: absolute, cycle: 6, pageCycle: 0},
	/* 0xEF */ {name: ISB, mode: absolute, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xF0 */ {name: BEQ, mode: relative, cycle: 2, pageCycle: 1},
	/* 0xF1 */ {name: SBC, mode: indirectIndexed, cycle: 5, pageCycle: 1},
	/* 0xF2 */ {name: KIL, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xF3 */ {name: ISB, mode: indirectIndexed_D, cycle: 8, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xF4 */ {name: NOP, mode: zeroPageX, cycle: 4, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xF5 */ {name: SBC, mode: zeroPageX, cycle: 4, pageCycle: 0},
	/* 0xF6 */ {name: INC, mode: zeroPageX, cycle: 6, pageCycle: 0},
	/* 0xF7 */ {name: ISB, mode: zeroPageX, cycle: 6, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xF8 */ {name: SED, mode: implied, cycle: 2, pageCycle: 0},
	/* 0xF9 */ {name: SBC, mode: absoluteY, cycle: 4, pageCycle: 1},
	/* 0xFA */ {name: NOP, mode: implied, cycle: 2, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xFB */ {name: ISB, mode: absoluteY_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
	/* 0xFC */ {name: NOP, mode: absoluteX, cycle: 4, pageCycle: 1, unofficial: true}, /* undocumented opcode */
	/* 0xFD */ {name: SBC, mode: absoluteX, cycle: 4, pageCycle: 1},
	/* 0xFE */ {name: INC, mode: absoluteX_D, cycle: 7, pageCycle: 0},
	/* 0xFF */ {name: ISB, mode: absoluteX_D, cycle: 7, pageCycle: 0, unofficial: true}, /* undocumented opcode */
}

// ref. https://www.qmtpro.com/~nes/misc/nestest.log
// absolute:$ZZZZ                       [JMP JSR]
// absolute:$ZZZZ = ZZ                  [ADC AND ASL BIT CMP CPX CPY DCP DEC EOR INC ISB LAX LDA LDX LDY LSR NOP ORA RLA ROL ROR RRA SAX SBC SLO SRE STA STX STY]
// absoluteX:$ZZZZ,X @ ZZZZ = ZZ         [ADC AND ASL CMP DCP DEC EOR INC ISB LDA LDY LSR NOP ORA RLA ROL ROR RRA SBC SLO SRE STA]
// absoluteY:$ZZZZ,Y @ ZZZZ = ZZ         [ADC AND CMP DCP EOR ISB LAX LDA LDX ORA RLA RRA SBC SLO SRE STA]
// accumulator:A                           [ASL LSR ROL ROR]
// immediate:#$ZZ                        [ADC AND CMP CPX CPY EOR LDA LDX LDY NOP ORA SBC]
// implied:                            [CLC CLD CLV DEX DEY INX INY NOP PHA PHP PLA PLP RTI RTS SEC SED SEI TAX TAY TSX TXA TXS TYA]
// indexedIndirect:($aa,X) @ bb = cccc = dd    [ADC AND CMP DCP EOR ISB LAX LDA ORA RLA RRA SAX SBC SLO SRE STA]
// indirect:($ZZZZ) = ZZZZ              [JMP]
// indirectIndexed:($ZZ),Y = ZZZZ @ ZZZZ = ZZ  [ADC AND CMP DCP EOR ISB LAX LDA ORA RLA RRA SBC SLO SRE STA]
// relative:$ZZZZ                       [BCC BCS BEQ BMI BNE BPL BVC BVS]
// zeroPage:$ZZ = ZZ                    [ADC AND ASL BIT CMP CPX CPY DCP DEC EOR INC ISB LAX LDA LDX LDY LSR NOP ORA RLA ROL ROR RRA SAX SBC SLO SRE STA STX STY]
// zeroPageX:$ZZ,X @ ZZ = ZZ             [ADC AND ASL CMP DCP DEC EOR INC ISB LDA LDY LSR NOP ORA RLA ROL ROR RRA SBC SLO SRE STA STY]
// zeroPageY:$ZZ,Y @ ZZ = ZZ             [LAX LDX SAX STX]
