package cpu

import (
	"github.com/ichirin2501/rgnes/nes/memory"
)

func (cpu *CPU) lda(addr uint16) {
	a := cpu.m.Read(addr)
	cpu.r.A = a
	cpu.r.UpdateZeroFlag(cpu.r.A)
	cpu.r.UpdateNegativeFlag(cpu.r.A)
}

func (cpu *CPU) ldx(addr uint16) {
	a := cpu.m.Read(addr)
	cpu.r.X = a
	cpu.r.UpdateZeroFlag(cpu.r.X)
	cpu.r.UpdateNegativeFlag(cpu.r.X)
}

func (cpu *CPU) ldy(addr uint16) {
	a := cpu.m.Read(addr)
	cpu.r.Y = a
	cpu.r.UpdateZeroFlag(cpu.r.Y)
	cpu.r.UpdateNegativeFlag(cpu.r.Y)
}

func (cpu *CPU) sta(addr uint16) {
	cpu.m.Write(addr, cpu.r.A)
}

func (cpu *CPU) stx(addr uint16) {
	cpu.m.Write(addr, cpu.r.X)
}

func (cpu *CPU) sty(addr uint16) {
	cpu.m.Write(addr, cpu.r.Y)
}

func (cpu *CPU) tax() {
	cpu.r.X = cpu.r.A
	cpu.r.UpdateNegativeFlag(cpu.r.X)
	cpu.r.UpdateZeroFlag(cpu.r.X)
}

func (cpu *CPU) tay() {
	cpu.r.Y = cpu.r.A
	cpu.r.UpdateNegativeFlag(cpu.r.Y)
	cpu.r.UpdateZeroFlag(cpu.r.Y)
}

func (cpu *CPU) tsx() {
	cpu.r.X = cpu.r.S
	cpu.r.UpdateNegativeFlag(cpu.r.X)
	cpu.r.UpdateZeroFlag(cpu.r.X)
}

func (cpu *CPU) txa() {
	cpu.r.A = cpu.r.X
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) txs() {
	cpu.r.S = cpu.r.X
}

func (cpu *CPU) tya() {
	cpu.r.A = cpu.r.Y
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) adc(addr uint16) {
	a := cpu.r.A
	b := cpu.m.Read(addr)
	c := byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	v := a + b + c
	cpu.r.A = v
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
	cpu.r.SetCarryFlag(uint16(a)+uint16(b)+uint16(c) > 0xFF)
	cpu.r.SetOverflowFlag((a^b)&0x80 == 0 && (a^v)&0x80 != 0)
}

func (cpu *CPU) and(addr uint16) {
	cpu.r.A &= cpu.m.Read(addr)
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) aslAcc() {
	cpu.r.SetCarryFlag((cpu.r.A & 0x80) == 0x80)
	cpu.r.A <<= 1
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) asl(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.r.SetCarryFlag((v & 0x80) == 0x80)
	v <<= 1
	cpu.m.Write(addr, v)
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
}

func (cpu *CPU) bit(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.r.SetOverflowFlag((v & 0x40) == 0x40)
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v & cpu.r.A)
}

func (cpu *CPU) cmp(addr uint16) {
	v := cpu.m.Read(addr)
	compare(cpu.r, cpu.r.A, v)
}

func (cpu *CPU) cpx(addr uint16) {
	v := cpu.m.Read(addr)
	compare(cpu.r, cpu.r.X, v)
}

func (cpu *CPU) cpy(addr uint16) {
	v := cpu.m.Read(addr)
	compare(cpu.r, cpu.r.Y, v)
}

func (cpu *CPU) dec(addr uint16) {
	v := cpu.m.Read(addr) - 1
	cpu.m.Write(addr, v)
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
}

func (cpu *CPU) dex() {
	cpu.r.X--
	cpu.r.UpdateNegativeFlag(cpu.r.X)
	cpu.r.UpdateZeroFlag(cpu.r.X)
}

func (cpu *CPU) dey() {
	cpu.r.Y--
	cpu.r.UpdateNegativeFlag(cpu.r.Y)
	cpu.r.UpdateZeroFlag(cpu.r.Y)
}

func (cpu *CPU) eor(addr uint16) {
	cpu.r.A ^= cpu.m.Read(addr)
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) inc(addr uint16) {
	v := cpu.m.Read(addr) + 1
	cpu.m.Write(addr, v)
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
}

func (cpu *CPU) inx() {
	cpu.r.X++
	cpu.r.UpdateNegativeFlag(cpu.r.X)
	cpu.r.UpdateZeroFlag(cpu.r.X)
}

func (cpu *CPU) iny() {
	cpu.r.Y++
	cpu.r.UpdateNegativeFlag(cpu.r.Y)
	cpu.r.UpdateZeroFlag(cpu.r.Y)
}

