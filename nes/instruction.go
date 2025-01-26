package nes

func (cpu *cpu) lda(addr uint16) {
	a := cpu.read(addr)
	cpu.A = a
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) ldx(addr uint16) {
	a := cpu.read(addr)
	cpu.X = a
	cpu.P.setZN(cpu.X)
}

func (cpu *cpu) ldy(addr uint16) {
	a := cpu.read(addr)
	cpu.Y = a
	cpu.P.setZN(cpu.Y)
}

func (cpu *cpu) sta(addr uint16) {
	cpu.write(addr, cpu.A)
}

func (cpu *cpu) stx(addr uint16) {
	cpu.write(addr, cpu.X)
}

func (cpu *cpu) sty(addr uint16) {
	cpu.write(addr, cpu.Y)
}

func (cpu *cpu) tax() {
	cpu.X = cpu.A
	cpu.P.setZN(cpu.X)
}

func (cpu *cpu) tay() {
	cpu.Y = cpu.A
	cpu.P.setZN(cpu.Y)
}

func (cpu *cpu) tsx() {
	cpu.X = cpu.S
	cpu.P.setZN(cpu.X)
}

func (cpu *cpu) txa() {
	cpu.A = cpu.X
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) txs() {
	cpu.S = cpu.X
}

func (cpu *cpu) tya() {
	cpu.A = cpu.Y
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) adc(addr uint16) {
	a := cpu.A
	b := cpu.read(addr)
	c := byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	v := a + b + c
	cpu.A = v
	cpu.P.setZN(v)
	cpu.P.setCarry(uint16(a)+uint16(b)+uint16(c) > 0xFF)
	cpu.P.setOverflow((a^b)&0x80 == 0 && (a^v)&0x80 != 0)
}

func (cpu *cpu) and(addr uint16) {
	cpu.A &= cpu.read(addr)
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) aslAcc() {
	cpu.P.setCarry((cpu.A & 0x80) == 0x80)
	cpu.A <<= 1
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) asl(addr uint16) {
	v := cpu.read(addr)
	cpu.P.setCarry((v & 0x80) == 0x80)
	cpu.write(addr, v) // dummy write
	v <<= 1
	cpu.write(addr, v)
	cpu.P.setZN(v)
}

func (cpu *cpu) bit(addr uint16) {
	v := cpu.read(addr)
	cpu.P.setOverflow((v & 0x40) == 0x40)
	cpu.P.setNegative(v&0x80 != 0)
	cpu.P.setZero((v & cpu.A) == 0x00)
}

func (cpu *cpu) cmp(addr uint16) {
	v := cpu.read(addr)
	cpu.compare(cpu.A, v)
}

func (cpu *cpu) cpx(addr uint16) {
	v := cpu.read(addr)
	cpu.compare(cpu.X, v)
}

func (cpu *cpu) cpy(addr uint16) {
	v := cpu.read(addr)
	cpu.compare(cpu.Y, v)
}

func (cpu *cpu) dec(addr uint16) {
	v := cpu.read(addr)
	cpu.write(addr, v) // dummy write
	v--
	cpu.write(addr, v)
	cpu.P.setZN(v)
}

func (cpu *cpu) dex() {
	cpu.X--
	cpu.P.setZN(cpu.X)
}

func (cpu *cpu) dey() {
	cpu.Y--
	cpu.P.setZN(cpu.Y)
}

func (cpu *cpu) eor(addr uint16) {
	cpu.A ^= cpu.read(addr)
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) inc(addr uint16) {
	v := cpu.read(addr)
	cpu.write(addr, v) // dummy write
	v++
	cpu.write(addr, v)
	cpu.P.setZN(v)
}

func (cpu *cpu) inx() {
	cpu.X++
	cpu.P.setZN(cpu.X)
}

func (cpu *cpu) iny() {
	cpu.Y++
	cpu.P.setZN(cpu.Y)
}

func (cpu *cpu) lsrAcc() {
	cpu.P.setCarry((cpu.A & 1) == 1)
	cpu.A >>= 1
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) lsr(addr uint16) {
	v := cpu.read(addr)
	cpu.write(addr, v) // dummy write
	cpu.P.setCarry((v & 1) == 1)
	v >>= 1
	cpu.write(addr, v)
	cpu.P.setZN(v)
}

