package nes

import (
	"fmt"
	"io"
)

type tracer struct {
	w io.Writer

	A                   byte
	X                   byte
	Y                   byte
	PC                  uint16
	S                   byte
	P                   byte
	byteCode            []byte
	opcode              opcode
	addressingResult    string
	instructionReadByte *byte
	ppux                uint16
	ppuy                uint16
}

func newTracer(w io.Writer) *tracer {
	return &tracer{
		w: w,
	}
}

func (t *tracer) print() {
	bc := ""
	switch len(t.byteCode) {
	case 1:
		bc = fmt.Sprintf("%02X      ", t.byteCode[0])
	case 2:
		bc = fmt.Sprintf("%02X %02X   ", t.byteCode[0], t.byteCode[1])
	case 3:
		bc = fmt.Sprintf("%02X %02X %02X", t.byteCode[0], t.byteCode[1], t.byteCode[2])
	}
	ar := ""
	if t.instructionReadByte == nil {
		ar = t.addressingResult
	} else {
		ar = fmt.Sprintf("%s = %02X", t.addressingResult, *t.instructionReadByte)
	}
	op := t.opcode.name.String()
	if t.opcode.unofficial {
		op = "*" + t.opcode.name.String()
	}

	// C000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:24 SP:FD PPU:  0, 45
	_, _ = fmt.Fprintf(t.w, "%04X  %s %4s %-27s A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d\n",
		t.PC,
		bc,
		op,
		ar,
		t.A,
		t.X,
		t.Y,
		t.P,
		t.S,
		t.ppuy,
		t.ppux,
	)
}

func (t *tracer) setCPURegisters(cpu *cpu) {
	t.A = cpu.A
	t.X = cpu.X
	t.Y = cpu.Y
	t.PC = cpu.PC
	t.S = cpu.S
	t.P = cpu.P.byte()
}

func (t *tracer) setCPUOpcode(v opcode)           { t.opcode = v }
func (t *tracer) setCPUAddressingResult(v string) { t.addressingResult = v }
func (t *tracer) setPPUX(v uint16)                { t.ppux = v }
func (t *tracer) setPPUY(v uint16)                { t.ppuy = v }
func (t *tracer) addCPUByteCode(v byte) {
	t.byteCode = append(t.byteCode, v)
}
func (t *tracer) reset() {
	t.addressingResult = ""
	t.instructionReadByte = nil
	t.byteCode = t.byteCode[:0]
}
