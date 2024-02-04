package nes

import (
	"fmt"
	"sync"
)

type CPUOption func(*CPU)

type CPU struct {
	// registers
	A  byte   // Accumulator
	X  byte   // Index
	Y  byte   // Index
	PC uint16 // Program Counter
	S  byte   // Stack Pointer
	P  processorStatus

	I   *interruptLines
	bus *Bus
	t   *Trace
	mu  *sync.Mutex

	irqNeeded     bool
	prevIRQNeeded bool
	nmiNeeded     bool
	prevNMINeeded bool
	prevNMILine   interruptLineStatus
}

func NewCPU(bus *Bus, i *interruptLines, opts ...CPUOption) *CPU {
	cpu := &CPU{
		I:   i,
		bus: bus,
		t:   nil,
		mu:  &sync.Mutex{},
	}
	for _, opt := range opts {
		opt(cpu)
	}
	return cpu
}

func WithTracer(tracer *Trace) CPUOption {
	return func(cpu *CPU) {
		cpu.t = tracer
	}
}

// debug
func (cpu *CPU) FetchCycles() int {
	return cpu.bus.clock
}

/*
ref: http://users.telenet.be/kim1-6502/6502/proman.html#92

	Cycles   Address Bus   Data Bus    External Operation     Internal Operation
	1           ?           ?        Don't Care             Hold During Reset
	2         ? + 1         ?        Don't Care             First Start State
	3        0100 + SP      ?        Don't Care             Second Start State
	4        0100 + SP-1    ?        Don't Care             Third Start State
	5        0100 + SP-2    ?        Don't Care             Fourth Start State
	6        FFFC        Start PCL   Fetch First Vector
	7        FFFD        Start PCH   Fetch Second Vector    Hold PCL
	8        PCH PCL     First       Load First OP CODE
*/
func (cpu *CPU) PowerUp() {
	cpu.bus.tick(5)
	cpu.A = 0x00
	cpu.X = 0x00
	cpu.Y = 0x00
	cpu.P = processorStatus(0x34)
	cpu.S = 0xFD
	cpu.PC = cpu.read16(0xFFFC)
}

func (cpu *CPU) Reset() {
	// とりあえず今はlock取っておく
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	cpu.bus.tick(5)
	cpu.PC = cpu.read16(0xFFFC)
	cpu.P.SetInterruptDisable(true)
	cpu.S -= 3
}

