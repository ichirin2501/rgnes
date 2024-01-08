package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
