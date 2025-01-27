package nes

import (
	"fmt"
	"sync"
)

const (
	// ref: https://www.nesdev.org/wiki/CPU#Frequencies
	// NTSC
	CPUClockFrequency = 1789773
)

type pendingInterruptType int

const (
	pendingInterruptNone pendingInterruptType = iota
	pendingInterruptIRQ
	pendingInterruptNMI
)

type cpu struct {
	// registers
	A  byte   // Accumulator
	X  byte   // Index
	Y  byte   // Index
	PC uint16 // Program Counter
	S  byte   // Stack Pointer
	P  processorStatus

	nmiLine *nmiInterruptLine
	irqLine *irqInterruptLine
	bus     *cpuBus
	tracer  *tracer
	mu      *sync.Mutex

	irqTriggered bool
	irqSignal    bool
	nmiTriggered bool
	nmiSignal    bool
	prevNMILine  nmiInterruptLine

	pendingInterrupt pendingInterruptType

	// https://www.nesdev.org/wiki/CPU_interrupts#Detailed_interrupt_behavior
	// > The interrupt sequences themselves do not perform interrupt polling,
	// > meaning at least one instruction from the interrupt handler will execute before another interrupt is serviced.
	interrupting bool
}

func newCPU(bus *cpuBus, nmiLine *nmiInterruptLine, irqLine *irqInterruptLine, tracer *tracer) *cpu {
	cpu := &cpu{
		nmiLine: nmiLine,
		irqLine: irqLine,
		bus:     bus,
		tracer:  tracer,
		mu:      &sync.Mutex{},
	}
	return cpu
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
func (cpu *cpu) powerUp() {
	cpu.bus.tick(5)
	cpu.A = 0x00
	cpu.X = 0x00
	cpu.Y = 0x00
	cpu.P = processorStatus(0x34)
	cpu.S = 0xFD
	cpu.PC = cpu.read16(0xFFFC)
}

func (cpu *cpu) reset() {
	// とりあえず今はlock取っておく
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	cpu.bus.tick(5)
	cpu.PC = cpu.read16(0xFFFC)
	cpu.P.setInterruptDisable(true)
	cpu.S -= 3
}

func (cpu *cpu) step() {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()

	if cpu.tracer != nil {
		cpu.tracer.reset()
		cpu.tracer.setPPUX(uint16(cpu.bus.ppu.cycle))
		cpu.tracer.setPPUY(uint16(cpu.bus.ppu.scanline))
		cpu.tracer.setCPURegisters(cpu)
	}

	beforeClock := cpu.bus.clock

	additionalCycle := 0

	if cpu.pendingInterrupt == pendingInterruptNMI {
		cpu.nmi()
		return
	} else if cpu.pendingInterrupt == pendingInterruptIRQ {
		cpu.irq()
		return
	}

	opcodeByte := cpu.fetch()
	opcode := opcodeMap[opcodeByte]

	if cpu.tracer != nil {
		cpu.tracer.setCPUOpcode(*opcode)
		cpu.tracer.addCPUByteCode(opcodeByte)
	}

	addr, pageCrossed := cpu.fetchOperand(opcode)
	if pageCrossed {
		additionalCycle += opcode.pageCycle
	}

	switch opcode.name {
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
		if opcode.mode == accumulator {
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
		if opcode.mode == accumulator {
			cpu.lsrAcc()
		} else {
			cpu.lsr(addr)
		}
	case ORA:
		cpu.ora(addr)
	case ROL:
		if opcode.mode == accumulator {
			cpu.rolAcc()
		} else {
			cpu.rol(addr)
		}
	case ROR:
		if opcode.mode == accumulator {
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
		if opcode.mode != implied {
			cpu.read(addr) // dummy read
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
		panic(fmt.Sprintf("Unable to reach: opcode.Name:%s", opcode.name))
	}

	afterClock := cpu.bus.clock
	if (opcode.cycle+additionalCycle)-(afterClock-beforeClock) > 0 {
		t := (opcode.cycle + additionalCycle) - (afterClock - beforeClock)
		cpu.bus.tick(t)
	}

	if (opcode.cycle+additionalCycle)-(afterClock-beforeClock) < 0 {
		fmt.Printf("panic: %02X\t%s\t%s\tcycle:%d\tclock:%d\tdiff:%d\tunoff:%v\n",
			opcodeByte,
			opcode.name,
			opcode.mode,
			opcode.cycle+additionalCycle,
			afterClock-beforeClock,
			(opcode.cycle+additionalCycle)-(afterClock-beforeClock),
			opcode.unofficial,
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

	if cpu.tracer != nil {
		cpu.tracer.print()
	}
}

func (cpu *cpu) pollInterruptSignals() {
	// https://www.nesdev.org/wiki/CPU_interrupts#Detailed_interrupt_behavior
	// The internal signals of NMI/IRQ inputs are detected in the φ1 of each CPU cycle, so I emulate that.
	// Once NMI internal signal is turned ON, it will not turn OFF until NMI is executed.
	if cpu.nmiTriggered {
		cpu.nmiSignal = true
	}
	cpu.irqSignal = cpu.irqTriggered
}

// https://www.nesdev.org/wiki/CPU_interrupts#Detailed_interrupt_behavior
// > As can be deduced from above, it's really the status of the interrupt lines at the end of the second-to-last cycle that matters.
// The polling process for interrupt status occurs during φ2 of each CPU cycle, but the actual generation of the interrupt is delayed by one instruction.
func (cpu *cpu) pollInterrupts() {
	// poll interrupt lines
	// The NMI input is connected to an edge detector
	cpu.nmiTriggered = false
	if cpu.nmiLine.isLow() && cpu.prevNMILine.isHigh() {
		cpu.nmiTriggered = true
	}
	cpu.prevNMILine = *cpu.nmiLine

	// The IRQ input is connected to a level detector
	cpu.irqTriggered = false
	if cpu.irqLine.isLow() {
		cpu.irqTriggered = true
	}

	// poll interrupt events
	cpu.pendingInterrupt = pendingInterruptNone
	if !cpu.interrupting {
		if cpu.nmiSignal {
			cpu.pendingInterrupt = pendingInterruptNMI
		} else if cpu.irqSignal && !cpu.P.isInterruptDisable() {
			cpu.pendingInterrupt = pendingInterruptIRQ
		}
	}
}

func (cpu *cpu) clearNMIInterruptState() {
	cpu.nmiSignal = false
	cpu.nmiTriggered = false
}

func (cpu *cpu) read(addr uint16) byte {
	cpu.bus.runDMAIfOccurred(true)
	cpu.pollInterruptSignals()
	cpu.bus.clock++
	cpu.bus.ppu.step()
	cpu.bus.ppu.step()
	cpu.bus.apu.step()
	ret := cpu.bus.read(addr)
	cpu.bus.ppu.step()
	cpu.pollInterrupts()
	return ret
}

func (cpu *cpu) write(addr uint16, val byte) {
	cpu.bus.runDMAIfOccurred(false)
	cpu.pollInterruptSignals()
	cpu.bus.clock++
	cpu.bus.ppu.step()
	cpu.bus.ppu.step()
	cpu.bus.apu.step()
	cpu.bus.write(addr, val)
	cpu.bus.ppu.step()
	cpu.pollInterrupts()
}

func (cpu *cpu) fetch() byte {
	v := cpu.read(cpu.PC)
	cpu.PC++
	return v
}

func (cpu *cpu) read16(addr uint16) uint16 {
	l := cpu.read(addr)
	h := cpu.read(addr + 1)
	return (uint16(h) << 8) | uint16(l)
}

func (cpu *cpu) fetchOperand(op *opcode) (uint16, bool) {
	switch op.mode {
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

func (cpu *cpu) addressingAbsolute(op *opcode) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	addr = uint16(h)<<8 | uint16(l)

	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(l)
		cpu.tracer.addCPUByteCode(h)
		if op.name == JMP || op.name == JSR {
			cpu.tracer.setCPUAddressingResult(fmt.Sprintf("$%04X", addr))
		} else {
			cpu.tracer.setCPUAddressingResult(fmt.Sprintf("$%04X = %02X", addr, cpu.bus.peek(addr)))
		}
	}

	return addr, false
}

func (cpu *cpu) addressingAbsoluteX(_ *opcode, forceDummyRead bool) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	a := uint16(h)<<8 | uint16(l)
	addr = a + uint16(cpu.X)
	pageCrossed = pagesCross(addr, addr-uint16(cpu.X))
	if pageCrossed || forceDummyRead {
		dummyAddr := uint16(h)<<8 | ((uint16(l) + uint16(cpu.X)) & 0xFF)
		cpu.read(dummyAddr)
	}
	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(l)
		cpu.tracer.addCPUByteCode(h)
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("$%04X,X @ %04X = %02X", a, addr, cpu.bus.peek(addr)))
	}
	return addr, pageCrossed
}

func (cpu *cpu) addressingAbsoluteY(_ *opcode, forceDummyRead bool) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	a := uint16(h)<<8 | uint16(l)
	addr = a + uint16(cpu.Y)
	pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
	if pageCrossed || forceDummyRead {
		dummyAddr := uint16(h)<<8 | ((uint16(l) + uint16(cpu.Y)) & 0xFF)
		cpu.read(dummyAddr)
	}
	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(l)
		cpu.tracer.addCPUByteCode(h)
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("$%04X,Y @ %04X = %02X", a, addr, cpu.bus.peek(addr)))
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
func (cpu *cpu) addressingAccumulator(_ *opcode) (addr uint16, pageCrossed bool) {
	cpu.read(cpu.PC) // dummy read
	if cpu.tracer != nil {
		cpu.tracer.setCPUAddressingResult("A")
	}
	return 0, false
}
func (cpu *cpu) addressingImplied(_ *opcode) (addr uint16, pageCrossed bool) {
	cpu.read(cpu.PC) // dummy read
	return 0, false
}

func (cpu *cpu) addressingImmediate(_ *opcode) (addr uint16, pageCrossed bool) {
	addr = cpu.PC
	cpu.PC++
	if cpu.tracer != nil {
		a := cpu.bus.peek(addr)
		cpu.tracer.addCPUByteCode(a)
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("#$%02X", a))
	}
	return addr, false
}

func (cpu *cpu) addressingIndexedIndirect(_ *opcode) (addr uint16, pageCrossed bool) {
	k := cpu.fetch()
	// https://www.nesdev.org/6502_cpu.txt
	// > pointer    R  read from the address, add X to it
	cpu.read(uint16(k)) // dummy read

	a := uint16(k + cpu.X)
	b := (a & 0xFF00) | uint16(byte(a)+1)
	addr = uint16(cpu.read(b))<<8 | uint16(cpu.read(a))

	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(k)
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", k, byte(a), addr, cpu.bus.peek(addr)))
	}
	return addr, false
}

