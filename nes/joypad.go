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

type joypad struct {
	strobe       bool
	buttonIndex  byte
	buttonStatus byte
	mu           *sync.RWMutex // for ButtonStatus
}

func newJoypad() *joypad {
	return &joypad{
		mu: &sync.RWMutex{},
	}
}

func (j *joypad) read() byte {
	if j.buttonIndex > 7 {
		return 1
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	res := (byte(j.buttonStatus) & (1 << j.buttonIndex)) >> j.buttonIndex
	if !j.strobe && j.buttonIndex <= 7 {
		j.buttonIndex++
	}
	return res
}

// peek is used for debugging
func (j *joypad) peek() byte {
	if j.buttonIndex > 7 {
		return 1
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	return (byte(j.buttonStatus) & (1 << j.buttonIndex)) >> j.buttonIndex
}

func (j *joypad) write(v byte) {
	j.strobe = (v & 1) == 1
	if j.strobe {
		j.buttonIndex = 0
	}
}

func (j *joypad) setButtonStatus(b byte, pressed bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if pressed {
		j.buttonStatus |= b
	} else {
		j.buttonStatus &= ^b
	}
}
