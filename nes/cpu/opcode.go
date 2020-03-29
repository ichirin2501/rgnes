package cpu

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
	ISC
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
	case ISC:
		return "ISC"
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
	/* 0x02 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x03 */ &opcode{Name: SLO, Mode: indexedIndirect, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0x04 */ &opcode{Name: NOP, Mode: zeroPage, Cycle: 3, PageCycle: 0}, /* undocumented opcode */
	/* 0x05 */ &opcode{Name: ORA, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x06 */ &opcode{Name: ASL, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x07 */ &opcode{Name: SLO, Mode: zeroPage, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0x08 */ &opcode{Name: PHP, Mode: implied, Cycle: 3, PageCycle: 0},
	/* 0x09 */ &opcode{Name: ORA, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x0A */ &opcode{Name: ASL, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x0B */ &opcode{Name: ANC, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x0C */ &opcode{Name: NOP, Mode: absolute, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0x0D */ &opcode{Name: ORA, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x0E */ &opcode{Name: ASL, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x0F */ &opcode{Name: SLO, Mode: absolute, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x10 */ &opcode{Name: BPL, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x11 */ &opcode{Name: ORA, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x12 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x13 */ &opcode{Name: SLO, Mode: indirectIndexed, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0x14 */ &opcode{Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0x15 */ &opcode{Name: ORA, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x16 */ &opcode{Name: ASL, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x17 */ &opcode{Name: SLO, Mode: zeroPageX, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x18 */ &opcode{Name: CLC, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x19 */ &opcode{Name: ORA, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x1A */ &opcode{Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x1B */ &opcode{Name: SLO, Mode: absoluteY, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0x1C */ &opcode{Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1}, /* undocumented opcode */
	/* 0x1D */ &opcode{Name: ORA, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x1E */ &opcode{Name: ASL, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0x1F */ &opcode{Name: SLO, Mode: absoluteX, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0x20 */ &opcode{Name: JSR, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x21 */ &opcode{Name: AND, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x22 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x23 */ &opcode{Name: RLA, Mode: indexedIndirect, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0x24 */ &opcode{Name: BIT, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x25 */ &opcode{Name: AND, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x26 */ &opcode{Name: ROL, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x27 */ &opcode{Name: RLA, Mode: zeroPage, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0x28 */ &opcode{Name: PLP, Mode: implied, Cycle: 4, PageCycle: 0},
	/* 0x29 */ &opcode{Name: AND, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x2A */ &opcode{Name: ROL, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x2B */ &opcode{Name: ANC, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x2C */ &opcode{Name: BIT, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x2D */ &opcode{Name: AND, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x2E */ &opcode{Name: ROL, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x2F */ &opcode{Name: RLA, Mode: absolute, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x30 */ &opcode{Name: BMI, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x31 */ &opcode{Name: AND, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x32 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x33 */ &opcode{Name: RLA, Mode: indirectIndexed, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0x34 */ &opcode{Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0x35 */ &opcode{Name: AND, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x36 */ &opcode{Name: ROL, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x37 */ &opcode{Name: RLA, Mode: zeroPageX, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x38 */ &opcode{Name: SEC, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x39 */ &opcode{Name: AND, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x3A */ &opcode{Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x3B */ &opcode{Name: RLA, Mode: absoluteY, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0x3C */ &opcode{Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1}, /* undocumented opcode */
	/* 0x3D */ &opcode{Name: AND, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x3E */ &opcode{Name: ROL, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0x3F */ &opcode{Name: RLA, Mode: absoluteX, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0x40 */ &opcode{Name: RTI, Mode: implied, Cycle: 6, PageCycle: 0},
	/* 0x41 */ &opcode{Name: EOR, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x42 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x43 */ &opcode{Name: SRE, Mode: indexedIndirect, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0x44 */ &opcode{Name: NOP, Mode: zeroPage, Cycle: 3, PageCycle: 0}, /* undocumented opcode */
	/* 0x45 */ &opcode{Name: EOR, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x46 */ &opcode{Name: LSR, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x47 */ &opcode{Name: SRE, Mode: zeroPage, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0x48 */ &opcode{Name: PHA, Mode: implied, Cycle: 3, PageCycle: 0},
	/* 0x49 */ &opcode{Name: EOR, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x4A */ &opcode{Name: LSR, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x4B */ &opcode{Name: ALR, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x4C */ &opcode{Name: JMP, Mode: absolute, Cycle: 3, PageCycle: 0},
	/* 0x4D */ &opcode{Name: EOR, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x4E */ &opcode{Name: LSR, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x4F */ &opcode{Name: SRE, Mode: absolute, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x50 */ &opcode{Name: BVC, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x51 */ &opcode{Name: EOR, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x52 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x53 */ &opcode{Name: SRE, Mode: indirectIndexed, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0x54 */ &opcode{Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0x55 */ &opcode{Name: EOR, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x56 */ &opcode{Name: LSR, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x57 */ &opcode{Name: SRE, Mode: zeroPageX, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x58 */ &opcode{Name: CLI, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x59 */ &opcode{Name: EOR, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x5A */ &opcode{Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x5B */ &opcode{Name: SRE, Mode: absoluteY, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0x5C */ &opcode{Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1}, /* undocumented opcode */
	/* 0x5D */ &opcode{Name: EOR, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x5E */ &opcode{Name: LSR, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0x5F */ &opcode{Name: SRE, Mode: absoluteX, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0x60 */ &opcode{Name: RTS, Mode: implied, Cycle: 6, PageCycle: 0},
	/* 0x61 */ &opcode{Name: ADC, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x62 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x63 */ &opcode{Name: RRA, Mode: indexedIndirect, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0x64 */ &opcode{Name: NOP, Mode: zeroPage, Cycle: 3, PageCycle: 0}, /* undocumented opcode */
	/* 0x65 */ &opcode{Name: ADC, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x66 */ &opcode{Name: ROR, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x67 */ &opcode{Name: RRA, Mode: zeroPage, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0x68 */ &opcode{Name: PLA, Mode: implied, Cycle: 4, PageCycle: 0},
	/* 0x69 */ &opcode{Name: ADC, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x6A */ &opcode{Name: ROR, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x6B */ &opcode{Name: ARR, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x6C */ &opcode{Name: JMP, Mode: indirect, Cycle: 5, PageCycle: 0},
	/* 0x6D */ &opcode{Name: ADC, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x6E */ &opcode{Name: ROR, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x6F */ &opcode{Name: RRA, Mode: absolute, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x70 */ &opcode{Name: BVS, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x71 */ &opcode{Name: ADC, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x72 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x73 */ &opcode{Name: RRA, Mode: indirectIndexed, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0x74 */ &opcode{Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0x75 */ &opcode{Name: ADC, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x76 */ &opcode{Name: ROR, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x77 */ &opcode{Name: RRA, Mode: zeroPageX, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x78 */ &opcode{Name: SEI, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x79 */ &opcode{Name: ADC, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x7A */ &opcode{Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x7B */ &opcode{Name: RRA, Mode: absoluteY, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0x7C */ &opcode{Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1}, /* undocumented opcode */
	/* 0x7D */ &opcode{Name: ADC, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x7E */ &opcode{Name: ROR, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0x7F */ &opcode{Name: RRA, Mode: absoluteX, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0x80 */ &opcode{Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x81 */ &opcode{Name: STA, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x82 */ &opcode{Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x83 */ &opcode{Name: SAX, Mode: indexedIndirect, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x84 */ &opcode{Name: STY, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x85 */ &opcode{Name: STA, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x86 */ &opcode{Name: STX, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x87 */ &opcode{Name: SAX, Mode: zeroPage, Cycle: 3, PageCycle: 0}, /* undocumented opcode */
	/* 0x88 */ &opcode{Name: DEY, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x89 */ &opcode{Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x8A */ &opcode{Name: TXA, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x8B */ &opcode{Name: XAA, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x8C */ &opcode{Name: STY, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x8D */ &opcode{Name: STA, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x8E */ &opcode{Name: STX, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x8F */ &opcode{Name: SAX, Mode: absolute, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0x90 */ &opcode{Name: BCC, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x91 */ &opcode{Name: STA, Mode: indirectIndexed, Cycle: 6, PageCycle: 0},
	/* 0x92 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0x93 */ &opcode{Name: AHX, Mode: indirectIndexed, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0x94 */ &opcode{Name: STY, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x95 */ &opcode{Name: STA, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x96 */ &opcode{Name: STX, Mode: zeroPageY, Cycle: 4, PageCycle: 0},
	/* 0x97 */ &opcode{Name: SAX, Mode: zeroPageY, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0x98 */ &opcode{Name: TYA, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x99 */ &opcode{Name: STA, Mode: absoluteY, Cycle: 5, PageCycle: 0},
	/* 0x9A */ &opcode{Name: TXS, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x9B */ &opcode{Name: TAS, Mode: absoluteY, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0x9C */ &opcode{Name: SHY, Mode: absoluteX, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0x9D */ &opcode{Name: STA, Mode: absoluteX, Cycle: 5, PageCycle: 0},
	/* 0x9E */ &opcode{Name: SHX, Mode: absoluteY, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0x9F */ &opcode{Name: AHX, Mode: absoluteY, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0xA0 */ &opcode{Name: LDY, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xA1 */ &opcode{Name: LDA, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0xA2 */ &opcode{Name: LDX, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xA3 */ &opcode{Name: LAX, Mode: indexedIndirect, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0xA4 */ &opcode{Name: LDY, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xA5 */ &opcode{Name: LDA, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xA6 */ &opcode{Name: LDX, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xA7 */ &opcode{Name: LAX, Mode: zeroPage, Cycle: 3, PageCycle: 0}, /* undocumented opcode */
	/* 0xA8 */ &opcode{Name: TAY, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xA9 */ &opcode{Name: LDA, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xAA */ &opcode{Name: TAX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xAB */ &opcode{Name: LAX, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xAC */ &opcode{Name: LDY, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xAD */ &opcode{Name: LDA, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xAE */ &opcode{Name: LDX, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xAF */ &opcode{Name: LAX, Mode: absolute, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0xB0 */ &opcode{Name: BCS, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0xB1 */ &opcode{Name: LDA, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0xB2 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xB3 */ &opcode{Name: LAX, Mode: indirectIndexed, Cycle: 5, PageCycle: 1}, /* undocumented opcode */
	/* 0xB4 */ &opcode{Name: LDY, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xB5 */ &opcode{Name: LDA, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xB6 */ &opcode{Name: LDX, Mode: zeroPageY, Cycle: 4, PageCycle: 0},
	/* 0xB7 */ &opcode{Name: LAX, Mode: zeroPageY, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0xB8 */ &opcode{Name: CLV, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xB9 */ &opcode{Name: LDA, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xBA */ &opcode{Name: TSX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xBB */ &opcode{Name: LAS, Mode: absoluteY, Cycle: 4, PageCycle: 1}, /* undocumented opcode */
	/* 0xBC */ &opcode{Name: LDY, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xBD */ &opcode{Name: LDA, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xBE */ &opcode{Name: LDX, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xBF */ &opcode{Name: LAX, Mode: absoluteY, Cycle: 4, PageCycle: 1}, /* undocumented opcode */
	/* 0xC0 */ &opcode{Name: CPY, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xC1 */ &opcode{Name: CMP, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0xC2 */ &opcode{Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xC3 */ &opcode{Name: DCP, Mode: indexedIndirect, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0xC4 */ &opcode{Name: CPY, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xC5 */ &opcode{Name: CMP, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xC6 */ &opcode{Name: DEC, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0xC7 */ &opcode{Name: DCP, Mode: zeroPage, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0xC8 */ &opcode{Name: INY, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xC9 */ &opcode{Name: CMP, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xCA */ &opcode{Name: DEX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xCB */ &opcode{Name: AXS, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xCC */ &opcode{Name: CPY, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xCD */ &opcode{Name: CMP, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xCE */ &opcode{Name: DEC, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0xCF */ &opcode{Name: DCP, Mode: absolute, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0xD0 */ &opcode{Name: BNE, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0xD1 */ &opcode{Name: CMP, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0xD2 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xD3 */ &opcode{Name: DCP, Mode: indirectIndexed, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0xD4 */ &opcode{Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0xD5 */ &opcode{Name: CMP, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xD6 */ &opcode{Name: DEC, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0xD7 */ &opcode{Name: DCP, Mode: zeroPageX, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0xD8 */ &opcode{Name: CLD, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xD9 */ &opcode{Name: CMP, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xDA */ &opcode{Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xDB */ &opcode{Name: DCP, Mode: absoluteY, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0xDC */ &opcode{Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1}, /* undocumented opcode */
	/* 0xDD */ &opcode{Name: CMP, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xDE */ &opcode{Name: DEC, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0xDF */ &opcode{Name: DCP, Mode: absoluteX, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0xE0 */ &opcode{Name: CPX, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xE1 */ &opcode{Name: SBC, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0xE2 */ &opcode{Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xE3 */ &opcode{Name: ISC, Mode: indexedIndirect, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0xE4 */ &opcode{Name: CPX, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xE5 */ &opcode{Name: SBC, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xE6 */ &opcode{Name: INC, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0xE7 */ &opcode{Name: ISC, Mode: zeroPage, Cycle: 5, PageCycle: 0}, /* undocumented opcode */
	/* 0xE8 */ &opcode{Name: INX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xE9 */ &opcode{Name: SBC, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xEA */ &opcode{Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xEB */ &opcode{Name: SBC, Mode: immediate, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xEC */ &opcode{Name: CPX, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xED */ &opcode{Name: SBC, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xEE */ &opcode{Name: INC, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0xEF */ &opcode{Name: ISC, Mode: absolute, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0xF0 */ &opcode{Name: BEQ, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0xF1 */ &opcode{Name: SBC, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0xF2 */ &opcode{Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xF3 */ &opcode{Name: ISC, Mode: indirectIndexed, Cycle: 8, PageCycle: 0}, /* undocumented opcode */
	/* 0xF4 */ &opcode{Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0}, /* undocumented opcode */
	/* 0xF5 */ &opcode{Name: SBC, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xF6 */ &opcode{Name: INC, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0xF7 */ &opcode{Name: ISC, Mode: zeroPageX, Cycle: 6, PageCycle: 0}, /* undocumented opcode */
	/* 0xF8 */ &opcode{Name: SED, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xF9 */ &opcode{Name: SBC, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xFA */ &opcode{Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0}, /* undocumented opcode */
	/* 0xFB */ &opcode{Name: ISC, Mode: absoluteY, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
	/* 0xFC */ &opcode{Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1}, /* undocumented opcode */
	/* 0xFD */ &opcode{Name: SBC, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xFE */ &opcode{Name: INC, Mode: absoluteX, Cycle: 7, PageCycle: 0},
	/* 0xFF */ &opcode{Name: ISC, Mode: absoluteX, Cycle: 7, PageCycle: 0}, /* undocumented opcode */
}