func (cpu *cpu) addressingIndirect(_ *opcode) (addr uint16, pageCrossed bool) {
	l := cpu.fetch()
	h := cpu.fetch()
	a := uint16(h)<<8 | uint16(l)
	b := (a & 0xFF00) | uint16(byte(a)+1)
	addr = uint16(cpu.read(b))<<8 | uint16(cpu.read(a))

	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(l)
		cpu.tracer.addCPUByteCode(h)
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("($%04X) = %04X", a, addr))
	}
	return addr, false

}
func (cpu *cpu) addressingIndirectIndexed(_ *opcode, forceDummyRead bool) (addr uint16, pageCrossed bool) {
	a := uint16(cpu.fetch())
	b := (a & 0xFF00) | uint16(byte(a)+1)
	baseAddr := uint16(cpu.read(b))<<8 | uint16(cpu.read(a))
	addr = baseAddr + uint16(cpu.Y)
	pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
	if pageCrossed || forceDummyRead {
		h := baseAddr & 0xFF00
		l := baseAddr & 0x00FF
		dummyAddr := h | ((l + uint16(cpu.Y)) & 0xFF)
		cpu.read(dummyAddr)
	}

	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(byte(a))
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", byte(a), baseAddr, addr, cpu.bus.peek(addr)))
	}
	return addr, pageCrossed
}
func (cpu *cpu) addressingRelative(_ *opcode) (addr uint16, pageCrossed bool) {
	offset := uint16(cpu.fetch())
	if offset < 0x80 {
		addr = cpu.PC + offset
	} else {
		addr = cpu.PC + offset - 0x100
	}

	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(byte(offset))
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("$%04X", addr))
	}
	return addr, false
}