func (cpu *cpu) ora(addr uint16) {
	cpu.A |= cpu.read(addr)
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) rolAcc() {
	c := byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	cpu.P.setCarry((cpu.A & 0x80) == 0x80)
	cpu.A = (cpu.A << 1) | c
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) rol(addr uint16) {
	c := byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	v := cpu.read(addr)
	cpu.write(addr, v) // dummy write
	cpu.P.setCarry((v & 0x80) == 0x80)
	v = (v << 1) | c
	cpu.write(addr, v)
	cpu.P.setZN(v)
}

func (cpu *cpu) rorAcc() {
	c := byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	cpu.P.setCarry((cpu.A & 1) == 1)
	cpu.A = (cpu.A >> 1) | (c << 7)
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) ror(addr uint16) {
	c := byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	v := cpu.read(addr)
	cpu.write(addr, v) // dummy write
	cpu.P.setCarry((v & 1) == 1)
	v = (v >> 1) | (c << 7)
	cpu.write(addr, v)
	cpu.P.setZN(v)
}

func (cpu *cpu) sbc(addr uint16) {
	a := cpu.A
	b := cpu.read(addr)
	c := byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	v := a - b - (1 - c)
	cpu.A = v
	cpu.P.setCarry(int(a)-int(b)-int(1-c) >= 0)
	cpu.P.setOverflow(((a^b)&0x80 != 0) && (a^v)&0x80 != 0)
	cpu.P.setZN(v)
}

/*
ref: https://www.nesdev.org/6502_cpu.txt

	PHA, PHP
	#  address R/W description
	--- ------- --- -----------------------------------------------
	1    PC     R  fetch opcode, increment PC
	2    PC     R  read next instruction byte (and throw it away)
	3  $0100,S  W  push register on stack, decrement S

ここはアドレッシングモード側でdummy readしてカバーされているのでdummy readはたぶん不要
*/
func (cpu *cpu) pha() {
	//cpu.Read(cpu.PC) // dummy read
	cpu.push(cpu.A)
}
func (cpu *cpu) php() {
	//cpu.Read(cpu.PC) // dummy read
	cpu.push(cpu.P.byte() | 0x30)
}

/*
ref: https://www.nesdev.org/6502_cpu.txt

	PLA, PLP
	#  address R/W description
	--- ------- --- -----------------------------------------------
	1    PC     R  fetch opcode, increment PC
	2    PC     R  read next instruction byte (and throw it away)
	3  $0100,S  R  increment S
	4  $0100,S  R  pull register from stack
*/
func (cpu *cpu) pla() {
	cpu.read(cpu.PC) // dummy read
	cpu.A = cpu.pop()
	cpu.P.setZN(cpu.A)
}
func (cpu *cpu) plp() {
	cpu.read(cpu.PC) // dummy read
	cpu.P = processorStatus((cpu.pop() & 0xEF) | (1 << 5))
}

func (cpu *cpu) jmp(addr uint16) {
	cpu.PC = addr
}

/*
ref: https://www.nesdev.org/6502_cpu.txt

	JSR
	#  address R/W description
	--- ------- --- -------------------------------------------------
	1    PC     R  fetch opcode, increment PC
	2    PC     R  fetch low address byte, increment PC
	3  $0100,S  R  internal operation (predecrement S?)
	4  $0100,S  W  push PCH on stack, decrement S
	5  $0100,S  W  push PCL on stack, decrement S
	6    PC     R  copy low address byte to PCL, fetch high address

	byte to PCH
*/
func (cpu *cpu) jsr(addr uint16) {
	cpu.push16(cpu.PC - 1)
	// dummy read
	cpu.read(cpu.PC)
	cpu.PC = addr
}

/*
ref: https://www.nesdev.org/6502_cpu.txt

	RTS
	#  address R/W description
	--- ------- --- -----------------------------------------------
	1    PC     R  fetch opcode, increment PC
	2    PC     R  read next instruction byte (and throw it away)
	3  $0100,S  R  increment S
	4  $0100,S  R  pull PCL from stack, increment S
	5  $0100,S  R  pull PCH from stack
	6    PC     R  increment PC
*/
func (cpu *cpu) rts() {
	// 3  $0100,S  R  increment S
	cpu.read(0x100 | uint16(cpu.S)) // dummy read

	cpu.PC = cpu.pop16()
	// 6    PC     R  increment PC
	cpu.read(cpu.PC) // dummy read

	cpu.PC++
}

