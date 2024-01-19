package e2e_test

import (
	"fmt"
	"image/color"
	"os"
	"testing"

	"github.com/ichirin2501/rgnes/nes"
	"github.com/stretchr/testify/assert"
)

type fakeRenderer struct{}

func (f *fakeRenderer) Render(x, y int, c color.Color) {}
func (f *fakeRenderer) Refresh()                       {}

type fakePlayer struct{}

func (f *fakePlayer) Sample(v float32)    {}
func (f *fakePlayer) SampleRate() float64 { return 44100 }

func Test_CPU_OUT_6000(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		rompath string
	}{
		{
			"instr_test-v5/rom_singles/01-basics.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/01-basics.nes",
		},
		{
			"instr_test-v5/rom_singles/02-implied.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/02-implied.nes",
		},
		{
			"instr_test-v5/rom_singles/03-immediate.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/03-immediate.nes",
		},
		{
			"instr_test-v5/rom_singles/04-zero_page.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/04-zero_page.nes",
		},
		{
			"instr_test-v5/rom_singles/05-zp_xy.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/05-zp_xy.nes",
		},
		{
			"instr_test-v5/rom_singles/06-absolute.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/06-absolute.nes",
		},
		{
			"instr_test-v5/rom_singles/07-abs_xy.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/07-abs_xy.nes",
		},
		{
			"instr_test-v5/rom_singles/08-ind_x.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/08-ind_x.nes",
		},
		{
			"instr_test-v5/rom_singles/09-ind_y.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/09-ind_y.nes",
		},
		{
			"instr_test-v5/rom_singles/10-branches.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/10-branches.nes",
		},
		{
			"instr_test-v5/rom_singles/11-stack.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/11-stack.nes",
		},
		{
			"instr_test-v5/rom_singles/12-jmp_jsr.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/12-jmp_jsr.nes",
		},
		{
			"instr_test-v5/rom_singles/13-rts.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/13-rts.nes",
		},
		{
			"instr_test-v5/rom_singles/14-rti.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/14-rti.nes",
		},
		{
			"instr_test-v5/rom_singles/15-brk.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/15-brk.nes",
		},
		{
			"instr_test-v5/rom_singles/16-special.nes",
			"../nes-test-roms/instr_test-v5/rom_singles/16-special.nes",
		},
		{
			"instr_misc/rom_singles/01-abs_x_wrap.nes",
			"../nes-test-roms/instr_misc/rom_singles/01-abs_x_wrap.nes",
		},
		{
			"instr_misc/rom_singles/02-branch_wrap.nes",
			"../nes-test-roms/instr_misc/rom_singles/02-branch_wrap.nes",
		},
		{
			"instr_misc/rom_singles/03-dummy_reads.nes",
			"../nes-test-roms/instr_misc/rom_singles/03-dummy_reads.nes",
		},
		// {
		// 	"instr_misc/rom_singles/04-dummy_reads_apu.nes",
		// 	"../nes-test-roms/instr_misc/rom_singles/04-dummy_reads_apu.nes",
		// },
		{
			"cpu_dummy_writes/cpu_dummy_writes_ppumem.nes",
			"../nes-test-roms/cpu_dummy_writes/cpu_dummy_writes_ppumem.nes",
		},
		{
			"cpu_dummy_writes/cpu_dummy_writes_oam.nes",
			"../nes-test-roms/cpu_dummy_writes/cpu_dummy_writes_oam.nes",
		},
		{
			"cpu_exec_space/test_cpu_exec_space_ppuio.nes",
			"../nes-test-roms/cpu_exec_space/test_cpu_exec_space_ppuio.nes",
		},
		// {
		// 	"cpu_exec_space/test_cpu_exec_space_apu.nes",
		// 	"../nes-test-roms/cpu_exec_space/test_cpu_exec_space_apu.nes",
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

func Test_NESTest(t *testing.T) {
	// https://www.qmtpro.com/~nes/misc/nestest.txt
	// > This test program, when run on "automation", (i.e. set your program counter
	// > to 0c000h) will perform all tests in sequence and shove the results of
	// > the tests into locations 02h and 03h.
	f, err := os.Open("../nes-test-roms/other/nestest.nes")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	mapper, err := nes.NewMapper(f)
	if err != nil {
		t.Fatal(err)
	}

	n := nes.New(mapper, &fakeRenderer{}, &fakePlayer{})
	n.SetCPUPC(0xC000)

	assert.Equal(t, byte(0), n.PeekMemory(0x02))
	assert.Equal(t, byte(0), n.PeekMemory(0x03))
	for i := 0; i < 8991; i++ {
		n.Step()
		if n.PeekMemory(0x02) != 0 {
			t.Fatal(fmt.Sprintf("0x02 is not 0 (0x%02x)", n.PeekMemory(0x02)))
		}
		if n.PeekMemory(0x03) != 0 {
			t.Fatal(fmt.Sprintf("0x03 is not 0 (0x%02x)", n.PeekMemory(0x03)))
		}
	}
	assert.Equal(t, byte(0), n.PeekMemory(0x02))
	assert.Equal(t, byte(0), n.PeekMemory(0x03))
}
