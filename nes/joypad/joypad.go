package joypad

import "sync"

// https://www.nesdev.org/wiki/Controller_reading_code
const (
	ButtonA = (1 << iota)
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonUP
	ButtonDown
	ButtonLeft
	ButtonRight
)

type Joypad struct {
	Strobe       bool
	ButtonIndex  byte
	ButtonStatus byte
	mu           *sync.RWMutex
}

func New() *Joypad {
	return &Joypad{
		mu: &sync.RWMutex{},
	}
}

func (j *Joypad) Read() byte {
	if j.ButtonIndex > 7 {
		return 1
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	res := (byte(j.ButtonStatus) & (1 << j.ButtonIndex)) >> j.ButtonIndex
	if !j.Strobe && j.ButtonIndex <= 7 {
		j.ButtonIndex++
	}
	return res
}

func (j *Joypad) Write(v byte) {
	j.Strobe = (v & 1) == 1
	if j.Strobe {
		j.ButtonIndex = 0
	}
}

func (j *Joypad) SetButtonStatus(b byte) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.ButtonStatus = b
}
