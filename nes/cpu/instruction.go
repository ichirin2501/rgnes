package cpu

import (
	"github.com/ichirin2501/rgnes/nes/memory"
)

func (cpu *CPU) lda(addr uint16) {
	a := cpu.m.Read(addr)
	cpu.A = a
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) ldx(addr uint16) {
	a := cpu.m.Read(addr)
	cpu.X = a
	cpu.P.SetZN(cpu.X)
}

func (cpu *CPU) ldy(addr uint16) {
	a := cpu.m.Read(addr)
	cpu.Y = a
	cpu.P.SetZN(cpu.Y)
}

func (cpu *CPU) sta(addr uint16) {
	cpu.m.Write(addr, cpu.A)
}

func (cpu *CPU) stx(addr uint16) {
	cpu.m.Write(addr, cpu.X)
}

func (cpu *CPU) sty(addr uint16) {
	cpu.m.Write(addr, cpu.Y)
}

func (cpu *CPU) tax() {
	cpu.X = cpu.A
	cpu.P.SetZN(cpu.X)
}

func (cpu *CPU) tay() {
	cpu.Y = cpu.A
	cpu.P.SetZN(cpu.Y)
}

func (cpu *CPU) tsx() {
	cpu.X = cpu.S
	cpu.P.SetZN(cpu.X)
}

func (cpu *CPU) txa() {
	cpu.A = cpu.X
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) txs() {
	cpu.S = cpu.X
}

func (cpu *CPU) tya() {
	cpu.A = cpu.Y
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) adc(addr uint16) {
	a := cpu.A
	b := cpu.m.Read(addr)
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	v := a + b + c
	cpu.A = v
	cpu.P.SetZN(v)
	cpu.P.SetCarry(uint16(a)+uint16(b)+uint16(c) > 0xFF)
	cpu.P.SetOverflow((a^b)&0x80 == 0 && (a^v)&0x80 != 0)
}

func (cpu *CPU) and(addr uint16) {
	cpu.A &= cpu.m.Read(addr)
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) aslAcc() {
	cpu.P.SetCarry((cpu.A & 0x80) == 0x80)
	cpu.A <<= 1
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) asl(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.P.SetCarry((v & 0x80) == 0x80)
	cpu.m.Write(addr, v) // dummy write
	v <<= 1
	cpu.m.Write(addr, v)
	cpu.P.SetZN(v)
}

func (cpu *CPU) bit(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.P.SetOverflow((v & 0x40) == 0x40)
	cpu.P.SetNegative(v&0x80 != 0)
	cpu.P.SetZero((v & cpu.A) == 0x00)
}

func (cpu *CPU) cmp(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.compare(cpu.A, v)
}

func (cpu *CPU) cpx(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.compare(cpu.X, v)
}

func (cpu *CPU) cpy(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.compare(cpu.Y, v)
}

func (cpu *CPU) dec(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.m.Write(addr, v) // dummy write
	v--
	cpu.m.Write(addr, v)
	cpu.P.SetZN(v)
}

func (cpu *CPU) dex() {
	cpu.X--
	cpu.P.SetZN(cpu.X)
}

func (cpu *CPU) dey() {
	cpu.Y--
	cpu.P.SetZN(cpu.Y)
}

func (cpu *CPU) eor(addr uint16) {
	cpu.A ^= cpu.m.Read(addr)
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) inc(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.m.Write(addr, v) // dummy write
	v++
	cpu.m.Write(addr, v)
	cpu.P.SetZN(v)
}

func (cpu *CPU) inx() {
	cpu.X++
	cpu.P.SetZN(cpu.X)
}

func (cpu *CPU) iny() {
	cpu.Y++
	cpu.P.SetZN(cpu.Y)
}

func (cpu *CPU) lsrAcc() {
	cpu.P.SetCarry((cpu.A & 1) == 1)
	cpu.A >>= 1
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) lsr(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.m.Write(addr, v) // dummy write
	cpu.P.SetCarry((v & 1) == 1)
	v >>= 1
	cpu.m.Write(addr, v)
	cpu.P.SetZN(v)
}