func (cpu *cpu) addressingZeroPage(_ *opcode) (addr uint16, pageCrossed bool) {
	a := cpu.fetch()
	addr = uint16(a)

	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(a)
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("$%02X = %02X", a, cpu.bus.peek(addr)))
	}
	return addr, false
}

func (cpu *cpu) addressingZeroPageX(_ *opcode) (addr uint16, pageCrossed bool) {
	a := cpu.fetch()
	// https://www.nesdev.org/6502_cpu.txt
	// > address   R  read from address, add index register to it
	cpu.read(uint16(a)) // dummy read

	addr = uint16(a+cpu.X) & 0xFF

	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(a)
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("$%02X,X @ %02X = %02X", a, addr, cpu.bus.peek(addr)))
	}
	return addr, false
}

func (cpu *cpu) addressingZeroPageY(_ *opcode) (addr uint16, pageCrossed bool) {
	a := cpu.fetch()
	// https://www.nesdev.org/6502_cpu.txt
	// > address   R  read from address, add index register to it
	cpu.read(uint16(a)) // dummy read

	addr = uint16(a+cpu.Y) & 0xFF

	if cpu.tracer != nil {
		cpu.tracer.addCPUByteCode(a)
		cpu.tracer.setCPUAddressingResult(fmt.Sprintf("$%02X,Y @ %02X = %02X", a, addr, cpu.bus.peek(addr)))
	}
	return addr, false
}
