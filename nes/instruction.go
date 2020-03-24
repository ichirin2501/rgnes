package nes

func lda(r *cpuRegister, m MemoryReader, addr uint16) {
	r.A = m.Read(addr)
	r.UpdateZeroFlag(r.A)
	r.UpdateNegativeFlag(r.A)
}

func ldx(r *cpuRegister, m MemoryReader, addr uint16) {
	r.X = m.Read(addr)
	r.UpdateZeroFlag(r.X)
	r.UpdateNegativeFlag(r.X)
}

func ldy(r *cpuRegister, m MemoryReader, addr uint16) {
	r.Y = m.Read(addr)
	r.UpdateZeroFlag(r.Y)
	r.UpdateNegativeFlag(r.Y)
}

func sta(r *cpuRegister, m MemoryWriter, addr uint16) {
	m.Write(addr, r.A)
}

func stx(r *cpuRegister, m MemoryWriter, addr uint16) {
	m.Write(addr, r.X)
}

func sty(r *cpuRegister, m MemoryWriter, addr uint16) {
	m.Write(addr, r.Y)
}

func tax(r *cpuRegister) {
	r.X = r.A
	r.UpdateNegativeFlag(r.X)
	r.UpdateZeroFlag(r.X)
}

func tay(r *cpuRegister) {
	r.Y = r.A
	r.UpdateNegativeFlag(r.Y)
	r.UpdateZeroFlag(r.Y)
}

func tsx(r *cpuRegister) {
	r.X = r.S
	r.UpdateNegativeFlag(r.X)
	r.UpdateZeroFlag(r.X)
}

