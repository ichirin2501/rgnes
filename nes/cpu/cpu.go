package cpu

import (
	"fmt"
)

type Option func(*CPU)

type CPU struct {
	A  byte   // Accumulator
	X  byte   // Index
	Y  byte   // Index
	PC uint16 // Program Counter
	S  byte   // Stack Pointer
	P  StatusRegister

	I   *Interrupter
	bus *Bus
	t   *Trace
}

func New(bus *Bus, i *Interrupter, opts ...Option) *CPU {
	cpu := &CPU{
		A:  0x00,
		X:  0x00,
		Y:  0x00,
		PC: 0x8000,
		S:  0xFD,
		P:  StatusRegister(0x24),

		I:   i,
		bus: bus,
		t:   nil,
	}
	for _, opt := range opts {
		opt(cpu)
	}
	return cpu
}

func WithTracer(tracer *Trace) Option {
	return func(cpu *CPU) {
		cpu.t = tracer
	}
}

// debug
func (cpu *CPU) FetchCycles() int {
	return cpu.bus.clock
}

// TODO: after reset
func (cpu *CPU) Reset() {
	cpu.PC = cpu.read16(0xFFFC)
	cpu.P = StatusRegister(0x24)
	// adjust
	cpu.bus.tick(5)
}

func (cpu *CPU) Step() {
	beforeClock := cpu.bus.clock

	additionalCycle := 0

	switch cpu.I.interrupt {
	case InterruptNMI:
		if !cpu.I.delayNMI {
			cpu.nmi()
			cpu.I.interrupt = InterruptNone
			// adjust
			cpu.bus.tick(2)
			return
		}
	case InterruptIRQ:
		cpu.irq()
		cpu.I.interrupt = InterruptNone
		// adjust
		cpu.bus.tick(2)
		return
	}

	cpu.I.delayNMI = false

	if cpu.t != nil {
		cpu.t.SetCPURegisterA(cpu.A)
		cpu.t.SetCPURegisterX(cpu.X)
		cpu.t.SetCPURegisterY(cpu.Y)
		cpu.t.SetCPURegisterPC(cpu.PC)
		cpu.t.SetCPURegisterS(cpu.S)
		cpu.t.SetCPURegisterP(cpu.P.Byte())
	}

	//fmt.Printf("[debug] before fetch(): cpu.PC = 0x%04X  _____ %s\n", cpu.PC, cpu.t.NESTestString())

	opcodeByte := cpu.fetch()
	opcode := opcodeMap[opcodeByte]

	//fmt.Printf("[debug] after fetch(): %s\n", cpu.t.NESTestString())
	//fmt.Println(opcode)

	if cpu.t != nil {
		cpu.t.SetCPUOpcode(*opcode)
		cpu.t.AddCPUByteCode(opcodeByte)
	}
	addr, pageCrossed := cpu.fetchOperand(opcode)
	if pageCrossed {
		additionalCycle += opcode.PageCycle
		// fetchOperand()内のdummyReadでppu 3step回ってるからここで回す必要はない
		// cpu.cycles += opcode.PageCycle
		// for i := 0; i < opcode.PageCycle*3; i++ {
		// 	cpu.bus.ppu.Step()
		// }
	}

	switch opcode.Name {
	case LDA:
		cpu.lda(addr)
	case LDX:
		cpu.ldx(addr)
	case LDY:
		cpu.ldy(addr)
	case STA:
		cpu.sta(addr)
	case STX:
		cpu.stx(addr)
	case STY:
		cpu.sty(addr)
	case TAX:
		cpu.tax()
	case TAY:
		cpu.tay()
	case TSX:
		cpu.tsx()
	case TXA:
		cpu.txa()
	case TXS:
		cpu.txs()
	case TYA:
		cpu.tya()
	case ADC:
		cpu.adc(addr)
	case AND:
		cpu.and(addr)
	case ASL:
		if opcode.Mode == accumulator {
			cpu.aslAcc()
		} else {
			cpu.asl(addr)
		}
	case BIT:
		cpu.bit(addr)
	case CMP:
		cpu.cmp(addr)
	case CPX:
		cpu.cpx(addr)
	case CPY:
		cpu.cpy(addr)
	case DEC:
		cpu.dec(addr)
	case DEX:
		cpu.dex()
	case DEY:
		cpu.dey()
	case EOR:
		cpu.eor(addr)
	case INC:
		cpu.inc(addr)
	case INX:
		cpu.inx()
	case INY:
		cpu.iny()
	case LSR:
		if opcode.Mode == accumulator {
			cpu.lsrAcc()
		} else {
			cpu.lsr(addr)
		}
	case ORA:
		cpu.ora(addr)
	case ROL:
		if opcode.Mode == accumulator {
			cpu.rolAcc()
		} else {
			cpu.rol(addr)
		}
	case ROR:
		if opcode.Mode == accumulator {
			cpu.rorAcc()
		} else {
			cpu.ror(addr)
		}
	case SBC:
		cpu.sbc(addr)
	case PHA:
		cpu.pha()
	case PHP:
		cpu.php()
	case PLA:
		cpu.pla()
	case PLP:
		cpu.plp()
	case JMP:
		cpu.jmp(addr)
	case JSR:
		cpu.jsr(addr)
	case RTS:
		cpu.rts()
	case RTI:
		cpu.rti()
	case BCC:
		additionalCycle += cpu.bcc(addr)
	case BCS:
		additionalCycle += cpu.bcs(addr)
	case BEQ:
		additionalCycle += cpu.beq(addr)
	case BMI:
		additionalCycle += cpu.bmi(addr)
	case BNE:
		additionalCycle += cpu.bne(addr)
	case BPL:
		additionalCycle += cpu.bpl(addr)
	case BVC:
		additionalCycle += cpu.bvc(addr)
	case BVS:
		additionalCycle += cpu.bvs(addr)
	case CLC:
		cpu.clc()
	case CLD:
		cpu.cld()
	case CLI:
		cpu.cli()
	case CLV:
		cpu.clv()
	case SEC:
		cpu.sec()
	case SED:
		cpu.sed()
	case SEI:
		cpu.sei()
	case BRK:
		cpu.brk()
	case NOP:
	// case KIL:
	// 	// TODO
	// 	cpu.kil()
	case SLO:
		cpu.slo(addr)
	case ANC:
		cpu.anc(addr)
	case RLA:
		cpu.rla(addr)
	case SRE:
		cpu.sre(addr)
	case ALR:
		cpu.alr(addr)
	case RRA:
		cpu.rra(addr)
	case ARR:
		cpu.arr(addr)
	case SAX:
		cpu.sax(addr)
	// case XAA:
	// case AHX:
	// case TAS:
	case SHY:
		cpu.shy(addr)
	case SHX:
		cpu.shx(addr)
	case LAX:
		cpu.lax(addr)
	// case LAS:
	case DCP:
		cpu.dcp(addr)
	case AXS:
		cpu.axs(addr)
	case ISB:
		cpu.isb(addr)
	default:
		panic(fmt.Sprintf("Unable to reach: opcode.Name:%s", opcode.Name))
	}

	afterClock := cpu.bus.clock
	if (opcode.Cycle+additionalCycle)-(afterClock-beforeClock) > 0 {
		t := (opcode.Cycle + additionalCycle) - (afterClock - beforeClock)
		cpu.bus.tick(t)
	}

	if (opcode.Cycle+additionalCycle)-(afterClock-beforeClock) < 0 {
		fmt.Printf("panic: %02X\t%s\t%s\tcycle:%d\tclock:%d\tdiff:%d\tunoff:%v\n",
			opcodeByte,
			opcode.Name,
			opcode.Mode,
			opcode.Cycle+additionalCycle,
			afterClock-beforeClock,
			(opcode.Cycle+additionalCycle)-(afterClock-beforeClock),
			opcode.Unofficial,
		)
		panic("wryyyyyyyyyyyyyy")
	}
	// if (opcode.Cycle+additionalCycle)-(afterClock-beforeClock) > 0 {
	// 	fmt.Printf("%02X\t%s\t%s\tcycle:%d\tclock:%d\tdiff:%d\tunoff:%v\n",
	// 		opcodeByte,
	// 		opcode.Name,
	// 		opcode.Mode,
	// 		opcode.Cycle+additionalCycle,
	// 		afterClock-beforeClock,
	// 		(opcode.Cycle+additionalCycle)-(afterClock-beforeClock),
	// 		opcode.Unofficial,
	// 	)
	// }
}