func (cpu *CPU) lsrAcc() {
	cpu.r.SetCarryFlag((cpu.r.A & 1) == 1)
	cpu.r.A >>= 1
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) lsr(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.r.SetCarryFlag((v & 1) == 1)
	v >>= 1
	cpu.m.Write(addr, v)
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
}

func (cpu *CPU) ora(addr uint16) {
	cpu.r.A |= cpu.m.Read(addr)
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) rolAcc() {
	c := byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	cpu.r.SetCarryFlag((cpu.r.A & 0x80) == 0x80)
	cpu.r.A = (cpu.r.A << 1) | c
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) rol(addr uint16) {
	c := byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	v := cpu.m.Read(addr)
	cpu.r.SetCarryFlag((v & 0x80) == 0x80)
	v = (v << 1) | c
	cpu.m.Write(addr, v)
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
}

func (cpu *CPU) rorAcc() {
	c := byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	cpu.r.SetCarryFlag((cpu.r.A & 1) == 1)
	cpu.r.A = (cpu.r.A >> 1) | (c << 7)
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) ror(addr uint16) {
	c := byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	v := cpu.m.Read(addr)
	cpu.r.SetCarryFlag((v & 1) == 1)
	v = (v >> 1) | (c << 7)
	cpu.m.Write(addr, v)
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
}

func (cpu *CPU) sbc(addr uint16) {
	a := cpu.r.A
	b := cpu.m.Read(addr)
	c := byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	v := a - b - (1 - c)
	cpu.r.A = v
	cpu.r.SetCarryFlag(int(a)-int(b)-int(1-c) >= 0)
	cpu.r.SetOverflowFlag(((a^b)&0x80 != 0) && (a^v)&0x80 != 0)
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
}

func (cpu *CPU) pha() {
	cpu.push(cpu.r.A)
}

func (cpu *CPU) php() {
	// TODO: 数値やめる
	cpu.push(cpu.r.P | 0x30)
}

func (cpu *CPU) pla() {
	cpu.r.A = cpu.pop()
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) plp() {
	cpu.r.P = (cpu.pop() & 0xEF) | reservedFlagMask
}

func (cpu *CPU) jmp(addr uint16) {
	cpu.r.PC = addr
}

func (cpu *CPU) jsr(addr uint16) {
	cpu.push16(cpu.r.PC - 1)
	cpu.r.PC = addr
}

func (cpu *CPU) rts() {
	cpu.r.PC = cpu.pop16() + 1
}

func (cpu *CPU) rti() {
	cpu.r.P = (cpu.pop() & 0xEF) | reservedFlagMask
	cpu.r.PC = cpu.pop16()
}

