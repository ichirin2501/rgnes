package cpu

import (
	"image/color"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ichirin2501/rgnes/nes"
	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/bus"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/memory"
	"github.com/ichirin2501/rgnes/nes/ppu"
)

type fakeRenderer struct {
}

func (f *fakeRenderer) Render(x, y int, c color.Color) {

}

// func TestCPU(t *testing.T) {
// 	//os.Setenv("DEBUG", "1")
// 	buf := bytes.NewBuffer(make([]byte, 0))
// 	//debugWriter = buf
// 	defer func() {
// 		os.Setenv("DEBUG", "0")
// 		debugWriter = os.Stdout
// 	}()

// 	f, err := os.Open("testdata/nestest.nes")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer f.Close()

// 	c, err := cassette.NewCassette(f)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	mapper := cassette.NewMapper(c)
// 	cycle := NewCPUCycle()
// 	ram := memory.NewMemory(0x8100)

// 	irp := &Interrupter{}

// 	trace := &Trace{}

// 	fake := &fakeRenderer{}
// 	ppu := ppu.NewPPU(fake, c.CHR, mapper, c.Mirror, irp, trace)

// 	apu := apu.NewAPU()
// 	joypad := nes.NewJoypad()
// 	cpuBus := bus.NewCPUBus(ram, ppu, apu, mapper, joypad)
// 	cpu := NewCPU(cpuBus, cycle, irp, trace)

// 	cpu.Reset()
// 	trace.AddCPUCycle(7)
// 	for i := 0; i < 7*3+1; i++ {
// 		ppu.Step()
// 	}

// 	cpu.r.PC = 0xC000
// 	assert.Equal(t, byte(0), cpuBus.Read(0x02))
// 	assert.Equal(t, byte(0), cpuBus.Read(0x03))

// 	mp := make(map[int]struct{}, 0)
// 	for i := 0; i < 8991; i++ {
// 		cpu.t.Reset()

// 		cycle := cpu.Step()
// 		fmt.Printf("%s\r\n", trace.NESTestString())
// 		for k, v := range opcodeMap {
// 			if *v == trace.Opcode {
// 				mp[k] = struct{}{}
// 			}
// 		}

// 		trace.AddCPUCycle(cycle)
// 		for i := 0; i < cycle*3; i++ {
// 			ppu.Step()
// 		}
// 		// if cpuBus.Read(0x02) != 0 {
// 		// 	t.Fatal(fmt.Sprintf("0x02 is not 0 (0x%02x)", cpuBus.Read(0x02)))
// 		// }
// 		// if cpuBus.Read(0x03) != 0 {
// 		// 	t.Fatal(fmt.Sprintf("0x03 is not 0 (0x%02x)", cpuBus.Read(0x03)))
// 		// }
// 	}

// 	fmt.Println(mp)

// 	f2, err := os.Open("testdata/nestest-formatted2.log")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer f2.Close()

// 	expectedReader := bufio.NewReader(f2)

// 	for {
// 		got, err1 := buf.ReadString('\n')
// 		want, err2 := expectedReader.ReadString('\n')

// 		if err1 == io.EOF && err2 == io.EOF {
// 			break
// 		}
// 		//fmt.Println(got)
// 		if err1 != nil || err2 != nil {
// 			t.Fatal("don't eq")
// 		}
// 		assert.Equal(t, want, got)

// 	}

// 	// last check
// 	// http://www.qmtpro.com/~nes/misc/nestest.txt
// 	// This test program, when run on "automation", (i.e. set your program counter
// 	// to 0c000h) will perform all tests in sequence and shove the results of
// 	// the tests into locations 02h and 03h
// 	assert.Equal(t, byte(0), cpuBus.Read(0x02))
// 	assert.Equal(t, byte(0), cpuBus.Read(0x03))
// }

