package e2e_test

import (
	"image/color"
	"os"
	"testing"

	"github.com/ichirin2501/rgnes/nes"
	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/bus"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/cpu"
	"github.com/ichirin2501/rgnes/nes/memory"
	"github.com/ichirin2501/rgnes/nes/ppu"

	"github.com/stretchr/testify/assert"
)

type fakeRenderer struct {
}

func (f *fakeRenderer) Render(x, y int, c color.Color) {

}

func Test_InstrTestV5(t *testing.T) {
	tests := []struct {
		name    string
		rompath string
	}{
		{
			"01-basics.nes",
			"../roms/instr_test-v5/rom_singles/01-basics.nes",
		},
		{
			"02-implied.nes",
			"../roms/instr_test-v5/rom_singles/02-implied.nes",
		},
		{
			"03-immediate.nes",
			"../roms/instr_test-v5/rom_singles/03-immediate.nes",
		},
		{
			"04-zero_page.nes",
			"../roms/instr_test-v5/rom_singles/04-zero_page.nes",
		},
		{
			"05-zp_xy.nes",
			"../roms/instr_test-v5/rom_singles/05-zp_xy.nes",
		},
		{
			"06-absolute.nes",
			"../roms/instr_test-v5/rom_singles/06-absolute.nes",
		},
		{
			"07-abs_xy.nes",
			"../roms/instr_test-v5/rom_singles/07-abs_xy.nes",
		},
		{
			"08-ind_x.nes",
			"../roms/instr_test-v5/rom_singles/08-ind_x.nes",
		},
		{
			"09-ind_y.nes",
			"../roms/instr_test-v5/rom_singles/09-ind_y.nes",
		},
		{
			"10-branches.nes",
			"../roms/instr_test-v5/rom_singles/10-branches.nes",
		},
		{
			"11-stack.nes",
			"../roms/instr_test-v5/rom_singles/11-stack.nes",
		},
		{
			"12-jmp_jsr.nes",
			"../roms/instr_test-v5/rom_singles/12-jmp_jsr.nes",
		},
		{
			"13-rts.nes",
			"../roms/instr_test-v5/rom_singles/13-rts.nes",
		},
		{
			"14-rti.nes",
			"../roms/instr_test-v5/rom_singles/14-rti.nes",
		},
		{
			"15-brk.nes",
			"../roms/instr_test-v5/rom_singles/15-brk.nes",
		},
		{
			"16-special.nes",
			"../roms/instr_test-v5/rom_singles/16-special.nes",
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
			ram := memory.NewMemory(0x8100)
			irp := &cpu.Interrupter{}
			fake := &fakeRenderer{}
			ppu := ppu.NewPPU(fake, mapper, c.Mirror, irp, nil)
			apu := apu.NewAPU()
			joypad := nes.NewJoypad()
			cpuBus := bus.NewCPUBus(ram, ppu, apu, mapper, joypad)
			cpu := cpu.NewCPU(cpuBus, irp, nil)
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