func (cpu *CPU) fetch() byte {
	v := cpu.bus.Read(cpu.PC)
	cpu.PC++
	return v
}

func (cpu *CPU) read16(addr uint16) uint16 {
	l := cpu.bus.Read(addr)
	h := cpu.bus.Read(addr + 1)
	return (uint16(h) << 8) | uint16(l)
}

func (cpu *CPU) fetchOperand(op *opcode) (uint16, bool) {
	switch op.Mode {
	case absolute:
		return cpu.AddressingAbsolute(op)
	case absoluteX:
		return cpu.AddressingAbsoluteX(op, false)
	case absoluteX_D:
		return cpu.AddressingAbsoluteX(op, true)
	case absoluteY:
		return cpu.AddressingAbsoluteY(op, false)
	case absoluteY_D:
		return cpu.AddressingAbsoluteY(op, true)
	case accumulator:
		return cpu.AddressingAccumulator(op)
	case immediate:
		return cpu.AddressingImmediate(op)
	case implied:
		return cpu.AddressingImplied(op)
	case indexedIndirect:
		return cpu.AddressingIndexedIndirect(op)
	case indirect:
		return cpu.AddressingIndirect(op)
	case indirectIndexed:
		return cpu.AddressingIndirectIndexed(op, false)
	case indirectIndexed_D:
		return cpu.AddressingIndirectIndexed(op, true)
	case relative:
		return cpu.AddressingRelative(op)
	case zeroPage:
		return cpu.AddressingZeroPage(op)
	case zeroPageX:
		return cpu.AddressingZeroPageX(op)
	case zeroPageY:
		return cpu.AddressingZeroPageY(op)
	default:
		panic("unknown addressing mode")
	}
}

