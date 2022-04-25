package nes

type ButtonType byte

// https://www.nesdev.org/wiki/Controller_reading_code
const (
	ButtonA ButtonType = (1 << iota)
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
	ButtonStatus ButtonType
}

func NewJoypad() *Joypad {
	return &Joypad{}
}

func (j *Joypad) Read() byte {
	if j.ButtonIndex > 7 {
		return 1
	}
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

func (j *Joypad) SetButtonPressedStatus(b ButtonType, pressed bool) {
	if pressed {
		j.ButtonStatus |= b
	} else {
		j.ButtonStatus &= ^b
	}
}
