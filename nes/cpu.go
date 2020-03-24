package nes

import "fmt"

type CPUCycle struct {
	stall  int
	cycles int

	noCopy noCopy
}

func NewCPUCycle() *CPUCycle {
	return &CPUCycle{}
}
func (c *CPUCycle) Stall() int {
	return c.stall
}
func (c *CPUCycle) Cycles() int {
	return c.cycles
}
func (c *CPUCycle) AddStall(x int) int {
	c.stall += x
	return c.stall
}
func (c *CPUCycle) AddCycles(x int) int {
	c.cycles += x
	return c.cycles
}

type CPU struct {
	r         *cpuRegister
	m         Memory
	cycle     *CPUCycle
	interrupt *Interrupt
	noCopy    noCopy
}

func NewCPU(mem Memory, cycle *CPUCycle, interrupt *Interrupt) *CPU {
	// ref. http://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_note-1
	return &CPU{
		r:         newCPURegister(),
		m:         mem,
		cycle:     cycle,
		interrupt: interrupt,
	}
}

// TODO: after reset
func (cpu *CPU) Reset() {
	cpu.r.PC = read16(cpu.m, 0xFFFC)
	cpu.r.P = reservedFlagMask | breakFlagMask | interruptDisableFlagMask
}

func (cpu *CPU) Step() {
	if cpu.cycle.Stall() > 0 {
		cpu.cycle.AddStall(-1)
	}

	// DEBUG
	prevPC := cpu.r.PC

	if cpu.interrupt.IsNMI() {
		nmi(cpu.r, cpu.m)
		cpu.interrupt.DeassertNMI()
	} else if cpu.interrupt.IsIRQ() {
		irq(cpu.r, cpu.m)
		cpu.interrupt.DeassertIRQ()
	}

	opcodeByte := fetch(cpu.r, cpu.m)
	opcode := opcodeMap[opcodeByte]
	if opcode.Name == UnknownMnemonic {
		panic(fmt.Sprintf("Unknown opcode: 0x%0x", opcodeByte))
	}
	additionalCycle := 0
	addr, pageCrossed := fetchOperand(cpu.r, cpu.m, opcode.Mode)
	if pageCrossed {
		additionalCycle += opcode.PageCycle
	}

	// debug
	bytes := cpu.r.PC - prevPC
	w0 := fmt.Sprintf("%02X", cpu.m.Read(prevPC))
	w1 := fmt.Sprintf("%02X", cpu.m.Read(prevPC+1))
	w2 := fmt.Sprintf("%02X", cpu.m.Read(prevPC+2))
	if bytes < 2 {
		w1 = "  "
	}
	if bytes < 3 {
		w2 = "  "
	}
	fmt.Printf("%04X  %s %s %s  %s A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d\n",
		prevPC, w0, w1, w2, opcode.Name, cpu.r.A, cpu.r.X, cpu.r.Y, cpu.r.P, cpu.r.S, (cpu.cycle.Cycles()*3)%341)

	switch opcode.Name {
	case LDA:
		lda(cpu.r, cpu.m, addr)
	case LDX:
		ldx(cpu.r, cpu.m, addr)
	case LDY:
		ldy(cpu.r, cpu.m, addr)
	case STA:
		sta(cpu.r, cpu.m, addr)
	case STX:
		stx(cpu.r, cpu.m, addr)
	case STY:
		sty(cpu.r, cpu.m, addr)
	case TAX:
		tax(cpu.r)
	case TAY:
		tay(cpu.r)
	case TSX:
		tsx(cpu.r)
	case TXA:
		txa(cpu.r)
	case TXS:
		txs(cpu.r)
	case TYA:
		tya(cpu.r)
	case ADC:
		adc(cpu.r, cpu.m, addr)
	case AND:
		and(cpu.r, cpu.m, addr)
	case ASL:
		if opcode.Mode == accumulator {
			aslAcc(cpu.r)
		} else {
			asl(cpu.r, cpu.m, addr)
		}
	case BIT:
		bit(cpu.r, cpu.m, addr)
	case CMP:
		cmp(cpu.r, cpu.m, addr)
	case CPX:
		cpx(cpu.r, cpu.m, addr)
	case CPY:
		cpy(cpu.r, cpu.m, addr)
	case DEC:
		dec(cpu.r, cpu.m, addr)
	case DEX:
		dex(cpu.r)
	case DEY:
		dey(cpu.r)
	case EOR:
		eor(cpu.r, cpu.m, addr)
	case INC:
		inc(cpu.r, cpu.m, addr)
	case INX:
		inx(cpu.r)
	case INY:
		iny(cpu.r)
	case LSR:
		if opcode.Mode == accumulator {
			lsrAcc(cpu.r)
		} else {
			lsr(cpu.r, cpu.m, addr)
		}
	case ORA:
		ora(cpu.r, cpu.m, addr)
	case ROL:
		if opcode.Mode == accumulator {
			rolAcc(cpu.r)
		} else {
			rol(cpu.r, cpu.m, addr)
		}
	case ROR:
		if opcode.Mode == accumulator {
			rorAcc(cpu.r)
		} else {
			ror(cpu.r, cpu.m, addr)
		}
	case SBC:
		sbc(cpu.r, cpu.m, addr)
	case PHA:
		pha(cpu.r, cpu.m)
	case PHP:
		php(cpu.r, cpu.m)
	case PLA:
		pla(cpu.r, cpu.m)
	case PLP:
		plp(cpu.r, cpu.m)
	case JMP:
		jmp(cpu.r, addr)
	case JSR:
		jsr(cpu.r, cpu.m, addr)
	case RTS:
		rts(cpu.r, cpu.m)
	case RTI:
		rti(cpu.r, cpu.m)
	case BCC:
		additionalCycle += bcc(cpu.r, addr)
	case BCS:
		additionalCycle += bcs(cpu.r, addr)
	case BEQ:
		additionalCycle += beq(cpu.r, addr)
	case BMI:
		additionalCycle += bmi(cpu.r, addr)
	case BNE:
		additionalCycle += bne(cpu.r, addr)
	case BPL:
		additionalCycle += bpl(cpu.r, addr)
	case BVC:
		additionalCycle += bvc(cpu.r, addr)
	case BVS:
		additionalCycle += bvs(cpu.r, addr)
	case CLC:
		clc(cpu.r)
	case CLD:
		cld(cpu.r)
	case CLI:
		cli(cpu.r)
	case CLV:
		clv(cpu.r)
	case SEC:
		sec(cpu.r)
	case SED:
		sed(cpu.r)
	case SEI:
		sei(cpu.r)
	case BRK:
		brk(cpu.r, cpu.m)
	case NOP:
	// case KIL:
	// case SLO:
	// case ANC:
	// case RLA:
	// case SRE:
	// case ALR:
	// case RRA:
	// case ARR:
	case SAX:
		sax(cpu.r, cpu.m, addr)
	// case XAA:
	// case AHX:
	// case TAS:
	// case SHY:
	// case SHX:
	case LAX:
		lax(cpu.r, cpu.m, addr)
	// case LAS:
	case DCP:
		dcp(cpu.r, cpu.m, addr)
	// case AXS:
	// case ISC:
	default:
		panic("Unable to reach here")
	}

	cpu.cycle.AddCycles(opcode.Cycle + additionalCycle)
}

