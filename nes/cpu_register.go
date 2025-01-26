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

func (s *processorStatus) byte() byte {
	return byte(*s)
}

func (s *processorStatus) isCarry() bool {
	return (byte(*s) & (1 << 0)) == (1 << 0)
}

func (s *processorStatus) setCarry(val bool) {
	s.updateBit(0, val)
}

func (s *processorStatus) isZero() bool {
	return (byte(*s) & (1 << 1)) == (1 << 1)
}

func (s *processorStatus) setZero(val bool) {
	s.updateBit(1, val)
}

func (s *processorStatus) isInterruptDisable() bool {
	return (byte(*s) & (1 << 2)) == (1 << 2)
}

func (s *processorStatus) setInterruptDisable(val bool) {
	s.updateBit(2, val)
}

func (s *processorStatus) isDecimal() bool {
	return (byte(*s) & (1 << 3)) == (1 << 3)
}

func (s *processorStatus) setDecimal(val bool) {
	s.updateBit(3, val)
}

// TODO: break
func (s *processorStatus) setBreak1(val bool) {
	s.updateBit(4, val)
}

func (s *processorStatus) setBreak2(val bool) {
	s.updateBit(5, val)
}

func (s *processorStatus) isOverflow() bool {
	return (byte(*s) & (1 << 6)) == (1 << 6)
}

func (s *processorStatus) setOverflow(val bool) {
	s.updateBit(6, val)
}

func (s *processorStatus) isNegative() bool {
	return (byte(*s) & (1 << 7)) == (1 << 7)
}

func (s *processorStatus) setNegative(val bool) {
	s.updateBit(7, val)
}

func (s *processorStatus) setZN(val byte) {
	s.setZero(val == 0x00)
	s.setNegative(val&0x80 != 0)
}
