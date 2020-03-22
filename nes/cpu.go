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
	cycle     *CPUCycle
	interrupt *Interrupt
	memory    Memory
	noCopy    noCopy
}

func NewCPU(mem Memory, cycle *CPUCycle, interrupt *Interrupt) *CPU {
	// ref. http://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_note-1
	return &CPU{
		r:         newCPURegister(),
		cycle:     cycle,
		interrupt: interrupt,
		memory:    mem,
	}
}

// TODO: after reset
func (cpu *CPU) reset() {
	cpu.r.PC = read16(cpu.memory, 0xFFFC)
	cpu.r.P = reservedFlagMask | breakFlagMask | interruptDisableFlagMask
}

func (cpu *CPU) Step() int {
	if cpu.cycle.Stall() > 0 {
		cpu.cycle.AddStall(-1)
		return 1
	}

	if cpu.interrupt.IsNMI() {
		nmi(cpu.r, cpu.memory)
		cpu.interrupt.DeassertNMI()
	} else if cpu.interrupt.IsIRQ() {
		irq(cpu.r, cpu.memory)
		cpu.interrupt.DeassertIRQ()
	}

	opcodeByte := fetch(cpu.r, cpu.memory)
	opcode, ok := opcodeMap[opcodeByte]
	if !ok {
		panic(fmt.Sprintf("Unknown opcode: %0x", opcodeByte))
	}
	additionalCycle := 0
	addr, pageCrossed := fetchOperand(cpu.r, cpu.memory, opcode.Mode)
	if pageCrossed {
		additionalCycle++
	}

	switch opcode.Name {
	case LDA:
		lda(cpu.r, cpu.memory, addr)
	case LDX:
		ldx(cpu.r, cpu.memory, addr)
	case LDY:
		ldy(cpu.r, cpu.memory, addr)
	case STA:
		sta(cpu.r, cpu.memory, addr)
	case STX:
		stx(cpu.r, cpu.memory, addr)
	case STY:
		sty(cpu.r, cpu.memory, addr)
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
		adc(cpu.r, cpu.memory, addr)
	case AND:
		and(cpu.r, cpu.memory, addr)
	case ASL:
		if opcode.Mode == accumulator {
			aslAcc(cpu.r)
		} else {
			asl(cpu.r, cpu.memory, addr)
		}
	case BIT:
		bit(cpu.r, cpu.memory, addr)
	case CMP:
		cmp(cpu.r, cpu.memory, addr)
	case CPX:
		cpx(cpu.r, cpu.memory, addr)
	case CPY:
		cpy(cpu.r, cpu.memory, addr)
	case DEC:
		dec(cpu.r, cpu.memory, addr)
	case DEX:
		dex(cpu.r)
	case DEY:
		dey(cpu.r)
	case EOR:
		eor(cpu.r, cpu.memory, addr)
	case INC:
		inc(cpu.r, cpu.memory, addr)
	case INX:
		inx(cpu.r)
	case INY:
		iny(cpu.r)
	case LSR:
		if opcode.Mode == accumulator {
			lsrAcc(cpu.r)
		} else {
			lsr(cpu.r, cpu.memory, addr)
		}
	case ORA:
		ora(cpu.r, cpu.memory, addr)
	case ROL:
		if opcode.Mode == accumulator {
			rolAcc(cpu.r)
		} else {
			rol(cpu.r, cpu.memory, addr)
		}
	case ROR:
		if opcode.Mode == accumulator {
			rorAcc(cpu.r)
		} else {
			ror(cpu.r, cpu.memory, addr)
		}
	case SBC:
		sbc(cpu.r, cpu.memory, addr)
	case PHA:
		pha(cpu.r, cpu.memory)
	case PHP:
		php(cpu.r, cpu.memory)
	case PLA:
		pla(cpu.r, cpu.memory)
	case PLP:
		plp(cpu.r, cpu.memory)
	case JMP:
		jmp(cpu.r, addr)
	case JSR:
		jsr(cpu.r, cpu.memory, addr)
	case RTS:
		rts(cpu.r, cpu.memory)
	case RTI:
		rti(cpu.r, cpu.memory)
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
		brk(cpu.r, cpu.memory)
	case NOP:
	default:
		panic("Unable to reach here")
	}

	return opcode.Cycle + additionalCycle
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
		baseAddr := uint16((fetch(r, m) + r.X) & 0xFF)
		addr = uint16(m.Read((baseAddr+1)&0xFF))<<8 | uint16(m.Read(baseAddr))
	case indirect:
		baseAddr := fetch16(r, m)
		addr = uint16(m.Read((baseAddr+1)&0xFF))<<8 | uint16(m.Read(baseAddr))
	case indirectIndexed:
		baseAddr := uint16(fetch(r, m))
		baseAddr2 := uint16(m.Read((baseAddr+1)&0xFF))<<8 | uint16(m.Read(baseAddr))
		addr = baseAddr2 + uint16(r.Y)
		pageCrossed = pagesCross(addr, addr-uint16(r.Y))
	case relative:
		offset := uint16(fetch(r, m))
		if offset < 0x80 {
			addr = r.PC + 2 + offset
		} else {
			addr = r.PC + 2 + offset - 0x100
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
