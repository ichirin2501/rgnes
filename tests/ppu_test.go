package e2e_test

import (
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

func Test_OAM_Read(t *testing.T) {
	tests := []struct {
		name    string
		rompath string
	}{
		{
			"oam_read.nes",
			"../roms/oam_read/oam_read.nes",
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
			ram := memory.NewMemory(0x800)
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
