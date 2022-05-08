package cpu

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
	Name       Mnemonic
	Mode       addressingMode
	Cycle      int
	PageCycle  int
	Unofficial bool
}

var opcodeMap = []*opcode{
	/* 0x00 */ {Name: BRK, Mode: implied, Cycle: 7, PageCycle: 0},
	/* 0x01 */ {Name: ORA, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x02 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x03 */ {Name: SLO, Mode: indexedIndirect, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x04 */ {Name: NOP, Mode: zeroPage, Cycle: 3, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x05 */ {Name: ORA, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x06 */ {Name: ASL, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x07 */ {Name: SLO, Mode: zeroPage, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x08 */ {Name: PHP, Mode: implied, Cycle: 3, PageCycle: 0},
	/* 0x09 */ {Name: ORA, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x0A */ {Name: ASL, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x0B */ {Name: ANC, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x0C */ {Name: NOP, Mode: absolute, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x0D */ {Name: ORA, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x0E */ {Name: ASL, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x0F */ {Name: SLO, Mode: absolute, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x10 */ {Name: BPL, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x11 */ {Name: ORA, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x12 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x13 */ {Name: SLO, Mode: indirectIndexed_D, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x14 */ {Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x15 */ {Name: ORA, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x16 */ {Name: ASL, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x17 */ {Name: SLO, Mode: zeroPageX, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x18 */ {Name: CLC, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x19 */ {Name: ORA, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x1A */ {Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x1B */ {Name: SLO, Mode: absoluteY_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x1C */ {Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1, Unofficial: true}, /* undocumented opcode */
	/* 0x1D */ {Name: ORA, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x1E */ {Name: ASL, Mode: absoluteX_D, Cycle: 7, PageCycle: 0},
	/* 0x1F */ {Name: SLO, Mode: absoluteX_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x20 */ {Name: JSR, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x21 */ {Name: AND, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x22 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x23 */ {Name: RLA, Mode: indexedIndirect, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x24 */ {Name: BIT, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x25 */ {Name: AND, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x26 */ {Name: ROL, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x27 */ {Name: RLA, Mode: zeroPage, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x28 */ {Name: PLP, Mode: implied, Cycle: 4, PageCycle: 0},
	/* 0x29 */ {Name: AND, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x2A */ {Name: ROL, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x2B */ {Name: ANC, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x2C */ {Name: BIT, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x2D */ {Name: AND, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x2E */ {Name: ROL, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x2F */ {Name: RLA, Mode: absolute, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x30 */ {Name: BMI, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x31 */ {Name: AND, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x32 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x33 */ {Name: RLA, Mode: indirectIndexed_D, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x34 */ {Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x35 */ {Name: AND, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x36 */ {Name: ROL, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x37 */ {Name: RLA, Mode: zeroPageX, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x38 */ {Name: SEC, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x39 */ {Name: AND, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x3A */ {Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x3B */ {Name: RLA, Mode: absoluteY_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x3C */ {Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1, Unofficial: true}, /* undocumented opcode */
	/* 0x3D */ {Name: AND, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x3E */ {Name: ROL, Mode: absoluteX_D, Cycle: 7, PageCycle: 0},
	/* 0x3F */ {Name: RLA, Mode: absoluteX_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x40 */ {Name: RTI, Mode: implied, Cycle: 6, PageCycle: 0},
	/* 0x41 */ {Name: EOR, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x42 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x43 */ {Name: SRE, Mode: indexedIndirect, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x44 */ {Name: NOP, Mode: zeroPage, Cycle: 3, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x45 */ {Name: EOR, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x46 */ {Name: LSR, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x47 */ {Name: SRE, Mode: zeroPage, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x48 */ {Name: PHA, Mode: implied, Cycle: 3, PageCycle: 0},
	/* 0x49 */ {Name: EOR, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x4A */ {Name: LSR, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x4B */ {Name: ALR, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x4C */ {Name: JMP, Mode: absolute, Cycle: 3, PageCycle: 0},
	/* 0x4D */ {Name: EOR, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x4E */ {Name: LSR, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x4F */ {Name: SRE, Mode: absolute, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x50 */ {Name: BVC, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x51 */ {Name: EOR, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x52 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x53 */ {Name: SRE, Mode: indirectIndexed_D, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x54 */ {Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x55 */ {Name: EOR, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x56 */ {Name: LSR, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x57 */ {Name: SRE, Mode: zeroPageX, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x58 */ {Name: CLI, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x59 */ {Name: EOR, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x5A */ {Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x5B */ {Name: SRE, Mode: absoluteY_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x5C */ {Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1, Unofficial: true}, /* undocumented opcode */
	/* 0x5D */ {Name: EOR, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x5E */ {Name: LSR, Mode: absoluteX_D, Cycle: 7, PageCycle: 0},
	/* 0x5F */ {Name: SRE, Mode: absoluteX_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x60 */ {Name: RTS, Mode: implied, Cycle: 6, PageCycle: 0},
	/* 0x61 */ {Name: ADC, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x62 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x63 */ {Name: RRA, Mode: indexedIndirect, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x64 */ {Name: NOP, Mode: zeroPage, Cycle: 3, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x65 */ {Name: ADC, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x66 */ {Name: ROR, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0x67 */ {Name: RRA, Mode: zeroPage, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x68 */ {Name: PLA, Mode: implied, Cycle: 4, PageCycle: 0},
	/* 0x69 */ {Name: ADC, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0x6A */ {Name: ROR, Mode: accumulator, Cycle: 2, PageCycle: 0},
	/* 0x6B */ {Name: ARR, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x6C */ {Name: JMP, Mode: indirect, Cycle: 5, PageCycle: 0},
	/* 0x6D */ {Name: ADC, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x6E */ {Name: ROR, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0x6F */ {Name: RRA, Mode: absolute, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x70 */ {Name: BVS, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x71 */ {Name: ADC, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0x72 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x73 */ {Name: RRA, Mode: indirectIndexed_D, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x74 */ {Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x75 */ {Name: ADC, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x76 */ {Name: ROR, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0x77 */ {Name: RRA, Mode: zeroPageX, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x78 */ {Name: SEI, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x79 */ {Name: ADC, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0x7A */ {Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x7B */ {Name: RRA, Mode: absoluteY_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x7C */ {Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1, Unofficial: true}, /* undocumented opcode */
	/* 0x7D */ {Name: ADC, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0x7E */ {Name: ROR, Mode: absoluteX_D, Cycle: 7, PageCycle: 0},
	/* 0x7F */ {Name: RRA, Mode: absoluteX_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x80 */ {Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x81 */ {Name: STA, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0x82 */ {Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x83 */ {Name: SAX, Mode: indexedIndirect, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x84 */ {Name: STY, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x85 */ {Name: STA, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x86 */ {Name: STX, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0x87 */ {Name: SAX, Mode: zeroPage, Cycle: 3, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x88 */ {Name: DEY, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x89 */ {Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x8A */ {Name: TXA, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x8B */ {Name: XAA, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x8C */ {Name: STY, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x8D */ {Name: STA, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x8E */ {Name: STX, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0x8F */ {Name: SAX, Mode: absolute, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x90 */ {Name: BCC, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0x91 */ {Name: STA, Mode: indirectIndexed_D, Cycle: 6, PageCycle: 0},
	/* 0x92 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x93 */ {Name: AHX, Mode: indirectIndexed_D, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x94 */ {Name: STY, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x95 */ {Name: STA, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0x96 */ {Name: STX, Mode: zeroPageY, Cycle: 4, PageCycle: 0},
	/* 0x97 */ {Name: SAX, Mode: zeroPageY, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x98 */ {Name: TYA, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x99 */ {Name: STA, Mode: absoluteY_D, Cycle: 5, PageCycle: 0},
	/* 0x9A */ {Name: TXS, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0x9B */ {Name: TAS, Mode: absoluteY_D, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x9C */ {Name: SHY, Mode: absoluteX_D, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x9D */ {Name: STA, Mode: absoluteX_D, Cycle: 5, PageCycle: 0},
	/* 0x9E */ {Name: SHX, Mode: absoluteY_D, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0x9F */ {Name: AHX, Mode: absoluteY_D, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xA0 */ {Name: LDY, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xA1 */ {Name: LDA, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0xA2 */ {Name: LDX, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xA3 */ {Name: LAX, Mode: indexedIndirect, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xA4 */ {Name: LDY, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xA5 */ {Name: LDA, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xA6 */ {Name: LDX, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xA7 */ {Name: LAX, Mode: zeroPage, Cycle: 3, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xA8 */ {Name: TAY, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xA9 */ {Name: LDA, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xAA */ {Name: TAX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xAB */ {Name: LAX, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xAC */ {Name: LDY, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xAD */ {Name: LDA, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xAE */ {Name: LDX, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xAF */ {Name: LAX, Mode: absolute, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xB0 */ {Name: BCS, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0xB1 */ {Name: LDA, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0xB2 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xB3 */ {Name: LAX, Mode: indirectIndexed, Cycle: 5, PageCycle: 1, Unofficial: true}, /* undocumented opcode */
	/* 0xB4 */ {Name: LDY, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xB5 */ {Name: LDA, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xB6 */ {Name: LDX, Mode: zeroPageY, Cycle: 4, PageCycle: 0},
	/* 0xB7 */ {Name: LAX, Mode: zeroPageY, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xB8 */ {Name: CLV, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xB9 */ {Name: LDA, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xBA */ {Name: TSX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xBB */ {Name: LAS, Mode: absoluteY, Cycle: 4, PageCycle: 1, Unofficial: true}, /* undocumented opcode */
	/* 0xBC */ {Name: LDY, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xBD */ {Name: LDA, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xBE */ {Name: LDX, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xBF */ {Name: LAX, Mode: absoluteY, Cycle: 4, PageCycle: 1, Unofficial: true}, /* undocumented opcode */
	/* 0xC0 */ {Name: CPY, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xC1 */ {Name: CMP, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0xC2 */ {Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xC3 */ {Name: DCP, Mode: indexedIndirect, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xC4 */ {Name: CPY, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xC5 */ {Name: CMP, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xC6 */ {Name: DEC, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0xC7 */ {Name: DCP, Mode: zeroPage, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xC8 */ {Name: INY, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xC9 */ {Name: CMP, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xCA */ {Name: DEX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xCB */ {Name: AXS, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xCC */ {Name: CPY, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xCD */ {Name: CMP, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xCE */ {Name: DEC, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0xCF */ {Name: DCP, Mode: absolute, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xD0 */ {Name: BNE, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0xD1 */ {Name: CMP, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0xD2 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xD3 */ {Name: DCP, Mode: indirectIndexed_D, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xD4 */ {Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xD5 */ {Name: CMP, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xD6 */ {Name: DEC, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0xD7 */ {Name: DCP, Mode: zeroPageX, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xD8 */ {Name: CLD, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xD9 */ {Name: CMP, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xDA */ {Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xDB */ {Name: DCP, Mode: absoluteY_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xDC */ {Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1, Unofficial: true}, /* undocumented opcode */
	/* 0xDD */ {Name: CMP, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xDE */ {Name: DEC, Mode: absoluteX_D, Cycle: 7, PageCycle: 0},
	/* 0xDF */ {Name: DCP, Mode: absoluteX_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xE0 */ {Name: CPX, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xE1 */ {Name: SBC, Mode: indexedIndirect, Cycle: 6, PageCycle: 0},
	/* 0xE2 */ {Name: NOP, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xE3 */ {Name: ISB, Mode: indexedIndirect, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xE4 */ {Name: CPX, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xE5 */ {Name: SBC, Mode: zeroPage, Cycle: 3, PageCycle: 0},
	/* 0xE6 */ {Name: INC, Mode: zeroPage, Cycle: 5, PageCycle: 0},
	/* 0xE7 */ {Name: ISB, Mode: zeroPage, Cycle: 5, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xE8 */ {Name: INX, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xE9 */ {Name: SBC, Mode: immediate, Cycle: 2, PageCycle: 0},
	/* 0xEA */ {Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xEB */ {Name: SBC, Mode: immediate, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xEC */ {Name: CPX, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xED */ {Name: SBC, Mode: absolute, Cycle: 4, PageCycle: 0},
	/* 0xEE */ {Name: INC, Mode: absolute, Cycle: 6, PageCycle: 0},
	/* 0xEF */ {Name: ISB, Mode: absolute, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xF0 */ {Name: BEQ, Mode: relative, Cycle: 2, PageCycle: 1},
	/* 0xF1 */ {Name: SBC, Mode: indirectIndexed, Cycle: 5, PageCycle: 1},
	/* 0xF2 */ {Name: KIL, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xF3 */ {Name: ISB, Mode: indirectIndexed_D, Cycle: 8, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xF4 */ {Name: NOP, Mode: zeroPageX, Cycle: 4, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xF5 */ {Name: SBC, Mode: zeroPageX, Cycle: 4, PageCycle: 0},
	/* 0xF6 */ {Name: INC, Mode: zeroPageX, Cycle: 6, PageCycle: 0},
	/* 0xF7 */ {Name: ISB, Mode: zeroPageX, Cycle: 6, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xF8 */ {Name: SED, Mode: implied, Cycle: 2, PageCycle: 0},
	/* 0xF9 */ {Name: SBC, Mode: absoluteY, Cycle: 4, PageCycle: 1},
	/* 0xFA */ {Name: NOP, Mode: implied, Cycle: 2, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xFB */ {Name: ISB, Mode: absoluteY_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
	/* 0xFC */ {Name: NOP, Mode: absoluteX, Cycle: 4, PageCycle: 1, Unofficial: true}, /* undocumented opcode */
	/* 0xFD */ {Name: SBC, Mode: absoluteX, Cycle: 4, PageCycle: 1},
	/* 0xFE */ {Name: INC, Mode: absoluteX_D, Cycle: 7, PageCycle: 0},
	/* 0xFF */ {Name: ISB, Mode: absoluteX_D, Cycle: 7, PageCycle: 0, Unofficial: true}, /* undocumented opcode */
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