/*
ref: https://www.nesdev.org/6502_cpu.txt

	RTI
	#  address R/W description
	--- ------- --- -----------------------------------------------
	1    PC     R  fetch opcode, increment PC
	2    PC     R  read next instruction byte (and throw it away)
	3  $0100,S  R  increment S
	4  $0100,S  R  pull P from stack, increment S
	5  $0100,S  R  pull PCL from stack, increment S
	6  $0100,S  R  pull PCH from stack
*/
func (cpu *cpu) rti() {
	cpu.read(0x100 | uint16(cpu.S)) // dummy read
	cpu.P = processorStatus((cpu.pop() & 0xEF) | (1 << 5))
	cpu.PC = cpu.pop16()
}

func (cpu *cpu) bcc(addr uint16) int {
	if !cpu.P.isCarry() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *cpu) bcs(addr uint16) int {
	if cpu.P.isCarry() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *cpu) beq(addr uint16) int {
	if cpu.P.isZero() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *cpu) bmi(addr uint16) int {
	if cpu.P.isNegative() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *cpu) bne(addr uint16) int {
	if !cpu.P.isZero() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *cpu) bpl(addr uint16) int {
	if !cpu.P.isNegative() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *cpu) bvc(addr uint16) int {
	if !cpu.P.isOverflow() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *cpu) bvs(addr uint16) int {
	if cpu.P.isOverflow() {
		return cpu.branch(addr)
	}
	return 0
}

func (cpu *cpu) clc() {
	cpu.P.setCarry(false)
}

func (cpu *cpu) cld() {
	cpu.P.setDecimal(false)
}

func (cpu *cpu) cli() {
	cpu.P.setInterruptDisable(false)
}

func (cpu *cpu) clv() {
	cpu.P.setOverflow(false)
}

func (cpu *cpu) sec() {
	cpu.P.setCarry(true)
}

func (cpu *cpu) sed() {
	cpu.P.setDecimal(true)
}

func (cpu *cpu) sei() {
	cpu.P.setInterruptDisable(true)
}

/*
BRK

	#  address R/W description
	--- ------- --- -----------------------------------------------
	1    PC     R  fetch opcode, increment PC
	2    PC     R  read next instruction byte (and throw it away),

	increment PC

	3  $0100,S  W  push PCH on stack (with B flag set), decrement S
	4  $0100,S  W  push PCL on stack, decrement S
	5  $0100,S  W  push P on stack, decrement S
	6   $FFFE   R  fetch PCL
	7   $FFFF   R  fetch PCH

The addressing mode of the BRK instruction is the Implied.
The first 2 CPU clocks are already working in advance addressing mode processing, so dummy read is not necessary.
*/
func (cpu *cpu) brk() {
	cpu.interrupting = true
	cpu.push16(cpu.PC + 1)

	// ref: https://www.nesdev.org/wiki/CPU_interrupts#Interrupt_hijacking
	// > For example, if NMI is asserted during the first four ticks of a BRK instruction,
	// > the BRK instruction will execute normally at first (PC increments will occur and the status word will be pushed with the B flag set),
	// > but execution will branch to the NMI vector instead of the IRQ/BRK vector:

	// > Each [] is a CPU tick. [...] is whatever tick precedes the BRK opcode fetch.
	// > Asserting NMI during the interval marked with * causes a branch to the NMI routine instead of the IRQ/BRK routine:
	// >      ********************
	// > [...][BRK][BRK][BRK][BRK][BRK][BRK][BRK]

	// It's not clear to me what timing "Asserting NMI" refers to.
	// Is it the timing when the NMI edge detector detects it, the timing when the NMI internal signal turns ON,
	// or the timing when the decision is made to actually execute the NMI?
	// As far as I checked with cpu_interrupts_v2/2-nmi_and_brk.nes,
	// it seems to be expecting the NMI edge detector timing.
	if cpu.nmiSignal || cpu.nmiTriggered {
		cpu.push(cpu.P.byte() | 0x30)
		cpu.P.setInterruptDisable(true)
		cpu.PC = cpu.read16(0xFFFA)
		cpu.clearNMIInterruptState()
	} else {
		cpu.push(cpu.P.byte() | 0x30)
		cpu.P.setInterruptDisable(true)
		cpu.PC = cpu.read16(0xFFFE)
	}
	cpu.interrupting = false
}

