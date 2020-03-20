package nes

const (
	carryFlagMask byte = (1 << iota)
	zeroFlagMask
	interruptDisableFlagMask
	decimalFlagMask
	breakFlagMask
	reservedFlagMask
	overflowFlagMask
	negativeFlagMask
)

type registerer interface {
	A() byte
	X() byte
	Y() byte
	PC() uint16
	S() byte
	P() byte
	SetA(val byte)
	SetX(val byte)
	SetY(val byte)
	SetPC(val uint16)
	SetS(val byte)
	SetP(val byte)
}

type register struct {
	a  byte   // Accumulator
	x  byte   // Index
	y  byte   // Index
	pc uint16 // Program Counter
	s  byte   // Stack Pointer
	p  byte   // Status Register
}

func (r *register) A() byte {
	return r.a
}
func (r *register) X() byte {
	return r.x
}
func (r *register) Y() byte {
	return r.y
}
func (r *register) PC() uint16 {
	return r.pc
}
func (r *register) S() byte {
	return r.s
}
func (r *register) P() byte {
	return r.p
}

func (r *register) SetA(val byte) {
	r.a = val
}
func (r *register) SetX(val byte) {
	r.x = val
}
func (r *register) SetY(val byte) {
	r.y = val
}
func (r *register) SetPC(val uint16) {
	r.pc = val
}
func (r *register) SetS(val byte) {
	r.s = val
}
func (r *register) SetP(val byte) {
	r.p = val
}

type cpuRegisterer interface {
	registerer
	CarryFlag() bool
	SetCarryFlag(cond bool)
	ZeroFlag() bool
	SetZeroFlag(cond bool)
	InterruptDisableFlag() bool
	SetInterruptDisableFlag(cond bool)
	DecimalFlag() bool
	SetDecimalFlag(cond bool)
	BreakFlag() bool
	SetBreakFlag(cond bool)
	OverflowFlag() bool
	SetOverflowFlag(cond bool)
	NegativeFlag() bool
	SetNegativeFlag(cond bool)

	UpdateNegativeFlag(val byte)
	UpdateZeroFlag(val byte)
}

type cpuRegister struct {
	*register
}

func newCPURegister() *cpuRegister {
	return &cpuRegister{&register{}}
}

func (r *cpuRegister) CarryFlag() bool {
	return (r.P() & carryFlagMask) == carryFlagMask
}
func (r *cpuRegister) SetCarryFlag(cond bool) {
	if cond {
		r.SetP(r.P() | carryFlagMask)
	} else {
		r.SetP(r.P() & ^carryFlagMask)
	}
}
func (r *cpuRegister) ZeroFlag() bool {
	return (r.P() & zeroFlagMask) == zeroFlagMask
}
func (r *cpuRegister) SetZeroFlag(cond bool) {
	if cond {
		r.SetP(r.P() | zeroFlagMask)
	} else {
		r.SetP(r.P() & ^zeroFlagMask)
	}
}
func (r *cpuRegister) UpdateZeroFlag(val byte) {
	r.SetZeroFlag(val == 0x00)
}
func (r *cpuRegister) InterruptDisableFlag() bool {
	return (r.P() & interruptDisableFlagMask) == interruptDisableFlagMask
}
func (r *cpuRegister) SetInterruptDisableFlag(cond bool) {
	if cond {
		r.SetP(r.P() | interruptDisableFlagMask)
	} else {
		r.SetP(r.P() & ^interruptDisableFlagMask)
	}
}
func (r *cpuRegister) DecimalFlag() bool {
	return (r.P() & decimalFlagMask) == decimalFlagMask
}
func (r *cpuRegister) SetDecimalFlag(cond bool) {
	if cond {
		r.SetP(r.P() | decimalFlagMask)
	} else {
		r.SetP(r.P() & ^decimalFlagMask)
	}
}
func (r *cpuRegister) BreakFlag() bool {
	return (r.P() & breakFlagMask) == breakFlagMask
}
func (r *cpuRegister) SetBreakFlag(cond bool) {
	if cond {
		r.SetP(r.P() | breakFlagMask)
	} else {
		r.SetP(r.P() & ^breakFlagMask)
	}
}
func (r *cpuRegister) OverflowFlag() bool {
	return (r.P() & overflowFlagMask) == overflowFlagMask
}
func (r *cpuRegister) SetOverflowFlag(cond bool) {
	if cond {
		r.SetP(r.P() | overflowFlagMask)
	} else {
		r.SetP(r.P() & ^overflowFlagMask)
	}
}
func (r *cpuRegister) NegativeFlag() bool {
	return (r.P() & negativeFlagMask) == negativeFlagMask
}
func (r *cpuRegister) SetNegativeFlag(cond bool) {
	if cond {
		r.SetP(r.P() | negativeFlagMask)
	} else {
		r.SetP(r.P() & ^negativeFlagMask)
	}
}
func (r *cpuRegister) UpdateNegativeFlag(val byte) {
	r.SetNegativeFlag(val&0x80 != 0)
}
