package e2e_test

import (
	"os"
	"testing"

	"github.com/ichirin2501/rgnes/nes"
	"github.com/stretchr/testify/assert"
)

func Test_APU_OUT_6000(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		rompath string
	}{
		{
			"apu_test/rom_singles/1-len_ctr.nes",
			"../roms/apu_test/rom_singles/1-len_ctr.nes",
		},
		{
			"apu_test/rom_singles/2-len_table.nes",
			"../roms/apu_test/rom_singles/2-len_table.nes",
		},
		{
			"apu_test/rom_singles/3-irq_flag.nes",
			"../roms/apu_test/rom_singles/3-irq_flag.nes",
		},
		{
			"apu_test/rom_singles/4-jitter.nes",
			"../roms/apu_test/rom_singles/4-jitter.nes",
		},
		{
			"apu_test/rom_singles/5-len_timing.nes",
			"../roms/apu_test/rom_singles/5-len_timing.nes",
		},
		{
			"apu_test/rom_singles/6-irq_flag_timing.nes",
			"../roms/apu_test/rom_singles/6-irq_flag_timing.nes",
		},
		// {
		// 	"apu_test/rom_singles/7-dmc_basics.nes",
		// 	"../roms/apu_test/rom_singles/7-dmc_basics.nes",
		// },
		// {
		// 	"apu_test/rom_singles/8-dmc_rates.nes",
		// 	"../roms/apu_test/rom_singles/8-dmc_rates.nes",
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
			mapper, err := nes.NewMapper(f)
			if err != nil {
				t.Fatal(err)
			}
			m := mapper.MirroingType()
			irp := &nes.Interrupter{}
			fake := &fakeRenderer{}
			fakePlayer := &fakePlayer{}
			ppu := nes.NewPPU(fake, mapper, m, irp)
			apu := nes.NewAPU(irp, fakePlayer)
			joypad := nes.NewJoypad()
			cpuBus := nes.NewBus(ppu, apu, mapper, joypad)
			cpu := nes.NewCPU(cpuBus, irp)
			cpu.PowerUp()

			ready := false
			done := false
			for {
				cpu.Step()
				got := cpuBus.Peek(0x6000)
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
