package nes

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
	mu           *sync.RWMutex // for ButtonStatus
}

func NewJoypad() *Joypad {
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

// Peek is used for debugging
func (j *Joypad) Peek() byte {
	if j.ButtonIndex > 7 {
		return 1
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	return (byte(j.ButtonStatus) & (1 << j.ButtonIndex)) >> j.ButtonIndex
}

func (j *Joypad) Write(v byte) {
	j.Strobe = (v & 1) == 1
	if j.Strobe {
		j.ButtonIndex = 0
	}
}

func (j *Joypad) SetButtonStatus(b byte, pressed bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if pressed {
		j.ButtonStatus |= b
	} else {
		j.ButtonStatus &= ^b
	}
}