func Test_AddressingMode(t *testing.T) {
	tests := []struct {
		name            string
		op              *opcode
		r               *cpuRegister
		m               []byte // MemoryReader
		wantAddr        uint16
		wantPageCrossed bool
	}{
		{
			"absolute",
			&opcode{Mode: absolute},
			&cpuRegister{},
			[]byte{0x01, 0x20},
			0x2001,
			false,
		},
		{
			"zeroPage",
			&opcode{Mode: zeroPage},
			&cpuRegister{},
			[]byte{0x01, 0x20},
			0x0001,
			false,
		},
		{
			"zeroPageX/1",
			&opcode{Mode: zeroPageX},
			&cpuRegister{X: 0xFF},
			[]byte{0x01, 0x20},
			0x0000,
			false,
		},
		{
			"zeroPageX/2",
			&opcode{Mode: zeroPageX},
			&cpuRegister{X: 0x10},
			[]byte{0x01, 0x20},
			0x0011,
			false,
		},
		{
			"zeroPageY/1",
			&opcode{Mode: zeroPageY},
			&cpuRegister{Y: 0x20},
			[]byte{0x01, 0x20},
			0x0021,
			false,
		},
		{
			"relative/1",
			&opcode{Mode: relative},
			&cpuRegister{PC: 0x5},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00},
			0x0007,
			false,
		},
		{
			"relative/2",
			&opcode{Mode: relative},
			&cpuRegister{PC: 0x5},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00},
			0x0005,
			false,
		},
		{
			"absoluteX/1",
			&opcode{Mode: absoluteX},
			&cpuRegister{PC: 0x0, X: 0x1},
			[]byte{0xFF, 0x00},
			0x0100,
			true,
		},
		{
			"absoluteX/2",
			&opcode{Mode: absoluteX},
			&cpuRegister{PC: 0x0, X: 0x1},
			[]byte{0xFE, 0x01},
			0x01FF,
			false,
		},
		{
			"absoluteY/2",
			&opcode{Mode: absoluteY},
			&cpuRegister{PC: 0x0, Y: 0x1},
			[]byte{0xFE, 0x02},
			0x02FF,
			false,
		},
		{
			"indirect",
			&opcode{Mode: indirect},
			&cpuRegister{PC: 0x00},
			[]byte{0x03, 0x00, 0x00, 0x22, 0x33},
			0x3322,
			false,
		},
		// これの追加テストしたい
		{
			"indexedIndirect/1",
			&opcode{Mode: indexedIndirect},
			&cpuRegister{PC: 0x01, X: 0xFF},
			[]byte{0x00, 0x04, 0x00, 0x40, 0x01},
			0x0140,
			false,
		},
		{
			"indirectIndexed/1",
			&opcode{Mode: indirectIndexed},
			&cpuRegister{PC: 0x00, Y: 0x20},
			[]byte{0x02, 0x00, 0x01, 0x50, 0x00},
			0x5021,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.MemoryType(tt.m)

			cpuBus := bus.NewCPUBus(mem, nil, nil, nil, nil)
			cpu := &CPU{r: tt.r, m: cpuBus}
			gotAddr, gotPageCrossed := cpu.fetchOperand(tt.op)
			assert.Equal(t, tt.wantAddr, gotAddr)
			assert.Equal(t, tt.wantPageCrossed, gotPageCrossed)
		})
	}
}

func Test_InstrTestV5(t *testing.T) {
	tests := []struct {
		name    string
		rompath string
	}{
		{
			"01-basics.nes",
			"../../roms/instr_test-v5/rom_singles/01-basics.nes",
		},
		{
			"02-implied.nes",
			"../../roms/instr_test-v5/rom_singles/02-implied.nes",
		},
		// {
		// 	"03-immediate.nes",
		// 	"../../roms/instr_test-v5/rom_singles/03-immediate.nes",
		// },
		{
			"04-zero_page.nes",
			"../../roms/instr_test-v5/rom_singles/04-zero_page.nes",
		},
		{
			"05-zp_xy.nes",
			"../../roms/instr_test-v5/rom_singles/05-zp_xy.nes",
		},
		{
			"06-absolute.nes",
			"../../roms/instr_test-v5/rom_singles/06-absolute.nes",
		},
		// {
		// 	"07-abs_xy.nes",
		// 	"../../roms/instr_test-v5/rom_singles/07-abs_xy.nes",
		// },
		{
			"08-ind_x.nes",
			"../../roms/instr_test-v5/rom_singles/08-ind_x.nes",
		},
		{
			"09-ind_y.nes",
			"../../roms/instr_test-v5/rom_singles/09-ind_y.nes",
		},
		{
			"10-branches.nes",
			"../../roms/instr_test-v5/rom_singles/10-branches.nes",
		},
		{
			"11-stack.nes",
			"../../roms/instr_test-v5/rom_singles/11-stack.nes",
		},
		{
			"12-jmp_jsr.nes",
			"../../roms/instr_test-v5/rom_singles/12-jmp_jsr.nes",
		},
		{
			"13-rts.nes",
			"../../roms/instr_test-v5/rom_singles/13-rts.nes",
		},
		{
			"14-rti.nes",
			"../../roms/instr_test-v5/rom_singles/14-rti.nes",
		},
		{
			"15-brk.nes",
			"../../roms/instr_test-v5/rom_singles/15-brk.nes",
		},
		{
			"16-special.nes",
			"../../roms/instr_test-v5/rom_singles/16-special.nes",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.rompath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			c, err := cassette.NewCassette(f)
			if err != nil {
				t.Fatal(err)
			}
			mapper := cassette.NewMapper(c)
			cycle := NewCPUCycle()
			ram := memory.NewMemory(0x8100)
			irp := &Interrupter{}
			fake := &fakeRenderer{}
			ppu := ppu.NewPPU(fake, c.CHR, mapper, c.Mirror, irp, nil)
			apu := apu.NewAPU()
			joypad := nes.NewJoypad()
			cpuBus := bus.NewCPUBus(ram, ppu, apu, mapper, joypad)
			cpu := NewCPU(cpuBus, cycle, irp, nil)
			cpu.Reset()

			ready := false
			done := false
			for {
				cycle := cpu.Step()
				for i := 0; i < cycle*3; i++ {
					ppu.Step()
				}
				got := cpuBus.Read(0x6000)
				switch got {
				case 0x80:
					ready = true
				default:
					if ready {
						assert.Equal(t, uint8(0x00), got)
						done = true
					}
				}
				if done {
					break
				}
			}
		})
	}
}
