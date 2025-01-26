package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IRQInterruptSource(t *testing.T) {

	i := irqInterruptLine(0)
	assert.Equal(t, false, i.isLow())

	i.setLow(irqSourceFrameCounter)
	assert.Equal(t, int(irqSourceFrameCounter), int(i))
	assert.Equal(t, true, i.isLow())

	i.setLow(irqSourceDMC)
	assert.Equal(t, int(irqSourceDMC)|int(irqSourceFrameCounter), int(i))

	i.setHigh(irqSourceDMC)
	assert.Equal(t, int(irqSourceFrameCounter), int(i))
}
