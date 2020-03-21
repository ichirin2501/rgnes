package nes

func lda(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := m.Read(addr)
	r.SetA(v)
	r.UpdateZeroFlag(v)
	r.UpdateNegativeFlag(v)
	return 0
}

func ldx(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := m.Read(addr)
	r.SetX(v)
	r.UpdateZeroFlag(v)
	r.UpdateNegativeFlag(v)
	return 0
}

func ldy(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := m.Read(addr)
	r.SetY(v)
	r.UpdateZeroFlag(v)
	r.UpdateNegativeFlag(v)
	return 0
}

func sta(r cpuRegisterer, m MemoryWriter, addr uint16) int {
	m.Write(addr, r.A())
	return 0
}

func stx(r cpuRegisterer, m MemoryWriter, addr uint16) int {
	m.Write(addr, r.X())
	return 0
}

func sty(r cpuRegisterer, m MemoryWriter, addr uint16) int {
	m.Write(addr, r.Y())
	return 0
}

func tax(r cpuRegisterer) int {
	v := r.A()
	r.SetX(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func tay(r cpuRegisterer) int {
	v := r.A()
	r.SetY(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func tsx(r cpuRegisterer) int {
	v := r.S()
	r.SetX(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func txa(r cpuRegisterer) int {
	v := r.X()
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func txs(r cpuRegisterer) int {
	v := r.X()
	r.SetS(v)
	return 0
}

func tya(r cpuRegisterer) int {
	v := r.Y()
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func adc(r cpuRegisterer, m MemoryReader, addr uint16) int {
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
	return 0
}

func and(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := r.A() & m.Read(addr)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func aslAcc(r cpuRegisterer) int {
	r.SetCarryFlag((r.A() & 0x80) == 0x80)
	v := r.A() << 1
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func asl(r cpuRegisterer, m Memory, addr uint16) int {
	v := m.Read(addr)
	r.SetCarryFlag((v & 0x80) == 0x80)
	v <<= 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func bit(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := m.Read(addr)
	r.SetOverflowFlag((v & 0x40) == 0x40)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v & r.A())
	return 0
}

func cmp(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := m.Read(addr)
	compare(r, r.A(), v)
	return 0
}

func cpx(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := m.Read(addr)
	compare(r, r.X(), v)
	return 0
}

func cpy(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := m.Read(addr)
	compare(r, r.Y(), v)
	return 0
}

func dec(r cpuRegisterer, m Memory, addr uint16) int {
	v := m.Read(addr) - 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func dex(r cpuRegisterer) int {
	v := r.X() - 1
	r.SetX(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func dey(r cpuRegisterer) int {
	v := r.Y() - 1
	r.SetY(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func eor(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := r.A() ^ m.Read(addr)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func inc(r cpuRegisterer, m Memory, addr uint16) int {
	v := m.Read(addr) + 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func inx(r cpuRegisterer) int {
	v := r.X() + 1
	r.SetX(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func iny(r cpuRegisterer) int {
	v := r.Y() + 1
	r.SetY(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func lsrAcc(r cpuRegisterer) int {
	r.SetCarryFlag((r.A() & 1) == 1)
	v := r.A() >> 1
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func lsr(r cpuRegisterer, m Memory, addr uint16) int {
	v := m.Read(addr)
	r.SetCarryFlag((v & 1) == 1)
	v >>= 1
	m.Write(addr, v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func ora(r cpuRegisterer, m MemoryReader, addr uint16) int {
	v := r.A() | m.Read(addr)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func rolAcc(r cpuRegisterer) int {
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	r.SetCarryFlag((r.A() & 0x80) == 0x80)
	v := (r.A() << 1) | c
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func rol(r cpuRegisterer, m Memory, addr uint16) int {
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
	return 0
}

func rorAcc(r cpuRegisterer) int {
	c := byte(0)
	if r.CarryFlag() {
		c = 1
	}
	r.SetCarryFlag((r.A() & 1) == 1)
	v := (r.A() >> 1) | (c << 7)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func ror(r cpuRegisterer, m Memory, addr uint16) int {
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
	return 0
}

func sbc(r cpuRegisterer, m MemoryReader, addr uint16) int {
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
	return 0
}

func pha(r cpuRegisterer, m MemoryWriter) int {
	push(r, m, r.A())
	return 0
}

func php(r cpuRegisterer, m MemoryWriter) int {
	push(r, m, r.P()|breakFlagMask)
	return 0
}

func pla(r cpuRegisterer, m MemoryReader) int {
	v := pop(r, m)
	r.SetA(v)
	r.UpdateNegativeFlag(v)
	r.UpdateZeroFlag(v)
	return 0
}

func plp(r cpuRegisterer, m MemoryReader) int {
	v := pop(r, m)
	r.SetP(v)
	return 0
}

func jmp(r cpuRegisterer, addr uint16) int {
	r.SetPC(addr)
	return 0
}

func jsr(r cpuRegisterer, m MemoryWriter, addr uint16) int {
	push16(r, m, r.PC()-1)
	r.SetPC(addr)
	return 0
}

func rts(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func rti(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func bcc(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func bcs(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func beq(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func bmi(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func bne(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func bpl(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func bvc(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func bvs(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func clc(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func cld(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func cli(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func clv(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func sec(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func sed(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func sei(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }
func brk(r cpuRegisterer, m MemoryReader, addr uint16) int { return 0 }

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
