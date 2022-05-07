package cpu

type StatusRegister byte

func (s *StatusRegister) updateBit(pos byte, val bool) {
	if val {
		*s |= (1 << pos)
	} else {
		*s &= ^(1 << pos)
	}
}

func (s *StatusRegister) Byte() byte {
	return byte(*s)
}

func (s *StatusRegister) IsCarry() bool {
	return (byte(*s) & (1 << 0)) == (1 << 0)
}

func (s *StatusRegister) SetCarry(val bool) {
	s.updateBit(0, val)
}

func (s *StatusRegister) IsZero() bool {
	return (byte(*s) & (1 << 1)) == (1 << 1)
}

func (s *StatusRegister) SetZero(val bool) {
	s.updateBit(1, val)
}

func (s *StatusRegister) IsInterruptDisable() bool {
	return (byte(*s) & (1 << 2)) == (1 << 2)
}

func (s *StatusRegister) SetInterruptDisable(val bool) {
	s.updateBit(2, val)
}

func (s *StatusRegister) IsDecimal() bool {
	return (byte(*s) & (1 << 3)) == (1 << 3)
}

func (s *StatusRegister) SetDecimal(val bool) {
	s.updateBit(3, val)
}

// TODO: break
func (s *StatusRegister) SetBreak1(val bool) {
	s.updateBit(4, val)
}

func (s *StatusRegister) SetBreak2(val bool) {
	s.updateBit(5, val)
}

func (s *StatusRegister) IsOverflow() bool {
	return (byte(*s) & (1 << 6)) == (1 << 6)
}

func (s *StatusRegister) SetOverflow(val bool) {
	s.updateBit(6, val)
}

func (s *StatusRegister) IsNegative() bool {
	return (byte(*s) & (1 << 7)) == (1 << 7)
}

func (s *StatusRegister) SetNegative(val bool) {
	s.updateBit(7, val)
}

func (s *StatusRegister) SetZN(val byte) {
	s.SetZero(val == 0x00)
	s.SetNegative(val&0x80 != 0)
}
