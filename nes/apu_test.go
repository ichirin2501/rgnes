package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakePlayer struct{}

func (f *fakePlayer) Sample(v float32)    {}
func (f *fakePlayer) SampleRate() float64 { return 44100 }

func Test_Divider_Tick(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		d           *divider
		wantDivider *divider
		want        bool
	}{
		{
			"1",
			&divider{counter: 0},
			&divider{counter: 0},
			true,
		},
		{
			"2",
			&divider{counter: 0, period: 2},
			&divider{counter: 2, period: 2},
			true,
		},
		{
			"4",
			&divider{counter: 1, period: 2},
			&divider{counter: 0, period: 2},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.d.tick()
			assert.Equal(t, tt.wantDivider, tt.d)
			assert.Equal(t, tt.want, got)
		})
	}
}
func Test_APU_TickFrameCounter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		mode  byte
		steps int
		want  int
	}{
		{
			"1",
			0,
			29830,
			29829,
		},
		{
			"2",
			0,
			29830 + 1,
			0,
		},
		{
			"3",
			1,
			37282,
			37281,
		},
		{
			"4",
			1,
			37282 + 1,
			0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			irqLine := irqInterruptLine(0)
			apu := newAPU(&irqLine, &fakePlayer{}, &dma{})
			apu.frameMode = tt.mode

			for i := 0; i < tt.steps; i++ {
				apu.tickFrameCounter()
			}
			assert.Equal(t, tt.want, apu.frameStep)

			if tt.mode == 0 {
				for i := 0; i < 29830; i++ {
					apu.tickFrameCounter()
				}
			} else {
				for i := 0; i < 37282; i++ {
					apu.tickFrameCounter()
				}
			}
			assert.Equal(t, tt.want, apu.frameStep)
		})
	}
}
