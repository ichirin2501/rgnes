package nes

/*
ref: https://www.nesdev.org/wiki/Status_flags

	7  bit  0
	---- ----
	NVss DIZC
	|||| ||||
	|||| |||+ - Carry
	|||| ||+- - Zero
	|||| |+-- - Interrupt Disable
	|||| +--- - Decimal
	||++----- - No CPU effect, see: the B flag
	|+------- - Overflow
	+-------- - Negative
*/
type processorStatus byte

func (s *processorStatus) updateBit(pos byte, val bool) {
	if val {
		*s |= (1 << pos)
	} else {
		*s &= ^(1 << pos)
	}
}

func (s *processorStatus) Byte() byte {
	return byte(*s)
}

func (s *processorStatus) IsCarry() bool {
	return (byte(*s) & (1 << 0)) == (1 << 0)
}

func (s *processorStatus) SetCarry(val bool) {
	s.updateBit(0, val)
}

func (s *processorStatus) IsZero() bool {
	return (byte(*s) & (1 << 1)) == (1 << 1)
}

func (s *processorStatus) SetZero(val bool) {
	s.updateBit(1, val)
}

func (s *processorStatus) IsInterruptDisable() bool {
	return (byte(*s) & (1 << 2)) == (1 << 2)
}

func (s *processorStatus) SetInterruptDisable(val bool) {
	s.updateBit(2, val)
}

func (s *processorStatus) IsDecimal() bool {
	return (byte(*s) & (1 << 3)) == (1 << 3)
}

func (s *processorStatus) SetDecimal(val bool) {
	s.updateBit(3, val)
}

// TODO: break
func (s *processorStatus) SetBreak1(val bool) {
	s.updateBit(4, val)
}

func (s *processorStatus) SetBreak2(val bool) {
	s.updateBit(5, val)
}

func (s *processorStatus) IsOverflow() bool {
	return (byte(*s) & (1 << 6)) == (1 << 6)
}

func (s *processorStatus) SetOverflow(val bool) {
	s.updateBit(6, val)
}

func (s *processorStatus) IsNegative() bool {
	return (byte(*s) & (1 << 7)) == (1 << 7)
}

func (s *processorStatus) SetNegative(val bool) {
	s.updateBit(7, val)
}

func (s *processorStatus) SetZN(val byte) {
	s.SetZero(val == 0x00)
	s.SetNegative(val&0x80 != 0)
}
