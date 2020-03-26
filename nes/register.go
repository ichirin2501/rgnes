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

type cpuRegister struct {
	A  byte   // Accumulator
	X  byte   // Index
	Y  byte   // Index
	PC uint16 // Program Counter
	S  byte   // Stack Pointer
	P  byte   // Status Register
}

func newCPURegister() *cpuRegister {
	return &cpuRegister{
		A:  0x00,
		X:  0x00,
		Y:  0x00,
		PC: 0x8000,
		S:  0xFD,
		P:  reservedFlagMask | interruptDisableFlagMask,
	}
}

func (r *cpuRegister) CarryFlag() bool {
	return (r.P & carryFlagMask) == carryFlagMask
}
func (r *cpuRegister) SetCarryFlag(cond bool) {
	if cond {
		r.P |= carryFlagMask
	} else {
		r.P &= ^carryFlagMask
	}
}
func (r *cpuRegister) ZeroFlag() bool {
	return (r.P & zeroFlagMask) == zeroFlagMask
}
func (r *cpuRegister) SetZeroFlag(cond bool) {
	if cond {
		r.P |= zeroFlagMask
	} else {
		r.P &= ^zeroFlagMask
	}
}
func (r *cpuRegister) UpdateZeroFlag(val byte) {
	r.SetZeroFlag(val == 0x00)
}
func (r *cpuRegister) InterruptDisableFlag() bool {
	return (r.P & interruptDisableFlagMask) == interruptDisableFlagMask
}
func (r *cpuRegister) SetInterruptDisableFlag(cond bool) {
	if cond {
		r.P |= interruptDisableFlagMask
	} else {
		r.P &= ^interruptDisableFlagMask
	}
}
func (r *cpuRegister) DecimalFlag() bool {
	return (r.P & decimalFlagMask) == decimalFlagMask
}
func (r *cpuRegister) SetDecimalFlag(cond bool) {
	if cond {
		r.P |= decimalFlagMask
	} else {
		r.P &= ^decimalFlagMask
	}
}
func (r *cpuRegister) BreakFlag() bool {
	return (r.P & breakFlagMask) == breakFlagMask
}
func (r *cpuRegister) SetBreakFlag(cond bool) {
	if cond {
		r.P |= breakFlagMask
	} else {
		r.P &= ^breakFlagMask
	}
}
func (r *cpuRegister) OverflowFlag() bool {
	return (r.P & overflowFlagMask) == overflowFlagMask
}
func (r *cpuRegister) SetOverflowFlag(cond bool) {
	if cond {
		r.P |= overflowFlagMask
	} else {
		r.P &= ^overflowFlagMask
	}
}
func (r *cpuRegister) NegativeFlag() bool {
	return (r.P & negativeFlagMask) == negativeFlagMask
}
func (r *cpuRegister) SetNegativeFlag(cond bool) {
	if cond {
		r.P |= negativeFlagMask
	} else {
		r.P &= ^negativeFlagMask
	}
}
func (r *cpuRegister) UpdateNegativeFlag(val byte) {
	r.SetNegativeFlag(val&0x80 != 0)
}

// TODO
type ppuRegister struct {
}

// TODO
func newPPURegister() *ppuRegister {
	return &ppuRegister{}
}