func fetch(r *cpuRegister, m MemoryReader) byte {
	v := m.Read(r.PC)
	r.PC++
	return v
}

func fetch16(r *cpuRegister, m MemoryReader) uint16 {
	l := fetch(r, m)
	h := fetch(r, m)
	return uint16(h)<<8 | uint16(l)
}

func fetchOperand(r *cpuRegister, m MemoryReader, mode addressingMode) (addr uint16, pageCrossed bool) {
	pageCrossed = false

	switch mode {
	case absolute:
		addr = fetch16(r, m)
	case absoluteX:
		addr = fetch16(r, m) + uint16(r.X)
		pageCrossed = pagesCross(addr, addr-uint16(r.X))
	case absoluteY:
		addr = fetch16(r, m) + uint16(r.Y)
		pageCrossed = pagesCross(addr, addr-uint16(r.Y))
	case accumulator:
		addr = 0
	case immediate:
		addr = r.PC
		r.PC++
	case implied:
		addr = 0
	case indexedIndirect:
		a := uint16(fetch(r, m) + r.X)
		b := (a & 0xFF00) | uint16(byte(a)+1)
		addr = uint16(m.Read(b))<<8 | uint16(m.Read(a))
	case indirect:
		a := fetch16(r, m)
		b := (a & 0xFF00) | uint16(byte(a)+1)
		addr = uint16(m.Read(b))<<8 | uint16(m.Read(a))
	case indirectIndexed:
		a := uint16(fetch(r, m))
		b := (a & 0xFF00) | uint16(byte(a)+1)
		baseAddr := uint16(m.Read(b))<<8 | uint16(m.Read(a))
		addr = baseAddr + uint16(r.Y)
		pageCrossed = pagesCross(addr, addr-uint16(r.Y))
	case relative:
		offset := uint16(fetch(r, m))
		if offset < 0x80 {
			addr = r.PC + offset
		} else {
			addr = r.PC + offset - 0x100
		}
	case zeroPage:
		addr = uint16(fetch(r, m))
	case zeroPageX:
		addr = uint16(fetch(r, m)+r.X) & 0xFF
	case zeroPageY:
		addr = uint16(fetch(r, m)+r.Y) & 0xFF
	}

	return addr, pageCrossed
}
