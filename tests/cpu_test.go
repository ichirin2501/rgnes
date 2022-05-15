package e2e_test

import (
	"image/color"
	"os"
	"testing"

	"github.com/ichirin2501/rgnes/nes"
	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/cpu"
	"github.com/ichirin2501/rgnes/nes/ppu"

	"github.com/stretchr/testify/assert"
)

type fakeRenderer struct {
}

func (f *fakeRenderer) Render(x, y int, c color.Color) {}
func (f *fakeRenderer) Refresh()                       {}

func Test_CPU_OUT_6000(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		rompath string
	}{
		{
			"instr_test-v5/rom_singles/01-basics.nes",
			"../roms/instr_test-v5/rom_singles/01-basics.nes",
		},
		{
			"instr_test-v5/rom_singles/02-implied.nes",
			"../roms/instr_test-v5/rom_singles/02-implied.nes",
		},
		{
			"instr_test-v5/rom_singles/03-immediate.nes",
			"../roms/instr_test-v5/rom_singles/03-immediate.nes",
		},
		{
			"instr_test-v5/rom_singles/04-zero_page.nes",
			"../roms/instr_test-v5/rom_singles/04-zero_page.nes",
		},
		{
			"instr_test-v5/rom_singles/05-zp_xy.nes",
			"../roms/instr_test-v5/rom_singles/05-zp_xy.nes",
		},
		{
			"instr_test-v5/rom_singles/06-absolute.nes",
			"../roms/instr_test-v5/rom_singles/06-absolute.nes",
		},
		{
			"instr_test-v5/rom_singles/07-abs_xy.nes",
			"../roms/instr_test-v5/rom_singles/07-abs_xy.nes",
		},
		{
			"instr_test-v5/rom_singles/08-ind_x.nes",
			"../roms/instr_test-v5/rom_singles/08-ind_x.nes",
		},
		{
			"instr_test-v5/rom_singles/09-ind_y.nes",
			"../roms/instr_test-v5/rom_singles/09-ind_y.nes",
		},
		{
			"instr_test-v5/rom_singles/10-branches.nes",
			"../roms/instr_test-v5/rom_singles/10-branches.nes",
		},
		{
			"instr_test-v5/rom_singles/11-stack.nes",
			"../roms/instr_test-v5/rom_singles/11-stack.nes",
		},
		{
			"instr_test-v5/rom_singles/12-jmp_jsr.nes",
			"../roms/instr_test-v5/rom_singles/12-jmp_jsr.nes",
		},
		{
			"instr_test-v5/rom_singles/13-rts.nes",
			"../roms/instr_test-v5/rom_singles/13-rts.nes",
		},
		{
			"instr_test-v5/rom_singles/14-rti.nes",
			"../roms/instr_test-v5/rom_singles/14-rti.nes",
		},
		{
			"instr_test-v5/rom_singles/15-brk.nes",
			"../roms/instr_test-v5/rom_singles/15-brk.nes",
		},
		{
			"instr_test-v5/rom_singles/16-special.nes",
			"../roms/instr_test-v5/rom_singles/16-special.nes",
		},
		{
			"instr_misc/rom_singles/01-abs_x_wrap.nes",
			"../roms/instr_misc/rom_singles/01-abs_x_wrap.nes",
		},
		{
			"instr_misc/rom_singles/02-branch_wrap.nes",
			"../roms/instr_misc/rom_singles/02-branch_wrap.nes",
		},
		{
			"instr_misc/rom_singles/03-dummy_reads.nes",
			"../roms/instr_misc/rom_singles/03-dummy_reads.nes",
		},
		// {
		// 	"instr_misc/rom_singles/04-dummy_reads_apu.nes",
		// 	"../roms/instr_misc/rom_singles/04-dummy_reads_apu.nes",
		// },
		{
			"cpu_dummy_writes/cpu_dummy_writes_ppumem.nes",
			"../roms/cpu_dummy_writes/cpu_dummy_writes_ppumem.nes",
		},
		{
			"cpu_dummy_writes/cpu_dummy_writes_oam.nes",
			"../roms/cpu_dummy_writes/cpu_dummy_writes_oam.nes",
		},
		{
			"cpu_exec_space/test_cpu_exec_space_ppuio.nes",
			"../roms/cpu_exec_space/test_cpu_exec_space_ppuio.nes",
		},
		// {
		// 	"cpu_exec_space/test_cpu_exec_space_apu.nes",
		// 	"../roms/cpu_exec_space/test_cpu_exec_space_apu.nes",
		// },
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
			irp := &cpu.Interrupter{}
			fake := &fakeRenderer{}
			ppu := ppu.NewPPU(fake, mapper, c.Mirror, irp, nil)
			apu := apu.NewAPU()
			joypad := nes.NewJoypad()
			cpuBus := cpu.NewBus(ppu, apu, mapper, joypad)
			cpu := cpu.NewCPU(cpuBus, irp, nil)
			cpu.Reset()
			for i := 0; i < 15; i++ {
				ppu.Step()
			}

			ready := false
			done := false
			for {
				cpu.Step()
				got := cpuBus.ReadForTest(0x6000)
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
