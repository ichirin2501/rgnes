package e2e_test

import (
	"os"
	"testing"

	"github.com/ichirin2501/rgnes/nes"
	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/cpu"
	"github.com/ichirin2501/rgnes/nes/ppu"
	"github.com/stretchr/testify/assert"
)

func Test_PPU_OUT_6000_By_blargg(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		rompath string
	}{
		{
			"oam_read/oam_read.nes",
			"../roms/oam_read/oam_read.nes",
		},
		{
			"oam_stress/oam_stress.nes",
			"../roms/oam_stress/oam_stress.nes",
		},
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