func (cpu *CPU) AddressingAbsolute(op *opcode) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	addr = uint16(h)<<8 | uint16(l)
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(l)
		cpu.t.AddCPUByteCode(h)
		if op.Name == JMP || op.Name == JSR {
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X", addr))
		} else {
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X = %02X", addr, cpu.bus.Peek(addr)))
		}
	}
	return addr, false
}

func (cpu *CPU) AddressingAbsoluteX(op *opcode, forceDummyRead bool) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	a := uint16(h)<<8 | uint16(l)
	addr = a + uint16(cpu.X)
	pageCrossed = pagesCross(addr, addr-uint16(cpu.X))
	if pageCrossed || forceDummyRead {
		dummyAddr := uint16(h)<<8 | ((uint16(l) + uint16(cpu.X)) & 0xFF)
		cpu.bus.Read(dummyAddr)
	}
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(l)
		cpu.t.AddCPUByteCode(h)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X,X @ %04X = %02X", a, addr, cpu.bus.Peek(addr)))
	}
	return addr, pageCrossed
}

func (cpu *CPU) AddressingAbsoluteY(op *opcode, forceDummyRead bool) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	a := uint16(h)<<8 | uint16(l)
	addr = a + uint16(cpu.Y)
	pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
	if pageCrossed || forceDummyRead {
		dummyAddr := uint16(h)<<8 | ((uint16(l) + uint16(cpu.Y)) & 0xFF)
		cpu.bus.Read(dummyAddr)
	}
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(l)
		cpu.t.AddCPUByteCode(h)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X,Y @ %04X = %02X", a, addr, cpu.bus.Peek(addr)))
	}
	return addr, pageCrossed
}

// https://www.nesdev.org/6502_cpu.txt
// Accumulator or implied addressing
//
// #  address R/W description
// --- ------- --- -----------------------------------------------
// 1    PC     R  fetch opcode, increment PC
// 2    PC     R  read next instruction byte (and throw it away)
func (cpu *CPU) AddressingAccumulator(op *opcode) (addr uint16, pageCrossed bool) {
	cpu.bus.Read(cpu.PC) // dummy read
	if cpu.t != nil {
		cpu.t.SetCPUAddressingResult("A")
	}
	return 0, false
}
func (cpu *CPU) AddressingImplied(op *opcode) (addr uint16, pageCrossed bool) {
	cpu.bus.Read(cpu.PC) // dummy read
	return 0, false
}