func (cpu *CPU) ora(addr uint16) {
	cpu.A |= cpu.m.Read(addr)
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) rolAcc() {
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	cpu.P.SetCarry((cpu.A & 0x80) == 0x80)
	cpu.A = (cpu.A << 1) | c
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) rol(addr uint16) {
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	v := cpu.m.Read(addr)
	cpu.m.Write(addr, v) // dummy write
	cpu.P.SetCarry((v & 0x80) == 0x80)
	v = (v << 1) | c
	cpu.m.Write(addr, v)
	cpu.P.SetZN(v)
}

func (cpu *CPU) rorAcc() {
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	cpu.P.SetCarry((cpu.A & 1) == 1)
	cpu.A = (cpu.A >> 1) | (c << 7)
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) ror(addr uint16) {
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	v := cpu.m.Read(addr)
	cpu.m.Write(addr, v) // dummy write
	cpu.P.SetCarry((v & 1) == 1)
	v = (v >> 1) | (c << 7)
	cpu.m.Write(addr, v)
	cpu.P.SetZN(v)
}

func (cpu *CPU) sbc(addr uint16) {
	a := cpu.A
	b := cpu.m.Read(addr)
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	v := a - b - (1 - c)
	cpu.A = v
	cpu.P.SetCarry(int(a)-int(b)-int(1-c) >= 0)
	cpu.P.SetOverflow(((a^b)&0x80 != 0) && (a^v)&0x80 != 0)
	cpu.P.SetZN(v)
}

// https://www.nesdev.org/6502_cpu.txt
// PHA, PHP
// #  address R/W description
// --- ------- --- -----------------------------------------------
// 1    PC     R  fetch opcode, increment PC
// 2    PC     R  read next instruction byte (and throw it away)
// 3  $0100,S  W  push register on stack, decrement S
// ここはアドレッシングモード側でdummy readしてカバーされているのでdummy readはたぶん不要
func (cpu *CPU) pha() {
	//cpu.m.Read(cpu.PC) // dummy read
	cpu.push(cpu.A)
}
func (cpu *CPU) php() {
	//cpu.m.Read(cpu.PC) // dummy read
	cpu.push(cpu.P.Byte() | 0x30)
}

// https://www.nesdev.org/6502_cpu.txt
// PLA, PLP
// #  address R/W description
// --- ------- --- -----------------------------------------------
// 1    PC     R  fetch opcode, increment PC
// 2    PC     R  read next instruction byte (and throw it away)
// 3  $0100,S  R  increment S
// 4  $0100,S  R  pull register from stack
func (cpu *CPU) pla() {
	cpu.m.Read(cpu.PC) // dummy read
	cpu.A = cpu.pop()
	cpu.P.SetZN(cpu.A)
}
func (cpu *CPU) plp() {
	cpu.m.Read(cpu.PC) // dummy read
	cpu.P = StatusRegister((cpu.pop() & 0xEF) | (1 << 5))
}

func (cpu *CPU) jmp(addr uint16) {
	cpu.PC = addr
}

func (cpu *CPU) jsr(addr uint16) {
	cpu.push16(cpu.PC - 1)
	cpu.PC = addr
}

// https://www.nesdev.org/6502_cpu.txt
// RTS
// #  address R/W description
// --- ------- --- -----------------------------------------------
// 1    PC     R  fetch opcode, increment PC
// 2    PC     R  read next instruction byte (and throw it away)
// 3  $0100,S  R  increment S
// 4  $0100,S  R  pull PCL from stack, increment S
// 5  $0100,S  R  pull PCH from stack
// 6    PC     R  increment PC
func (cpu *CPU) rts() {
	cpu.m.Read(cpu.PC) // dummy read
	cpu.PC = cpu.pop16() + 1
}

// https://www.nesdev.org/6502_cpu.txt
// RTI
// #  address R/W description
// --- ------- --- -----------------------------------------------
// 1    PC     R  fetch opcode, increment PC
// 2    PC     R  read next instruction byte (and throw it away)
// 3  $0100,S  R  increment S
// 4  $0100,S  R  pull P from stack, increment S
// 5  $0100,S  R  pull PCL from stack, increment S
// 6  $0100,S  R  pull PCH from stack
func (cpu *CPU) rti() {
	cpu.m.Read(cpu.PC) // dummy read
	cpu.P = StatusRegister((cpu.pop() & 0xEF) | (1 << 5))
	cpu.PC = cpu.pop16()
}

