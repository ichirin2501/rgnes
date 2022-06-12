package e2e_test

import (
	"os"
	"testing"

	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/cpu"
	"github.com/ichirin2501/rgnes/nes/joypad"
	"github.com/ichirin2501/rgnes/nes/ppu"
	"github.com/stretchr/testify/assert"
)

func Test_PPU_OUT_6000(t *testing.T) {
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
		{
			"ppu_vbl_nmi/rom_singles/01-vbl_basics.nes",
			"../roms/ppu_vbl_nmi/rom_singles/01-vbl_basics.nes",
		},
		{
			"ppu_vbl_nmi/rom_singles/02-vbl_set_time.nes",
			"../roms/ppu_vbl_nmi/rom_singles/02-vbl_set_time.nes",
		},
		{
			"ppu_vbl_nmi/rom_singles/03-vbl_clear_time.nes",
			"../roms/ppu_vbl_nmi/rom_singles/03-vbl_clear_time.nes",
		},
		{
			"ppu_vbl_nmi/rom_singles/04-nmi_control.nes",
			"../roms/ppu_vbl_nmi/rom_singles/04-nmi_control.nes",
		},
		{
			"ppu_vbl_nmi/rom_singles/05-nmi_timing.nes",
			"../roms/ppu_vbl_nmi/rom_singles/05-nmi_timing.nes",
		},
		{
			"ppu_vbl_nmi/rom_singles/06-suppression.nes",
			"../roms/ppu_vbl_nmi/rom_singles/06-suppression.nes",
		},
		{
			"ppu_vbl_nmi/rom_singles/07-nmi_on_timing.nes",
			"../roms/ppu_vbl_nmi/rom_singles/07-nmi_on_timing.nes",
		},
		{
			"ppu_vbl_nmi/rom_singles/08-nmi_off_timing.nes",
			"../roms/ppu_vbl_nmi/rom_singles/08-nmi_off_timing.nes",
		},
		{
			"ppu_vbl_nmi/rom_singles/09-even_odd_frames.nes",
			"../roms/ppu_vbl_nmi/rom_singles/09-even_odd_frames.nes",
		},
		// {
		// 	"ppu_vbl_nmi/rom_singles/10-even_odd_timing.nes",
		// 	"../roms/ppu_vbl_nmi/rom_singles/10-even_odd_timing.nes",
		// },
		{
			"ppu_read_buffer/test_ppu_read_buffer.nes",
			"../roms/ppu_read_buffer/test_ppu_read_buffer.nes",
		},
		{
			"ppu_open_bus/ppu_open_bus.nes",
			"../roms/ppu_open_bus/ppu_open_bus.nes",
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
			mapper, err := cassette.NewMapper(f)
			if err != nil {
				t.Fatal(err)
			}
			m := mapper.MirroingType()
			irp := &cpu.Interrupter{}
			fake := &fakeRenderer{}
			fakePlayer := &fakePlayer{}
			ppu := ppu.New(fake, mapper, &m, irp)
			apu := apu.New(irp, fakePlayer)
			joypad := joypad.New()
			cpuBus := cpu.NewBus(ppu, apu, mapper, joypad)
			cpu := cpu.New(cpuBus, irp)
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