func (cpu *CPU) AddressingImmediate(op *opcode) (addr uint16, pageCrossed bool) {
	addr = cpu.PC
	cpu.PC++
	if cpu.t != nil {
		a := cpu.bus.Peek(addr)
		cpu.t.AddCPUByteCode(a)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("#$%02X", a))
	}
	return addr, false
}

func (cpu *CPU) AddressingIndexedIndirect(op *opcode) (addr uint16, pageCrossed bool) {
	k := cpu.fetch()
	// https://www.nesdev.org/6502_cpu.txt
	// > pointer    R  read from the address, add X to it
	cpu.bus.Read(uint16(k)) // dummy read

	a := uint16(k + cpu.X)
	b := (a & 0xFF00) | uint16(byte(a)+1)
	addr = uint16(cpu.bus.Read(b))<<8 | uint16(cpu.bus.Read(a))
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(k)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", k, byte(a), addr, cpu.bus.Peek(addr)))
	}
	return addr, false
}

func (cpu *CPU) AddressingIndirect(op *opcode) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	a := uint16(h)<<8 | uint16(l)
	b := (a & 0xFF00) | uint16(byte(a)+1)
	addr = uint16(cpu.bus.Read(b))<<8 | uint16(cpu.bus.Read(a))
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(l)
		cpu.t.AddCPUByteCode(h)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("($%04X) = %04X", a, addr))
	}
	return addr, false

}
func (cpu *CPU) AddressingIndirectIndexed(op *opcode, forceDummyRead bool) (addr uint16, pageCrossed bool) {
	a := uint16(cpu.fetch())
	b := (a & 0xFF00) | uint16(byte(a)+1)
	baseAddr := uint16(cpu.bus.Read(b))<<8 | uint16(cpu.bus.Read(a))
	addr = baseAddr + uint16(cpu.Y)
	pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
	if pageCrossed || forceDummyRead {
		h := baseAddr & 0xFF00
		l := baseAddr & 0x00FF
		dummyAddr := h | ((l + uint16(cpu.Y)) & 0xFF)
		cpu.bus.Read(dummyAddr)
	}
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(byte(a))
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", byte(a), baseAddr, addr, cpu.bus.Peek(addr)))
	}
	return addr, pageCrossed
}
func (cpu *CPU) AddressingRelative(op *opcode) (addr uint16, pageCrossed bool) {
	offset := uint16(cpu.fetch())
	if offset < 0x80 {
		addr = cpu.PC + offset
	} else {
		addr = cpu.PC + offset - 0x100
	}
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(byte(offset))
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X", addr))
	}
	return addr, false
}

func (cpu *CPU) AddressingZeroPage(op *opcode) (addr uint16, pageCrossed bool) {
	a := cpu.fetch()
	addr = uint16(a)
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(a)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%02X = %02X", a, cpu.bus.Peek(addr)))
	}
	return addr, false
}

func (cpu *CPU) AddressingZeroPageX(op *opcode) (addr uint16, pageCrossed bool) {
	a := cpu.fetch()
	// https://www.nesdev.org/6502_cpu.txt
	// > address   R  read from address, add index register to it
	cpu.bus.Read(uint16(a)) // dummy read

	addr = uint16(a+cpu.X) & 0xFF
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(a)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%02X,X @ %02X = %02X", a, addr, cpu.bus.Peek(addr)))
	}
	return addr, false
}

func (cpu *CPU) AddressingZeroPageY(op *opcode) (addr uint16, pageCrossed bool) {
	a := cpu.fetch()
	// https://www.nesdev.org/6502_cpu.txt
	// > address   R  read from address, add index register to it
	cpu.bus.Read(uint16(a)) // dummy read

	addr = uint16(a+cpu.Y) & 0xFF
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(a)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%02X,Y @ %02X = %02X", a, addr, cpu.bus.Peek(addr)))
	}
	return addr, false
}
