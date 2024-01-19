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
			"../nes-test-roms/apu_test/rom_singles/1-len_ctr.nes",
		},
		{
			"apu_test/rom_singles/2-len_table.nes",
			"../nes-test-roms/apu_test/rom_singles/2-len_table.nes",
		},
		{
			"apu_test/rom_singles/3-irq_flag.nes",
			"../nes-test-roms/apu_test/rom_singles/3-irq_flag.nes",
		},
		{
			"apu_test/rom_singles/4-jitter.nes",
			"../nes-test-roms/apu_test/rom_singles/4-jitter.nes",
		},
		{
			"apu_test/rom_singles/5-len_timing.nes",
			"../nes-test-roms/apu_test/rom_singles/5-len_timing.nes",
		},
		{
			"apu_test/rom_singles/6-irq_flag_timing.nes",
			"../nes-test-roms/apu_test/rom_singles/6-irq_flag_timing.nes",
		},
		// {
		// 	"apu_test/rom_singles/7-dmc_basics.nes",
		// 	"../nes-test-roms/apu_test/rom_singles/7-dmc_basics.nes",
		// },
		// {
		// 	"apu_test/rom_singles/8-dmc_rates.nes",
		// 	"../nes-test-roms/apu_test/rom_singles/8-dmc_rates.nes",
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

			n := nes.New(mapper, &fakeRenderer{}, &fakePlayer{})
			n.PowerUp()

			ready := false
			done := false
			for {
				n.Step()
				got := n.PeekMemory(0x6000)
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
