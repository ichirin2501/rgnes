package cpu

import (
	"fmt"
	"io"
	"os"

	"github.com/ichirin2501/rgnes/nes/bus"
	"github.com/ichirin2501/rgnes/nes/memory"
)

type CPUCycle struct {
	stall  int
	cycles int
}

type InterruptType byte

const (
	InterruptNone InterruptType = iota
	InterruptNMI
	InterruptIRQ
)

type Interrupter struct {
	I InterruptType
}

func (i *Interrupter) TriggerNMI() {
	i.I = InterruptNMI
}
func (i *Interrupter) TriggerIRQ() {
	i.I = InterruptIRQ
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

var debugWriter io.Writer = os.Stdout

type Trace struct {
	A                   byte
	X                   byte
	Y                   byte
	PC                  uint16
	S                   byte
	P                   byte
	ByteCode            []byte
	Opcode              opcode
	AddressingResult    string
	InstructionReadByte *byte
	Cycle               int
	PPUX                uint16
	PPUY                uint16
	PPUVBlankState      bool
}

func (t *Trace) NESTestString() string {
	bc := ""
	switch len(t.ByteCode) {
	case 1:
		bc = fmt.Sprintf("%02X      ", t.ByteCode[0])
	case 2:
		bc = fmt.Sprintf("%02X %02X   ", t.ByteCode[0], t.ByteCode[1])
	case 3:
		bc = fmt.Sprintf("%02X %02X %02X", t.ByteCode[0], t.ByteCode[1], t.ByteCode[2])
	}
	ar := ""
	if t.InstructionReadByte == nil {
		ar = t.AddressingResult
	} else {
		ar = fmt.Sprintf("%s = %02X", t.AddressingResult, *t.InstructionReadByte)
	}
	op := t.Opcode.Name.String()
	if t.Opcode.Unofficial {
		op = "*" + t.Opcode.Name.String()
	}

	// C000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:24 SP:FD PPU:  0, 45 CYC:15
	return fmt.Sprintf("%04X  %s %4s %-27s A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d",
		t.PC,
		bc,
		op,
		ar,
		t.A,
		t.X,
		t.Y,
		t.P,
		t.S,
		t.PPUY,
		t.PPUX,
		t.Cycle,
	)
	// return fmt.Sprintf("%04X  %s %4s %-27s A:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%3d",
	// 	t.PC,
	// 	bc,
	// 	t.Opcode.Name, // op,
	// 	"",            //ar,
	// 	t.A,
	// 	t.X,
	// 	t.Y,
	// 	t.P,
	// 	t.S,
	// 	(t.Cycle-21)%341,
	// )
}

func (t *Trace) SetCPURegisterA(v byte)            { t.A = v }
func (t *Trace) SetCPURegisterX(v byte)            { t.X = v }
func (t *Trace) SetCPURegisterY(v byte)            { t.Y = v }
func (t *Trace) SetCPURegisterPC(v uint16)         { t.PC = v }
func (t *Trace) SetCPURegisterS(v byte)            { t.S = v }
func (t *Trace) SetCPURegisterP(v byte)            { t.P = v }
func (t *Trace) SetCPUOpcode(v opcode)             { t.Opcode = v }
func (t *Trace) SetCPUAddressingResult(v string)   { t.AddressingResult = v }
func (t *Trace) SetCPUInstructionReadByte(v *byte) { t.InstructionReadByte = v }
func (t *Trace) SetPPUX(v uint16)                  { t.PPUX = v }
func (t *Trace) SetPPUY(v uint16)                  { t.PPUY = v }
func (t *Trace) SetPPUVBlankState(v bool)          { t.PPUVBlankState = v }
func (t *Trace) AddCPUCycle(v int)                 { t.Cycle += v }
func (t *Trace) AddCPUByteCode(v byte) {
	t.ByteCode = append(t.ByteCode, v)
}
func (t *Trace) Reset() {
	t.AddressingResult = ""
	t.InstructionReadByte = nil
	t.ByteCode = t.ByteCode[:0]
}

type CPU struct {
	A  byte   // Accumulator
	X  byte   // Index
	Y  byte   // Index
	PC uint16 // Program Counter
	S  byte   // Stack Pointer
	P  StatusRegister

	m         *bus.CPUBus
	cycle     *CPUCycle
	interrupt *Interrupter
	t         *Trace
}

func NewCPU(mem *bus.CPUBus, cycle *CPUCycle, i *Interrupter, t *Trace) *CPU {
	return &CPU{
		A:  0x00,
		X:  0x00,
		Y:  0x00,
		PC: 0x8000,
		S:  0xFD,
		P:  StatusRegister(0x24),

		m:         mem,
		cycle:     cycle,
		interrupt: i,
		t:         t,
	}
}

// TODO: after reset
func (cpu *CPU) Reset() {
	cpu.PC = memory.Read16(cpu.m, 0xFFFC)
	cpu.P = StatusRegister(0x24)
}

func (cpu *CPU) Step() int {
	if cpu.cycle.Stall() > 0 {
		cpu.cycle.AddStall(-1)
	}

	additionalCycle := 0

	switch cpu.interrupt.I {
	case InterruptNMI:
		cpu.nmi()
		// https://www.nesdev.org/wiki/Consistent_frame_synchronization#Ideal_NMI
		additionalCycle += 7
	case InterruptIRQ:
		cpu.irq()
		additionalCycle += 7
	}
	cpu.interrupt.I = InterruptNone

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

	cpu.cycle.AddCycles(opcode.Cycle + additionalCycle)
	return opcode.Cycle + additionalCycle
}

func (cpu *CPU) fetch() byte {
	v := cpu.m.Read(cpu.PC)
	cpu.PC++
	return v
}

func (cpu *CPU) fetchOperand(op *opcode) (addr uint16, pageCrossed bool) {
	pageCrossed = false

	switch op.Mode {
	case absolute:
		l := cpu.fetch()
		h := cpu.fetch()
		addr = uint16(h)<<8 | uint16(l)
		if cpu.t != nil {
			cpu.t.AddCPUByteCode(l)
			cpu.t.AddCPUByteCode(h)
			if op.Name == JMP || op.Name == JSR {
				cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X", addr))
			} else {
				cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X = %02X", addr, cpu.m.ReadForTest(addr)))
			}
		}
	case absoluteX:
		l := cpu.fetch()
		h := cpu.fetch()
		a := uint16(h)<<8 | uint16(l)
		addr = a + uint16(cpu.X)
		pageCrossed = pagesCross(addr, addr-uint16(cpu.X))
		if pageCrossed {
			// dummy read
			dummyAddr := uint16(h)<<8 | ((uint16(l) + uint16(cpu.X)) & 0xFF)
			cpu.m.Read(dummyAddr)
		}
		if cpu.t != nil {
			cpu.t.AddCPUByteCode(l)
			cpu.t.AddCPUByteCode(h)
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X,X @ %04X = %02X", a, addr, cpu.m.ReadForTest(addr)))
		}
	case absoluteY:
		l := cpu.fetch()
		h := cpu.fetch()
		a := uint16(h)<<8 | uint16(l)
		addr = a + uint16(cpu.Y)
		pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
		if pageCrossed {
			// dummy read
			dummyAddr := uint16(h)<<8 | ((uint16(l) + uint16(cpu.Y)) & 0xFF)
			cpu.m.Read(dummyAddr)
		}
		if cpu.t != nil {
			cpu.t.AddCPUByteCode(l)
			cpu.t.AddCPUByteCode(h)
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%04X,Y @ %04X = %02X", a, addr, cpu.m.ReadForTest(addr)))
		}
	case accumulator:
		addr = 0
		if cpu.t != nil {
			cpu.t.SetCPUAddressingResult("A")
		}
	case immediate:
		addr = cpu.PC
		cpu.PC++
		if cpu.t != nil {
			a := cpu.m.ReadForTest(addr)
			cpu.t.AddCPUByteCode(a)
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("#$%02X", a))
		}
	case implied:
		addr = 0
	case indexedIndirect:
		k := cpu.fetch()
		a := uint16(k + cpu.X)
		b := (a & 0xFF00) | uint16(byte(a)+1)
		addr = uint16(cpu.m.Read(b))<<8 | uint16(cpu.m.Read(a))
		if cpu.t != nil {
			cpu.t.AddCPUByteCode(k)
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", k, byte(a), addr, cpu.m.ReadForTest(addr)))
		}
	case indirect:
		l := cpu.fetch()
		h := cpu.fetch()
		a := uint16(h)<<8 | uint16(l)
		b := (a & 0xFF00) | uint16(byte(a)+1)
		addr = uint16(cpu.m.Read(b))<<8 | uint16(cpu.m.Read(a))
		if cpu.t != nil {
			cpu.t.AddCPUByteCode(l)
			cpu.t.AddCPUByteCode(h)
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("($%04X) = %04X", a, addr))
		}
	case indirectIndexed:
		a := uint16(cpu.fetch())
		b := (a & 0xFF00) | uint16(byte(a)+1)
		baseAddr := uint16(cpu.m.Read(b))<<8 | uint16(cpu.m.Read(a))
		addr = baseAddr + uint16(cpu.Y)
		pageCrossed = pagesCross(addr, addr-uint16(cpu.Y))
		if pageCrossed {
			// dummy read
			h := baseAddr & 0xFF00
			l := baseAddr & 0x00FF
			dummyAddr := h | ((l + uint16(cpu.Y)) & 0xFF)
			cpu.m.Read(dummyAddr)
		}
		if cpu.t != nil {
			cpu.t.AddCPUByteCode(byte(a))
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", byte(a), baseAddr, addr, cpu.m.ReadForTest(addr)))
		}
	case relative:
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
	case zeroPage:
		a := cpu.fetch()
		addr = uint16(a)
		if cpu.t != nil {
			cpu.t.AddCPUByteCode(a)
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%02X = %02X", a, cpu.m.ReadForTest(addr)))
		}
	case zeroPageX:
		a := cpu.fetch()
		addr = uint16(a+cpu.X) & 0xFF
		if cpu.t != nil {
			cpu.t.AddCPUByteCode(a)
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%02X,X @ %02X = %02X", a, addr, cpu.m.ReadForTest(addr)))
		}
	case zeroPageY:
		a := cpu.fetch()
		addr = uint16(a+cpu.Y) & 0xFF
		if cpu.t != nil {
			cpu.t.AddCPUByteCode(a)
			cpu.t.SetCPUAddressingResult(fmt.Sprintf("$%02X,Y @ %02X = %02X", a, addr, cpu.m.ReadForTest(addr)))
		}
	}

	return addr, pageCrossed
}
