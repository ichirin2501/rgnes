package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IRQInterruptSource(t *testing.T) {

	i := IRQInterruptLine(0)
	assert.Equal(t, false, i.IsLow())

	i.SetLow(IRQSourceFrameCounter)
	assert.Equal(t, int(IRQSourceFrameCounter), int(i))
	assert.Equal(t, true, i.IsLow())

	i.SetLow(IRQSourceDMC)
	assert.Equal(t, int(IRQSourceDMC)|int(IRQSourceFrameCounter), int(i))

	i.SetHigh(IRQSourceDMC)
	assert.Equal(t, int(IRQSourceFrameCounter), int(i))
}