/*
Two interrupts (/IRQ and /NMI) and two instructions (PHP and BRK) push the flags to the stack.
In the byte pushed, bit 5 is always set to 1, and bit 4 is 1 if from an instruction (PHP or BRK) or 0 if from an interrupt line being pulled low (/IRQ or /NMI).

ref: https://www.nesdev.org/wiki/CPU_interrupts#IRQ_and_NMI_tick-by-tick_execution

	#  address R/W description
	--- ------- --- -----------------------------------------------

	1    PC     R  fetch opcode (and discard it - $00 (BRK) is forced into the opcode register instead)
	2    PC     R  read next instruction byte (actually the same as above, since PC increment is suppressed. Also discarded.)
	3  $0100,S  W  push PCH on stack, decrement S
	4  $0100,S  W  push PCL on stack, decrement S

	*** At this point, the signal status determines which interrupt vector is used ***

	5  $0100,S  W  push P on stack (with B flag *clear*), decrement S
	6   A       R  fetch PCL (A = FFFE for IRQ, A = FFFA for NMI), set I flag
	7   A       R  fetch PCH (A = FFFF for IRQ, A = FFFB for NMI)
*/
func (cpu *cpu) nmi() {
	cpu.interrupting = true

	cpu.read(cpu.PC) // dummy read
	cpu.read(cpu.PC) // dummy read
	cpu.push16(cpu.PC)
	cpu.push(cpu.P.byte() | 0x20)
	cpu.P.setInterruptDisable(true)
	cpu.PC = cpu.read16(0xFFFA)

	cpu.clearNMIInterruptState()
	cpu.interrupting = false
}

func (cpu *cpu) irq() {
	cpu.interrupting = true
	cpu.read(cpu.PC) // dummy read
	cpu.read(cpu.PC) // dummy read
	cpu.push16(cpu.PC)
	if cpu.nmiSignal || cpu.nmiTriggered {
		cpu.push(cpu.P.byte() | 0x20)
		cpu.P.setInterruptDisable(true)
		cpu.PC = cpu.read16(0xFFFA)
		cpu.clearNMIInterruptState()
	} else {
		cpu.push(cpu.P.byte() | 0x20)
		cpu.P.setInterruptDisable(true)
		cpu.PC = cpu.read16(0xFFFE)
	}
	cpu.interrupting = false
}

func (cpu *cpu) compare(a byte, b byte) {
	cpu.P.setZN(a - b)
	cpu.P.setCarry(a >= b)
}

func (cpu *cpu) push(val byte) {
	cpu.write(0x100|uint16(cpu.S), val)
	cpu.S--
}

func (cpu *cpu) pop() byte {
	cpu.S++
	return cpu.read(0x100 | uint16(cpu.S))
}

func (cpu *cpu) push16(val uint16) {
	l := byte(val & 0xFF)
	h := byte(val >> 8)
	cpu.push(h)
	cpu.push(l)
}

func (cpu *cpu) pop16() uint16 {
	l := cpu.pop()
	h := cpu.pop()
	return uint16(h)<<8 | uint16(l)
}

func (cpu *cpu) branch(addr uint16) int {
	cycle := 1
	cpu.read(cpu.PC) // dummy read
	if pagesCross(cpu.PC, addr) {
		cpu.read(cpu.PC) // dummy read
		cycle++
	}
	cpu.PC = addr
	return cycle
}

