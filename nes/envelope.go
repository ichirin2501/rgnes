package nes

type envelope struct {
	constantVolume    bool
	loop              bool
	start             bool
	divider           divider
	decayLevelCounter byte
}

func (e *envelope) output() byte {
	// https://www.nesdev.org/wiki/APU_Envelope
	// > The envelope unit's volume output depends on the constant volume flag:
	// > if set, the envelope parameter directly sets the volume, otherwise the decay level is the current volume.
	if e.constantVolume {
		// > bits 3-0	---- VVVV	Used as the volume in constant volume (C set) mode.
		// > Also used as the reload value for the envelope's divider (the period becomes V + 1 quarter frames).
		return byte(e.divider.period)
	} else {
		return e.decayLevelCounter
	}
}

func (e *envelope) tick() {
	// https://www.nesdev.org/wiki/APU_Envelope
	// > When clocked by the frame counter, one of two actions occurs:
	// > if the start flag is clear, the divider is clocked,
	// > otherwise the start flag is cleared, the decay level counter is loaded with 15, and the divider's period is immediately reloaded.
	if e.start {
		e.start = false
		e.decayLevelCounter = 15
		e.divider.reload()
	} else if e.divider.tick() {
		// > When the divider is clocked while at 0, it is loaded with V and clocks the decay level counter.
		// > Then one of two actions occurs: If the counter is non-zero, it is decremented, otherwise if the loop flag is set, the decay level counter is loaded with 15.
		if e.decayLevelCounter > 0 {
			e.decayLevelCounter--
		} else if e.loop {
			e.decayLevelCounter = 15
		}
	}
}
