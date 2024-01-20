package nes

// https://www.nesdev.org/wiki/APU_Triangle
var triangleTable = []byte{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
}

type triangle struct {
	seqPos              byte
	lc                  lengthCounter
	linearCounterCtrl   bool
	linearCounter       byte
	linearCounterPeriod byte
	linearCounterReload bool

	timer timer
}

func newTriangle() *triangle {
	return &triangle{}
}

func (t *triangle) output() byte {
	// > Write a period value of 0 or 1 to $400A/$400B, causing a very high frequency.
	// > Due to the averaging effect of the lowpass filter, the resulting value is halfway between 7 and 8.
	if t.timer.period < 2 {
		return 7
	}
	return triangleTable[t.seqPos]
}

func (t *triangle) tickTimer() {
	if t.timer.tick() {
		if t.lc.value > 0 && t.linearCounter > 0 {
			t.seqPos = (t.seqPos + 1) % 32
		}
	}
}

func (t *triangle) tickLinearCounter() {
	// > 1. If the linear counter reload flag is set, the linear counter is reloaded with the counter reload value, otherwise if the linear counter is non-zero, it is decremented.
	// > 2. If the control flag is clear, the linear counter reload flag is cleared.
	if t.linearCounterReload {
		t.linearCounter = t.linearCounterPeriod
	} else if t.linearCounter > 0 {
		t.linearCounter--
	}
	if !t.linearCounterCtrl {
		t.linearCounterReload = false
	}
}

func (t *triangle) tickLengthCounter() {
	t.lc.tick()
}
