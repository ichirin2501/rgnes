package nes

// https://www.nesdev.org/wiki/APU_Length_Counter
var lengthTable = []byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

type lengthCounter struct {
	enabled bool
	halt    bool
	value   byte
}

func (lc *lengthCounter) setEnabled(v bool) {
	if v {
		lc.enabled = true
	} else {
		lc.enabled = false
		lc.value = 0
	}
}

func (lc *lengthCounter) load(v byte) {
	// > If the enabled flag is set, the length counter is loaded with entry L of the length table
	if lc.enabled {
		lc.value = lengthTable[v]
	}
}

func (lc *lengthCounter) tick() {
	if lc.value > 0 && !lc.halt {
		lc.value--
	}
}