func pagesCross(a uint16, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

// undocumented opcode

func (cpu *cpu) kil() {}

func (cpu *cpu) slo(addr uint16) {
	v := cpu.read(addr)
	cpu.write(addr, v) // dummy write
	cpu.P.setCarry((v & 0x80) == 0x80)
	v <<= 1
	cpu.write(addr, v)

	cpu.A |= v
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) anc(addr uint16) {
	a := cpu.read(addr)
	cpu.A &= a
	cpu.P.setZN(cpu.A)
	cpu.P.setCarry(cpu.P.isNegative())
}

func (cpu *cpu) rla(addr uint16) {
	c := byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	v := cpu.read(addr)
	cpu.write(addr, v) // dummy write
	cpu.P.setCarry((v & 0x80) == 0x80)
	v = (v << 1) | c
	cpu.write(addr, v)

	cpu.A &= v
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) sre(addr uint16) {
	v := cpu.read(addr)
	cpu.write(addr, v) // dummy write
	cpu.P.setCarry((v & 1) == 1)
	v >>= 1
	cpu.write(addr, v)

	cpu.A ^= v
	cpu.P.setZN(cpu.A)
}

func (cpu *cpu) alr(addr uint16) {
	// A =(A&#{imm})/2
	// N, Z, C
	v := cpu.A & cpu.read(addr)
	cpu.P.setCarry((v & 1) == 1)
	v >>= 1
	cpu.P.setZN(v)
	cpu.A = v
}

func (cpu *cpu) rra(addr uint16) {
	c := byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	k := cpu.read(addr)
	cpu.write(addr, k) // dummy write
	cpu.P.setCarry((k & 1) == 1)
	k = (k >> 1) | (c << 7)
	cpu.write(addr, k)

	a := cpu.A
	b := k
	c = byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	v := a + b + c
	cpu.A = v
	cpu.P.setZN(v)
	cpu.P.setCarry(uint16(a)+uint16(b)+uint16(c) > 0xFF)
	cpu.P.setOverflow((a^b)&0x80 == 0 && (a^v)&0x80 != 0)
}

func (cpu *cpu) arr(addr uint16) {
	// A:=(A&#{imm})/2
	// N V Z C
	// N and Z are normal, but C is bit 6 and V is bit 6 xor bit 5.
	c := byte(0)
	if cpu.P.isCarry() {
		c = 0x80
	}
	v := ((cpu.A & cpu.read(addr)) >> 1) | c
	cpu.P.setZN(v)
	cpu.P.setCarry((v & 0x40) == 0x40)
	cpu.P.setOverflow(((v & 0x40) ^ ((v & 0x20) << 1)) == 0x40)
	cpu.A = v
}

func (cpu *cpu) sax(addr uint16) {
	cpu.write(addr, cpu.A&cpu.X)
}

func (cpu *cpu) xaa() {}

func (cpu *cpu) ahx() {}

func (cpu *cpu) tas() {}

func (cpu *cpu) shy(addr uint16) {
	if pagesCross(addr, addr-uint16(cpu.X)) {
		addr &= uint16(cpu.Y) << 8
	}
	res := cpu.Y & (byte(addr>>8) + 1)
	cpu.write(addr, res)
}

func (cpu *cpu) shx(addr uint16) {
	if pagesCross(addr, addr-uint16(cpu.Y)) {
		addr &= uint16(cpu.X) << 8
	}
	res := cpu.X & (byte(addr>>8) + 1)
	cpu.write(addr, res)
}

func (cpu *cpu) lax(addr uint16) {
	v := cpu.read(addr)
	cpu.X = v
	cpu.A = v
	cpu.P.setZN(v)
}

func (cpu *cpu) las() {}

func (cpu *cpu) dcp(addr uint16) {
	v := cpu.read(addr)
	cpu.write(addr, v) // dummy write
	v--
	cpu.compare(cpu.A, v)
	cpu.write(addr, v)
}

func (cpu *cpu) axs(addr uint16) {
	// X:=A&X-#{imm}
	// N, Z, C
	t := cpu.A & cpu.X
	v := cpu.read(addr)
	cpu.X = t - v
	cpu.P.setCarry(int(t)-int(v) >= 0)
	cpu.P.setZN(cpu.X)
}

func (cpu *cpu) isb(addr uint16) {
	k := cpu.read(addr)
	cpu.write(addr, k) // dummy write
	k++
	cpu.write(addr, k)

	a := cpu.A
	b := k
	c := byte(0)
	if cpu.P.isCarry() {
		c = 1
	}
	v := a - b - (1 - c)
	cpu.A = v
	cpu.P.setCarry(int(a)-int(b)-int(1-c) >= 0)
	cpu.P.setOverflow(((a^b)&0x80 != 0) && (a^v)&0x80 != 0)
	cpu.P.setZN(v)
}
