package apu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Envelope_Tick(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		e    *envelope
		want *envelope
	}{
		{
			"1",
			&envelope{
				start:             true,
				decayLevelCounter: 5,
				divider:           divider{counter: 3, period: 10},
			},
			&envelope{
				start:             false,
				decayLevelCounter: 15,
				divider:           divider{counter: 10, period: 10},
			},
		},
		{
			"2",
			&envelope{
				start:             false,
				loop:              true,
				decayLevelCounter: 5,
				divider:           divider{counter: 1, period: 10},
			},
			&envelope{
				start:             false,
				loop:              true,
				decayLevelCounter: 5,
				divider:           divider{counter: 0, period: 10},
			},
		},
		{
			"3",
			&envelope{
				start:             false,
				loop:              true,
				decayLevelCounter: 5,
				divider:           divider{counter: 0, period: 10},
			},
			&envelope{
				start:             false,
				loop:              true,
				decayLevelCounter: 4,
				divider:           divider{counter: 10, period: 10},
			},
		},
		{
			"4",
			&envelope{
				start:             false,
				loop:              true,
				decayLevelCounter: 0,
				divider:           divider{counter: 0, period: 10},
			},
			&envelope{
				start:             false,
				loop:              true,
				decayLevelCounter: 15,
				divider:           divider{counter: 10, period: 10},
			},
		},
		{
			"5",
			&envelope{
				start:             false,
				loop:              false,
				decayLevelCounter: 0,
				divider:           divider{counter: 0, period: 10},
			},
			&envelope{
				start:             false,
				loop:              false,
				decayLevelCounter: 0,
				divider:           divider{counter: 10, period: 10},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.e.tick()
			assert.Equal(t, tt.want, tt.e)
		})
	}
}