func (cpu *CPU) Step() {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()

	beforeClock := cpu.bus.clock

	additionalCycle := 0

	if cpu.prevNMINeeded {
		cpu.nmiNeeded = false
		cpu.nmi()
		return
	} else if cpu.prevIRQNeeded {
		// In the internal processing of IRQ(), read is executed twice after turning the InterruptDisable ON,
		// so prevIRQNeeded is false at the beginning of the next Step().
		cpu.irq()
		return
	}

	if cpu.t != nil {
		cpu.t.SetCPURegisterA(cpu.A)
		cpu.t.SetCPURegisterX(cpu.X)
		cpu.t.SetCPURegisterY(cpu.Y)
		cpu.t.SetCPURegisterPC(cpu.PC)
		cpu.t.SetCPURegisterS(cpu.S)
		cpu.t.SetCPURegisterP(cpu.P.Byte())
	}

	opcodeByte := cpu.fetch()
	opcode := opcodeMap[opcodeByte]

	if cpu.t != nil {
		cpu.t.SetCPUOpcode(*opcode)
		cpu.t.AddCPUByteCode(opcodeByte)
	}
	addr, pageCrossed := cpu.fetchOperand(opcode)
	if pageCrossed {
		additionalCycle += opcode.PageCycle
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
		// ref: https://www.nesdev.org/6502_cpu.txt
		// It seems to be processed in the same way as LDA and LDX, etc,
		// So one dummy read operation is required in addition to fetch opcode and addressing process.
		// However, NOP also has Implied addressing cases, so ignore the case only.
		if opcode.Mode != implied {
			cpu.Read(addr) // dummy read
		}
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

// https://www.nesdev.org/wiki/CPU_interrupts#Detailed_interrupt_behavior
// > As can be deduced from above, it's really the status of the interrupt lines at the end of the second-to-last cycle that matters.
// The polling process for interrupt status occurs during φ2 of each CPU cycle, but the actual generation of the interrupt is delayed by one instruction.
func (cpu *CPU) pollInterruptLines() {
	// > The NMI input is connected to an edge detector.
	cpu.prevNMINeeded = cpu.nmiNeeded
	if cpu.I.nmiLine == interruptLineLow && cpu.prevNMILine == interruptLineHigh {
		cpu.nmiNeeded = true
	}
	cpu.prevNMILine = cpu.I.nmiLine

	// > The IRQ input is connected to a level detector.
	cpu.prevIRQNeeded = cpu.irqNeeded
	cpu.irqNeeded = false
	if cpu.I.irqLine == interruptLineLow && !cpu.P.IsInterruptDisable() {
		cpu.irqNeeded = true
	}
}

func (cpu *CPU) Read(addr uint16) byte {
	cpu.bus.RunDMAIfOccurred(true)
	cpu.bus.clock++
	cpu.bus.ppu.Step()
	cpu.bus.ppu.Step()
	cpu.bus.apu.Step()
	ret := cpu.bus.read(addr)
	cpu.bus.ppu.Step()
	cpu.pollInterruptLines()
	return ret
}

func (cpu *CPU) Write(addr uint16, val byte) {
	cpu.bus.RunDMAIfOccurred(false)
	cpu.bus.clock++
	cpu.bus.ppu.Step()
	cpu.bus.ppu.Step()
	cpu.bus.apu.Step()
	cpu.bus.write(addr, val)
	cpu.bus.ppu.Step()
	cpu.pollInterruptLines()
}

func (cpu *CPU) fetch() byte {
	v := cpu.Read(cpu.PC)
	cpu.PC++
	return v
}

func (cpu *CPU) read16(addr uint16) uint16 {
	l := cpu.Read(addr)
	h := cpu.Read(addr + 1)
	return (uint16(h) << 8) | uint16(l)
}

func (cpu *CPU) fetchOperand(op *opcode) (uint16, bool) {
	switch op.Mode {
	case absolute:
		return cpu.addressingAbsolute(op)
	case absoluteX:
		return cpu.addressingAbsoluteX(op, false)
	case absoluteX_D:
		return cpu.addressingAbsoluteX(op, true)
	case absoluteY:
		return cpu.addressingAbsoluteY(op, false)
	case absoluteY_D:
		return cpu.addressingAbsoluteY(op, true)
	case accumulator:
		return cpu.addressingAccumulator(op)
	case immediate:
		return cpu.addressingImmediate(op)
	case implied:
		return cpu.addressingImplied(op)
	case indexedIndirect:
		return cpu.addressingIndexedIndirect(op)
	case indirect:
		return cpu.addressingIndirect(op)
	case indirectIndexed:
		return cpu.addressingIndirectIndexed(op, false)
	case indirectIndexed_D:
		return cpu.addressingIndirectIndexed(op, true)
	case relative:
		return cpu.addressingRelative(op)
	case zeroPage:
		return cpu.addressingZeroPage(op)
	case zeroPageX:
		return cpu.addressingZeroPageX(op)
	case zeroPageY:
		return cpu.addressingZeroPageY(op)
	default:
		panic("unknown addressing mode")
	}
}

func (cpu *CPU) addressingAbsolute(op *opcode) (addr uint16, pageCrossed bool) {
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

func (cpu *CPU) addressingAbsoluteX(op *opcode, forceDummyRead bool) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	a := uint16(h)<<8 | uint16(l)
	addr = a + uint16(cpu.X)
	pageCrossed = pagesCross(addr, addr-uint16(cpu.X))
	if pageCrossed || forceDummyRead {
		dummyAddr := uint16(h)<<8 | ((uint16(l) + uint16(cpu.X)) & 0xFF)
		cpu.Read(dummyAddr)
	}
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(l)
		cpu.t.AddCPUByteCode(h)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X,X @ %04X = %02X", a, addr, cpu.bus.Peek(addr)))
	}
	return addr, pageCrossed
}

func (cpu *CPU) addressingAbsoluteY(op *opcode, forceDummyRead bool) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	a := uint16(h)<<8 | uint16(l)
	addr = a + uint16(cpu.Y)
	pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
	if pageCrossed || forceDummyRead {
		dummyAddr := uint16(h)<<8 | ((uint16(l) + uint16(cpu.Y)) & 0xFF)
		cpu.Read(dummyAddr)
	}
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(l)
		cpu.t.AddCPUByteCode(h)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X,Y @ %04X = %02X", a, addr, cpu.bus.Peek(addr)))
	}
	return addr, pageCrossed
}