func (cpu *CPU) bcc(addr uint16) int {
	if !cpu.r.CarryFlag() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bcs(addr uint16) int {
	if cpu.r.CarryFlag() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) beq(addr uint16) int {
	if cpu.r.ZeroFlag() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bmi(addr uint16) int {
	if cpu.r.NegativeFlag() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bne(addr uint16) int {
	if !cpu.r.ZeroFlag() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bpl(addr uint16) int {
	if !cpu.r.NegativeFlag() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bvc(addr uint16) int {
	if !cpu.r.OverflowFlag() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) bvs(addr uint16) int {
	if cpu.r.OverflowFlag() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *CPU) clc() {
	cpu.r.SetCarryFlag(false)
}

func (cpu *CPU) cld() {
	cpu.r.SetDecimalFlag(false)
}

func (cpu *CPU) cli() {
	cpu.r.SetInterruptDisableFlag(false)
}

func (cpu *CPU) clv() {
	cpu.r.SetOverflowFlag(false)
}

func (cpu *CPU) sec() {
	cpu.r.SetCarryFlag(true)
}

func (cpu *CPU) sed() {
	cpu.r.SetDecimalFlag(true)
}

func (cpu *CPU) sei() {
	cpu.r.SetInterruptDisableFlag(true)
}

func (cpu *CPU) brk() {
	cpu.push16(cpu.r.PC + 1)
	cpu.push(cpu.r.P | 0x30)
	cpu.r.SetInterruptDisableFlag(true)
	cpu.r.PC = memory.Read16(cpu.m, 0xFFFE)
}

// Two interrupts (/IRQ and /NMI) and two instructions (PHP and BRK) push the flags to the stack.
// In the byte pushed, bit 5 is always set to 1, and bit 4 is 1 if from an instruction (PHP or BRK) or 0 if from an interrupt line being pulled low (/IRQ or /NMI).
func (cpu *CPU) nmi() {
	cpu.push16(cpu.r.PC)
	cpu.push(cpu.r.P | 0x20)
	cpu.r.SetInterruptDisableFlag(true)
	cpu.r.PC = memory.Read16(cpu.m, 0xFFFA)
}

func (cpu *CPU) irq() {
	if cpu.r.InterruptDisableFlag() {
		return
	}
	cpu.r.SetBreakFlag(false)
	cpu.push16(cpu.r.PC)
	cpu.push(cpu.r.P)
	cpu.r.SetInterruptDisableFlag(true)
	cpu.r.PC = memory.Read16(cpu.m, 0xFFFE)
}

func compare(r *cpuRegister, a byte, b byte) {
	r.SetCarryFlag(a >= b)
	r.UpdateNegativeFlag(a - b)
	r.UpdateZeroFlag(a - b)
}

func (cpu *CPU) push(val byte) {
	cpu.m.Write(0x100|uint16(cpu.r.S), val)
	cpu.r.S--
}

func (cpu *CPU) pop() byte {
	cpu.r.S++
	return cpu.m.Read(0x100 | uint16(cpu.r.S))
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
	if pagesCross(cpu.r.PC, addr) {
		cycle++
	}
	cpu.r.PC = addr
	return cycle
}

func pagesCross(a uint16, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

// undocumented opcode

func (cpu *CPU) kil() {}

func (cpu *CPU) slo(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.r.SetCarryFlag((v & 0x80) == 0x80)
	v <<= 1
	cpu.m.Write(addr, v)

	cpu.r.A |= v
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) anc(addr uint16) {
	a := cpu.m.Read(addr)
	cpu.r.A &= a
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
	cpu.r.SetCarryFlag(cpu.r.NegativeFlag())
}

func (cpu *CPU) rla(addr uint16) {
	c := byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	v := cpu.m.Read(addr)
	cpu.r.SetCarryFlag((v & 0x80) == 0x80)
	v = (v << 1) | c
	cpu.m.Write(addr, v)

	cpu.r.A &= v
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) sre(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.r.SetCarryFlag((v & 1) == 1)
	v >>= 1
	cpu.m.Write(addr, v)

	cpu.r.A ^= v
	cpu.r.UpdateNegativeFlag(cpu.r.A)
	cpu.r.UpdateZeroFlag(cpu.r.A)
}

func (cpu *CPU) alr() {}

func (cpu *CPU) rra(addr uint16) {
	c := byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	k := cpu.m.Read(addr)
	cpu.r.SetCarryFlag((k & 1) == 1)
	k = (k >> 1) | (c << 7)
	cpu.m.Write(addr, k)

	a := cpu.r.A
	b := k
	c = byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	v := a + b + c
	cpu.r.A = v
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
	cpu.r.SetCarryFlag(uint16(a)+uint16(b)+uint16(c) > 0xFF)
	cpu.r.SetOverflowFlag((a^b)&0x80 == 0 && (a^v)&0x80 != 0)
}

// func (cpu *CPU) arr(addr uint16) {
// 	// TODO
// 	a := cpu.m.Read(addr)
// 	cpu.r.A = a
// 	cpu.r.UpdateNegativeFlag(cpu.r.A)
// 	cpu.r.UpdateZeroFlag(cpu.r.A)
// 	cpu.rorAcc()
// }

func (cpu *CPU) sax(addr uint16) {
	cpu.m.Write(addr, cpu.r.A&cpu.r.X)
}

func (cpu *CPU) xaa() {}

func (cpu *CPU) ahx() {}

func (cpu *CPU) tas() {}

func (cpu *CPU) shy() {}

func (cpu *CPU) shx() {}

func (cpu *CPU) lax(addr uint16) {
	v := cpu.m.Read(addr)
	cpu.r.X = v
	cpu.r.A = v
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
}

func (cpu *CPU) las() {}

func (cpu *CPU) dcp(addr uint16) {
	v := cpu.m.Read(addr) - 1
	compare(cpu.r, cpu.r.A, v)
	cpu.m.Write(addr, v)
}

func (cpu *CPU) axs() {}

func (cpu *CPU) isb(addr uint16) {
	k := cpu.m.Read(addr) + 1
	cpu.m.Write(addr, k)

	a := cpu.r.A
	b := k
	c := byte(0)
	if cpu.r.CarryFlag() {
		c = 1
	}
	v := a - b - (1 - c)
	cpu.r.A = v
	cpu.r.SetCarryFlag(int(a)-int(b)-int(1-c) >= 0)
	cpu.r.SetOverflowFlag(((a^b)&0x80 != 0) && (a^v)&0x80 != 0)
	cpu.r.UpdateNegativeFlag(v)
	cpu.r.UpdateZeroFlag(v)
}
