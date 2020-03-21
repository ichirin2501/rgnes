package nes

func lda(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	r.SetA(v)
	r.UpdateZeroFlag(v)
	r.UpdateNegativeFlag(v)
}

func ldx(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	r.SetX(v)
	r.UpdateZeroFlag(v)
	r.UpdateNegativeFlag(v)
}

func ldy(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	r.SetY(v)
	r.UpdateZeroFlag(v)
	r.UpdateNegativeFlag(v)
}

func sta(r cpuRegisterer, m MemoryWriter, addr uint16) {
	m.Write(addr, r.A())
}

func stx(r cpuRegisterer, m MemoryWriter, addr uint16) {
	m.Write(addr, r.X())
}

func sty(r cpuRegisterer, m MemoryWriter, addr uint16) {
	m.Write(addr, r.Y())
}

func tax(r cpuRegisterer) {
	v := r.A()
	r.SetX(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func tay(r cpuRegisterer) {
	v := r.A()
	r.SetY(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func tsx(r cpuRegisterer) {
	v := r.S()
	r.SetX(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func txa(r cpuRegisterer) {
	v := r.X()
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func txs(r cpuRegisterer) {
	v := r.X()
	r.SetS(v)
}

func tya(r cpuRegisterer) {
	v := r.Y()
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func adc(r cpuRegisterer, m MemoryReader, addr uint16) {
	a := r.A()
	b := m.Read(addr)
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	v := a + b + c
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	r.SetCarryFlag(uint16(a)+uint16(b)+uint16(c) > 0xFF)
	r.SetOverflowFlag((a^b)&0x80 == 0 && (a^v)&0x80 != 0)
}

func and(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := r.A() & m.Read(addr)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func aslAcc(r cpuRegisterer) {
	r.SetCarryFlag((r.A() & 0x80) == 0x80)
	v := r.A() << 1
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func asl(r cpuRegisterer, m Memory, addr uint16) {
	v := m.Read(addr)
	r.SetCarryFlag((v & 0x80) == 0x80)
	v <<= 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func bit(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	r.SetOverflowFlag((v & 0x40) == 0x40)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v & r.A())
}

func cmp(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	compare(r, r.A(), v)
}

func cpx(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	compare(r, r.X(), v)
}

func cpy(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := m.Read(addr)
	compare(r, r.Y(), v)
}

func dec(r cpuRegisterer, m Memory, addr uint16) {
	v := m.Read(addr) - 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func dex(r cpuRegisterer) {
	v := r.X() - 1
	r.SetX(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func dey(r cpuRegisterer) {
	v := r.Y() - 1
	r.SetY(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func eor(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := r.A() ^ m.Read(addr)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func inc(r cpuRegisterer, m Memory, addr uint16) {
	v := m.Read(addr) + 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func inx(r cpuRegisterer) {
	v := r.X() + 1
	r.SetX(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func iny(r cpuRegisterer) {
	v := r.Y() + 1
	r.SetY(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func lsrAcc(r cpuRegisterer) {
	r.SetCarryFlag((r.A() & 1) == 1)
	v := r.A() >> 1
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func lsr(r cpuRegisterer, m Memory, addr uint16) {
	v := m.Read(addr)
	r.SetCarryFlag((v & 1) == 1)
	v >>= 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func ora(r cpuRegisterer, m MemoryReader, addr uint16) {
	v := r.A() | m.Read(addr)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func rolAcc(r cpuRegisterer) {
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	r.SetCarryFlag((r.A() & 0x80) == 0x80)
	v := (r.A() << 1) | c
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func rol(r cpuRegisterer, m Memory, addr uint16) {
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

func rorAcc(r cpuRegisterer) {
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	r.SetCarryFlag((r.A() & 1) == 1)
	v := (r.A() >> 1) | (c << 7)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func ror(r cpuRegisterer, m Memory, addr uint16) {
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

func sbc(r cpuRegisterer, m MemoryReader, addr uint16) {
	a := r.A()
	b := m.Read(addr)
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	v := a - b - (1 - c)
	r.SetA(v)
	r.SetCarryFlag(uint16(a)-uint16(b)-uint16(1-c) >= 0)
	r.SetOverflowFlag(((a^b)&0x80 != 0) && (a^v)&0x80 != 0)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func pha(r cpuRegisterer, m MemoryWriter) {
	push(r, m, r.A())
}

func php(r cpuRegisterer, m MemoryWriter) {
	push(r, m, r.P()|breakFlagMask)
}

func pla(r cpuRegisterer, m MemoryReader) {
	v := pop(r, m)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
}

func plp(r cpuRegisterer, m MemoryReader) {
	v := pop(r, m)
	r.SetP(v)
}

func jmp(r cpuRegisterer, addr uint16) {
	r.SetPC(addr)
}

func jsr(r cpuRegisterer, m MemoryWriter, addr uint16) {
	push16(r, m, r.PC()-1)
	r.SetPC(addr)
}

func rts(r cpuRegisterer, m MemoryReader) {
	v := pop16(r, m) + 1
	r.SetPC(v)
}

func rti(r cpuRegisterer, m MemoryReader) {
	r.SetP(pop(r, m))
	r.SetPC(pop16(r, m))
}

func bcc(r cpuRegisterer, addr uint16) int {
	if !r.CarryFlag() {
		return branch(r, addr)
	}
	return 0
}

func bcs(r cpuRegisterer, addr uint16) int {
	if r.CarryFlag() {
		return branch(r, addr)
	}
	return 0
}

func beq(r cpuRegisterer, addr uint16) int {
	if r.ZeroFlag() {
		return branch(r, addr)
	}
	return 0
}

func bmi(r cpuRegisterer, addr uint16) int {
	if r.NegativeFlag() {
		return branch(r, addr)
	}
	return 0
}

func bne(r cpuRegisterer, addr uint16) int {
	if !r.ZeroFlag() {
		return branch(r, addr)
	}
	return 0
}

func bpl(r cpuRegisterer, addr uint16) int {
	if !r.NegativeFlag() {
		return branch(r, addr)
	}
	return 0
}

func bvc(r cpuRegisterer, addr uint16) int {
	if !r.OverflowFlag() {
		return branch(r, addr)
	}
	return 0
}

func bvs(r cpuRegisterer, addr uint16) int {
	if r.OverflowFlag() {
		return branch(r, addr)
	}
	return 0
}

func clc(r cpuRegisterer) {
	r.SetCarryFlag(false)
}

func cld(r cpuRegisterer) {
	r.SetDecimalFlag(false)
}

func cli(r cpuRegisterer) {
	r.SetInterruptDisableFlag(false)
}

func clv(r cpuRegisterer) {
	r.SetOverflowFlag(false)
}

func sec(r cpuRegisterer) {
	r.SetCarryFlag(true)
}

func sed(r cpuRegisterer) {
	r.SetDecimalFlag(true)
}

func sei(r cpuRegisterer) {
	r.SetInterruptDisableFlag(true)
}

func brk(r cpuRegisterer, m Memory, addr uint16) {
	push16(r, m, r.PC())
	php(r, m)
	sei(r)
	r.SetPC(read16(m, 0xFFFE))
}

func compare(r cpuRegisterer, a byte, b byte) {
	r.SetCarryFlag(a >= b)
	r.UpdateNegativeFlag(a - b)
	r.UpdateZeroFlag(a - b)
}

func push(r registerer, m MemoryWriter, val byte) {
	m.Write(0x100|uint16(r.S()), val)
	r.SetS(r.S() - 1)
}

func pop(r registerer, m MemoryReader) byte {
	r.SetS(r.S() + 1)
	return m.Read(0x100 | uint16(r.S()))
}

func push16(r registerer, m MemoryWriter, val uint16) {
	l := byte(val & 0xFF)
	h := byte(val >> 8)
	push(r, m, h)
	push(r, m, l)
}

func pop16(r registerer, m MemoryReader) uint16 {
	l := pop(r, m)
	h := pop(r, m)
	return uint16(h)<<8 | uint16(l)
}

func branch(r registerer, addr uint16) int {
	cycle := 0
	pc := r.PC()
	r.SetPC(addr)
	if pagesCross(pc, addr) {
		cycle++
	}
	return cycle
}