/*
ref: https://www.nesdev.org/6502_cpu.txt

	Accumulator or implied addressing

	#  address R/W description
	--- ------- --- -----------------------------------------------
	1    PC     R  fetch opcode, increment PC
	2    PC     R  read next instruction byte (and throw it away)
*/
func (cpu *CPU) addressingAccumulator(op *opcode) (addr uint16, pageCrossed bool) {
	cpu.Read(cpu.PC) // dummy read
	if cpu.t != nil {
		cpu.t.SetCPUAddressingResult("A")
	}
	return 0, false
}
func (cpu *CPU) addressingImplied(op *opcode) (addr uint16, pageCrossed bool) {
	cpu.Read(cpu.PC) // dummy read
	return 0, false
}

func (cpu *CPU) addressingImmediate(op *opcode) (addr uint16, pageCrossed bool) {
	addr = cpu.PC
	cpu.PC++
	if cpu.t != nil {
		a := cpu.bus.Peek(addr)
		cpu.t.AddCPUByteCode(a)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("#$%02X", a))
	}
	return addr, false
}

func (cpu *CPU) addressingIndexedIndirect(op *opcode) (addr uint16, pageCrossed bool) {
	k := cpu.fetch()
	// https://www.nesdev.org/6502_cpu.txt
	// > pointer    R  read from the address, add X to it
	cpu.Read(uint16(k)) // dummy read

	a := uint16(k + cpu.X)
	b := (a & 0xFF00) | uint16(byte(a)+1)
	addr = uint16(cpu.Read(b))<<8 | uint16(cpu.Read(a))
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(k)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", k, byte(a), addr, cpu.bus.Peek(addr)))
	}
	return addr, false
}

func (cpu *CPU) addressingIndirect(op *opcode) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	a := uint16(h)<<8 | uint16(l)
	b := (a & 0xFF00) | uint16(byte(a)+1)
	addr = uint16(cpu.Read(b))<<8 | uint16(cpu.Read(a))
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(l)
		cpu.t.AddCPUByteCode(h)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("($%04X) = %04X", a, addr))
	}
	return addr, false

}
func (cpu *CPU) addressingIndirectIndexed(op *opcode, forceDummyRead bool) (addr uint16, pageCrossed bool) {
	a := uint16(cpu.fetch())
	b := (a & 0xFF00) | uint16(byte(a)+1)
	baseAddr := uint16(cpu.Read(b))<<8 | uint16(cpu.Read(a))
	addr = baseAddr + uint16(cpu.Y)
	pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
	if pageCrossed || forceDummyRead {
		h := baseAddr & 0xFF00
		l := baseAddr & 0x00FF
		dummyAddr := h | ((l + uint16(cpu.Y)) & 0xFF)
		cpu.Read(dummyAddr)
	}
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(byte(a))
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", byte(a), baseAddr, addr, cpu.bus.Peek(addr)))
	}
	return addr, pageCrossed
}
func (cpu *CPU) addressingRelative(op *opcode) (addr uint16, pageCrossed bool) {
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

func (cpu *CPU) addressingZeroPage(op *opcode) (addr uint16, pageCrossed bool) {
	a := cpu.fetch()
	addr = uint16(a)
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(a)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%02X = %02X", a, cpu.bus.Peek(addr)))
	}
	return addr, false
}

func (cpu *CPU) addressingZeroPageX(op *opcode) (addr uint16, pageCrossed bool) {
	a := cpu.fetch()
	// https://www.nesdev.org/6502_cpu.txt
	// > address   R  read from address, add index register to it
	cpu.Read(uint16(a)) // dummy read

	addr = uint16(a+cpu.X) & 0xFF
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(a)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%02X,X @ %02X = %02X", a, addr, cpu.bus.Peek(addr)))
	}
	return addr, false
}

func (cpu *CPU) addressingZeroPageY(op *opcode) (addr uint16, pageCrossed bool) {
	a := cpu.fetch()
	// https://www.nesdev.org/6502_cpu.txt
	// > address   R  read from address, add index register to it
	cpu.Read(uint16(a)) // dummy read

	addr = uint16(a+cpu.Y) & 0xFF
	if cpu.t != nil {
		cpu.t.AddCPUByteCode(a)
		cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%02X,Y @ %02X = %02X", a, addr, cpu.bus.Peek(addr)))
	}
	return addr, false
}