func (cpu *CPU) bcc(addr uint16) int {
	if !cpu.P.IsCarry() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bcs(addr uint16) int {
	if cpu.P.IsCarry() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) beq(addr uint16) int {
	if cpu.P.IsZero() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bmi(addr uint16) int {
	if cpu.P.IsNegative() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bne(addr uint16) int {
	if !cpu.P.IsZero() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bpl(addr uint16) int {
	if !cpu.P.IsNegative() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bvc(addr uint16) int {
	if !cpu.P.IsOverflow() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bvs(addr uint16) int {
	if cpu.P.IsOverflow() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) clc() {
	cpu.P.SetCarry(false)
}

func (cpu *CPU) cld() {
	cpu.P.SetDecimal(false)
}

func (cpu *CPU) cli() {
	cpu.P.SetInterruptDisable(false)
}

func (cpu *CPU) clv() {
	cpu.P.SetOverflow(false)
}

func (cpu *CPU) sec() {
	cpu.P.SetCarry(true)
}

func (cpu *CPU) sed() {
	cpu.P.SetDecimal(true)
}

func (cpu *CPU) sei() {
	cpu.P.SetInterruptDisable(true)
}

// BRK
// #  address R/W description
// --- ------- --- -----------------------------------------------
// 1    PC     R  fetch opcode, increment PC
// 2    PC     R  read next instruction byte (and throw it away),
// 	       increment PC
// 3  $0100,S  W  push PCH on stack (with B flag set), decrement S
// 4  $0100,S  W  push PCL on stack, decrement S
// 5  $0100,S  W  push P on stack, decrement S
// 6   $FFFE   R  fetch PCL
// 7   $FFFF   R  fetch PCH
// ここはアドレッシングモード側でdummy readしてカバーされているのでdummy readはたぶん不要
func (cpu *CPU) brk() {
	// cpu.m.Read(cpu.PC) // dummy read
	cpu.push16(cpu.PC + 1)
	cpu.push(cpu.P.Byte() | 0x30)
	cpu.P.SetInterruptDisable(true)
	cpu.PC = memory.Read16(cpu.m, 0xFFFE)
}

// Two interrupts (/IRQ and /NMI) and two instructions (PHP and BRK) push the flags to the stack.
// In the byte pushed, bit 5 is always set to 1, and bit 4 is 1 if from an instruction (PHP or BRK) or 0 if from an interrupt line being pulled low (/IRQ or /NMI).
func (cpu *CPU) nmi() {
	cpu.push16(cpu.PC)
	cpu.push(cpu.P.Byte() | 0x20)
	cpu.P.SetInterruptDisable(true)
	cpu.PC = memory.Read16(cpu.m, 0xFFFA)
}

func (cpu *CPU) irq() {
	if cpu.P.IsInterruptDisable() {
		return
	}
	cpu.P.SetBreak1(false)
	cpu.push16(cpu.PC)
	cpu.push(cpu.P.Byte())
	cpu.P.SetInterruptDisable(true)
	cpu.PC = memory.Read16(cpu.m, 0xFFFE)
}

func (cpu *CPU) compare(a byte, b byte) {
	cpu.P.SetCarry(a >= b)
	cpu.P.SetNegative((a-b)&0x80 != 0)
	cpu.P.SetZero((a - b) == 0x00)
}

func (cpu *CPU) push(val byte) {
	cpu.m.Write(0x100|uint16(cpu.S), val)
	cpu.S--
}

func (cpu *CPU) pop() byte {
	cpu.S++
	return cpu.m.Read(0x100 | uint16(cpu.S))
}

func (cpu *CPU) push16(val uint16) {
	l := byte(val & 0xFF)
	h := byte(val >> 8)
	cpu.push(h)
	cpu.push(l)
}

func (cpu *CPU) pop16() uint16 {
	l := cpu.pop()
	h := cpu.pop()
	return uint16(h)<<8 | uint16(l)
}

func (cpu *CPU) branch(addr uint16) int {
	cycle := 1
	if pagesCross(cpu.PC, addr) {
		cycle++
	}
	cpu.PC = addr
	return cycle
}

func pagesCross(a uint16, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

// undocumented opcode

func (cpu *CPU) kil() {}

func (cpu *CPU) slo(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.m.Write(addr, v) // dummy write
	cpu.P.SetCarry((v & 0x80) == 0x80)
	v <<= 1
	cpu.m.Write(addr, v)

	cpu.A |= v
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) anc(addr uint16) {
	a := cpu.m.Read(addr)
	cpu.A &= a
	cpu.P.SetZN(cpu.A)
	cpu.P.SetCarry(cpu.P.IsNegative())
}

func (cpu *CPU) rla(addr uint16) {
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	v := cpu.m.Read(addr)
	cpu.m.Write(addr, v) // dummy write
	cpu.P.SetCarry((v & 0x80) == 0x80)
	v = (v << 1) | c
	cpu.m.Write(addr, v)

	cpu.A &= v
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) sre(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.m.Write(addr, v) // dummy write
	cpu.P.SetCarry((v & 1) == 1)
	v >>= 1
	cpu.m.Write(addr, v)

	cpu.A ^= v
	cpu.P.SetZN(cpu.A)
}

func (cpu *CPU) alr(addr uint16) {
	// A =(A&#{imm})/2
	// N, Z, C
	v := cpu.A & cpu.m.Read(addr)
	cpu.P.SetCarry((v & 1) == 1)
	v >>= 1
	cpu.P.SetZN(v)
	cpu.A = v
}

func (cpu *CPU) rra(addr uint16) {
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	k := cpu.m.Read(addr)
	cpu.m.Write(addr, k) // dummy write
	cpu.P.SetCarry((k & 1) == 1)
	k = (k >> 1) | (c << 7)
	cpu.m.Write(addr, k)

	a := cpu.A
	b := k
	c = byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	v := a + b + c
	cpu.A = v
	cpu.P.SetZN(v)
	cpu.P.SetCarry(uint16(a)+uint16(b)+uint16(c) > 0xFF)
	cpu.P.SetOverflow((a^b)&0x80 == 0 && (a^v)&0x80 != 0)
}

func (cpu *CPU) arr(addr uint16) {
	// A:=(A&#{imm})/2
	// N V Z C
	// N and Z are normal, but C is bit 6 and V is bit 6 xor bit 5.
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 0x80
	}
	v := ((cpu.A & cpu.m.Read(addr)) >> 1) | c
	cpu.P.SetZN(v)
	cpu.P.SetCarry((v & 0x40) == 0x40)
	cpu.P.SetOverflow(((v & 0x40) ^ ((v & 0x20) << 1)) == 0x40)
	cpu.A = v
}

func (cpu *CPU) sax(addr uint16) {
	cpu.m.Write(addr, cpu.A&cpu.X)
}

func (cpu *CPU) xaa() {}

func (cpu *CPU) ahx() {}

func (cpu *CPU) tas() {}

func (cpu *CPU) shy(addr uint16) {
	if pagesCross(addr, addr-uint16(cpu.X)) {
		addr &= uint16(cpu.Y) << 8
	}
	res := cpu.Y & (byte(addr>>8) + 1)
	cpu.m.Write(addr, res)
}

func (cpu *CPU) shx(addr uint16) {
	if pagesCross(addr, addr-uint16(cpu.Y)) {
		addr &= uint16(cpu.X) << 8
	}
	res := cpu.X & (byte(addr>>8) + 1)
	cpu.m.Write(addr, res)
}

func (cpu *CPU) lax(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.X = v
	cpu.A = v
	cpu.P.SetZN(v)
}

func (cpu *CPU) las() {}

func (cpu *CPU) dcp(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.m.Write(addr, v) // dummy write
	v--
	cpu.compare(cpu.A, v)
	cpu.m.Write(addr, v)
}

func (cpu *CPU) axs(addr uint16) {
	// X:=A&X-#{imm}
	// N, Z, C
	t := cpu.A & cpu.X
	v := cpu.m.Read(addr)
	cpu.X = t - v
	cpu.P.SetCarry(int(t)-int(v) >= 0)
	cpu.P.SetZN(cpu.X)
}

func (cpu *CPU) isb(addr uint16) {
	k := cpu.m.Read(addr)
	cpu.m.Write(addr, k) // dummy write
	k++
	cpu.m.Write(addr, k)

	a := cpu.A
	b := k
	c := byte(0)
	if cpu.P.IsCarry() {
		c = 1
	}
	v := a - b - (1 - c)
	cpu.A = v
	cpu.P.SetCarry(int(a)-int(b)-int(1-c) >= 0)
	cpu.P.SetOverflow(((a^b)&0x80 != 0) && (a^v)&0x80 != 0)
	cpu.P.SetZN(v)
}