func txa(r *cpuRegister) {
	r.A = r.X
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func txs(r *cpuRegister) {
	r.S = r.X
}

func tya(r *cpuRegister) {
	r.A = r.Y
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func adc(r *cpuRegister, m MemoryReader, addr uint16) {
	a := r.A
	b := m.Read(addr)
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	v := a + b + c
	r.A = v
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	r.SetCarryFlag(uint16(a)+uint16(b)+uint16(c) > 0xFF)
	r.SetOverflowFlag((a^b)&0x80 == 0 && (a^v)&0x80 != 0)
}

func and(r *cpuRegister, m MemoryReader, addr uint16) {
	r.A = r.A & m.Read(addr)
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func aslAcc(r *cpuRegister) {
	r.SetCarryFlag((r.A & 0x80) == 0x80)
	r.A <<= 1
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func asl(r *cpuRegister, m Memory, addr uint16) {
	v := m.Read(addr)
	r.SetCarryFlag((v & 0x80) == 0x80)
	v <<= 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func bit(r *cpuRegister, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	r.SetOverflowFlag((v & 0x40) == 0x40)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v & r.A)
}

func cmp(r *cpuRegister, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	compare(r, r.A, v)
}

func cpx(r *cpuRegister, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	compare(r, r.X, v)
}

func cpy(r *cpuRegister, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	compare(r, r.Y, v)
}

func dec(r *cpuRegister, m Memory, addr uint16) {
	v := m.Read(addr) - 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func dex(r *cpuRegister) {
	r.X--
	r.UpdateNegativeFlag(r.X)
	r.UpdateZeroFlag(r.X)
}

func dey(r *cpuRegister) {
	r.Y--
	r.UpdateNegativeFlag(r.Y)
	r.UpdateZeroFlag(r.Y)
}

func eor(r *cpuRegister, m MemoryReader, addr uint16) {
	r.A ^= m.Read(addr)
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func inc(r *cpuRegister, m Memory, addr uint16) {
	v := m.Read(addr) + 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func inx(r *cpuRegister) {
	r.X++
	r.UpdateNegativeFlag(r.X)
	r.UpdateZeroFlag(r.X)
}

func iny(r *cpuRegister) {
	r.Y++
	r.UpdateNegativeFlag(r.Y)
	r.UpdateZeroFlag(r.Y)
}

func lsrAcc(r *cpuRegister) {
	r.SetCarryFlag((r.A & 1) == 1)
	r.A >>= 1
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func lsr(r *cpuRegister, m Memory, addr uint16) {
	v := m.Read(addr)
	r.SetCarryFlag((v & 1) == 1)
	v >>= 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func ora(r *cpuRegister, m MemoryReader, addr uint16) {
	r.A |= m.Read(addr)
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func rolAcc(r *cpuRegister) {
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	r.SetCarryFlag((r.A & 0x80) == 0x80)
	r.A = (r.A << 1) | c
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func rol(r *cpuRegister, m Memory, addr uint16) {
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	v := m.Read(addr)
	r.SetCarryFlag((v & 0x80) == 0x80)
	v = (v << 1) | c
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func rorAcc(r *cpuRegister) {
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	r.SetCarryFlag((r.A & 1) == 1)
	r.A = (r.A >> 1) | (c << 7)
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func ror(r *cpuRegister, m Memory, addr uint16) {
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	v := m.Read(addr)
	r.SetCarryFlag((v & 1) == 1)
	v = (v >> 1) | (c << 7)
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func sbc(r *cpuRegister, m MemoryReader, addr uint16) {
	a := r.A
	b := m.Read(addr)
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	v := a - b - (1 - c)
	r.A = v
	r.SetCarryFlag(uint16(a)-uint16(b)-uint16(1-c) >= 0)
	r.SetOverflowFlag(((a^b)&0x80 != 0) && (a^v)&0x80 != 0)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func pha(r *cpuRegister, m MemoryWriter) {
	push(r, m, r.A)
}

func php(r *cpuRegister, m MemoryWriter) {
	push(r, m, r.P|breakFlagMask)
}

func pla(r *cpuRegister, m MemoryReader) {
	r.A = pop(r, m)
	r.UpdateNegativeFlag(r.A)
	r.UpdateZeroFlag(r.A)
}

func plp(r *cpuRegister, m MemoryReader) {
	r.P = (pop(r, m) & 0xEF) | reservedFlagMask
}

func jmp(r *cpuRegister, addr uint16) {
	r.PC = addr
}

func jsr(r *cpuRegister, m MemoryWriter, addr uint16) {
	push16(r, m, r.PC-1)
	r.PC = addr
}

func rts(r *cpuRegister, m MemoryReader) {
	r.PC = pop16(r, m) + 1
}

func rti(r *cpuRegister, m MemoryReader) {
	r.P = (pop(r, m) & 0xEF) | reservedFlagMask
	r.PC = pop16(r, m)
}

func bcc(r *cpuRegister, addr uint16) int {
	if !r.CarryFlag() {
		return branch(r, addr)
	}
	return 0
}

func bcs(r *cpuRegister, addr uint16) int {
	if r.CarryFlag() {
		return branch(r, addr)
	}
	return 0
}

func beq(r *cpuRegister, addr uint16) int {
	if r.ZeroFlag() {
		return branch(r, addr)
	}
	return 0
}

func bmi(r *cpuRegister, addr uint16) int {
	if r.NegativeFlag() {
		return branch(r, addr)
	}
	return 0
}

func bne(r *cpuRegister, addr uint16) int {
	if !r.ZeroFlag() {
		return branch(r, addr)
	}
	return 0
}

func bpl(r *cpuRegister, addr uint16) int {
	if !r.NegativeFlag() {
		return branch(r, addr)
	}
	return 0
}

func bvc(r *cpuRegister, addr uint16) int {
	if !r.OverflowFlag() {
		return branch(r, addr)
	}
	return 0
}

func bvs(r *cpuRegister, addr uint16) int {
	if r.OverflowFlag() {
		return branch(r, addr)
	}
	return 0
}

func clc(r *cpuRegister) {
	r.SetCarryFlag(false)
}

func cld(r *cpuRegister) {
	r.SetDecimalFlag(false)
}

func cli(r *cpuRegister) {
	r.SetInterruptDisableFlag(false)
}

func clv(r *cpuRegister) {
	r.SetOverflowFlag(false)
}

func sec(r *cpuRegister) {
	r.SetCarryFlag(true)
}

func sed(r *cpuRegister) {
	r.SetDecimalFlag(true)
}

func sei(r *cpuRegister) {
	r.SetInterruptDisableFlag(true)
}

func brk(r *cpuRegister, m Memory) {
	push16(r, m, r.PC)
	push(r, m, r.P)
	r.SetInterruptDisableFlag(true)
	r.PC = read16(m, 0xFFFE)
}

func nmi(r *cpuRegister, m Memory) {
	r.SetBreakFlag(false)
	push16(r, m, r.PC)
	push(r, m, r.P)
	r.SetInterruptDisableFlag(true)
	r.PC = read16(m, 0xFFFA)
}

func irq(r *cpuRegister, m Memory) {
	if r.InterruptDisableFlag() {
		return
	}
	r.SetBreakFlag(false)
	push16(r, m, r.PC)
	push(r, m, r.P)
	r.SetInterruptDisableFlag(true)
	r.PC = read16(m, 0xFFFE)
}

func compare(r *cpuRegister, a byte, b byte) {
	r.SetCarryFlag(a >= b)
	r.UpdateNegativeFlag(a - b)
	r.UpdateZeroFlag(a - b)
}

func push(r *cpuRegister, m MemoryWriter, val byte) {
	m.Write(0x100|uint16(r.S), val)
	r.S--
}

func pop(r *cpuRegister, m MemoryReader) byte {
	r.S++
	return m.Read(0x100 | uint16(r.S))
}

func push16(r *cpuRegister, m MemoryWriter, val uint16) {
	l := byte(val & 0xFF)
	h := byte(val >> 8)
	push(r, m, h)
	push(r, m, l)
}

func pop16(r *cpuRegister, m MemoryReader) uint16 {
	l := pop(r, m)
	h := pop(r, m)
	return uint16(h)<<8 | uint16(l)
}

func branch(r *cpuRegister, addr uint16) int {
	cycle := 1
	if pagesCross(r.PC, addr) {
		cycle++
	}
	r.PC = addr
	return cycle
}

func pagesCross(a uint16, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}
