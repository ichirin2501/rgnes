package cpu

import "fmt"

type Trace struct {
	A                   byte
	X                   byte
	Y                   byte
	PC                  uint16
	S                   byte
	P                   byte
	ByteCode            []byte
	Opcode              opcode
	AddressingResult    string
	InstructionReadByte *byte
	Cycle               int
	PPUX                uint16
	PPUY                uint16
	PPUVBlankState      bool
}

func (t *Trace) NESTestString() string {
	bc := ""
	switch len(t.ByteCode) {
	case 1:
		bc = fmt.Sprintf("%02X      ", t.ByteCode[0])
	case 2:
		bc = fmt.Sprintf("%02X %02X   ", t.ByteCode[0], t.ByteCode[1])
	case 3:
		bc = fmt.Sprintf("%02X %02X %02X", t.ByteCode[0], t.ByteCode[1], t.ByteCode[2])
	}
	ar := ""
	if t.InstructionReadByte == nil {
		ar = t.AddressingResult
	} else {
		ar = fmt.Sprintf("%s = %02X", t.AddressingResult, *t.InstructionReadByte)
	}
	op := t.Opcode.Name.String()
	if t.Opcode.Unofficial {
		op = "*" + t.Opcode.Name.String()
	}

	// C000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:24 SP:FD PPU:  0, 45 CYC:15
	return fmt.Sprintf("%04X  %s %4s %-27s A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d",
		t.PC,
		bc,
		op,
		ar,
		t.A,
		t.X,
		t.Y,
		t.P,
		t.S,
		t.PPUY,
		t.PPUX,
		t.Cycle,
	)
}

func (t *Trace) SetCPURegisterA(v byte)            { t.A = v }
func (t *Trace) SetCPURegisterX(v byte)            { t.X = v }
func (t *Trace) SetCPURegisterY(v byte)            { t.Y = v }
func (t *Trace) SetCPURegisterPC(v uint16)         { t.PC = v }
func (t *Trace) SetCPURegisterS(v byte)            { t.S = v }
func (t *Trace) SetCPURegisterP(v byte)            { t.P = v }
func (t *Trace) SetCPUOpcode(v opcode)             { t.Opcode = v }
func (t *Trace) SetCPUAddressingResult(v string)   { t.AddressingResult = v }
func (t *Trace) SetCPUInstructionReadByte(v *byte) { t.InstructionReadByte = v }
func (t *Trace) SetPPUX(v uint16)                  { t.PPUX = v }
func (t *Trace) SetPPUY(v uint16)                  { t.PPUY = v }
func (t *Trace) SetPPUVBlankState(v bool)          { t.PPUVBlankState = v }
func (t *Trace) AddCPUCycle(v int)                 { t.Cycle += v }
func (t *Trace) AddCPUByteCode(v byte) {
	t.ByteCode = append(t.ByteCode, v)
}
func (t *Trace) Reset() {
	t.AddressingResult = ""
	t.InstructionReadByte = nil
	t.ByteCode = t.ByteCode[:0]
}
