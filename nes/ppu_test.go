package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: `H` の文字のテスト

type fakeCPU struct{}

func (f *fakeCPU) SetDelayNMI()    {}
func (f *fakeCPU) SetNMI(val bool) {}

type fakeMirroringType int

const (
	fakeMirroringVertical fakeMirroringType = iota
	fakeMirroringHorizontal
	fakeMirroringFourScreen
)

func (f *fakeMirroringType) IsVertical() bool {
	return *f == fakeMirroringVertical
}
func (f *fakeMirroringType) IsHorizontal() bool {
	return *f == fakeMirroringHorizontal
}
func (f *fakeMirroringType) IsFourScreen() bool {
	return *f == fakeMirroringFourScreen
}

func getFakeMirroringVertical() fakeMirroringType {
	m := fakeMirroringVertical
	return m
}
func getFakeMirroringHorizontal() fakeMirroringType {
	m := fakeMirroringHorizontal
	return m
}
func getFakeMirroringFourScreen() fakeMirroringType {
	m := fakeMirroringFourScreen
	return m
}

func Test_PPU_MirrorVRAMAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ram  *vram
		addr uint16
		want uint16
	}{
		{
			"1",
			newVRAM(MirroringHorizontal),
			// newVRAM(getFakeMirroringHorizontal()),
			0x2003,
			0x0003,
		},
		{
			"2",
			newVRAM(MirroringHorizontal),
			//newVRAM(getFakeMirroringHorizontal()),
			0x2403,
			0x0003,
		},
		{
			"3",
			newVRAM(MirroringHorizontal),
			//newVRAM(getFakeMirroringHorizontal()),
			0x2800,
			0x0400,
		},
		{
			"4",
			newVRAM(MirroringHorizontal),
			//newVRAM(getFakeMirroringHorizontal()),
			0x2C00,
			0x0400,
		},
		{
			"5",
			newVRAM(MirroringVertical),
			//newVRAM(getFakeMirroringVertical()),
			0x2000,
			0x0000,
		},
		{
			"6",
			newVRAM(MirroringVertical),
			//newVRAM(getFakeMirroringVertical()),
			0x2801,
			0x0001,
		},
		{
			"7",
			newVRAM(MirroringVertical),
			//newVRAM(getFakeMirroringVertical()),
			0x2400,
			0x0400,
		},
		{
			"8",
			newVRAM(MirroringVertical),
			//newVRAM(getFakeMirroringVertical()),
			0x2C01,
			0x0401,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.ram.mirrorAddr(tt.addr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_PPU_IncrementY(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ppu  *PPU
		want uint16
	}{
		{
			"1",
			&PPU{
				v: 0x77A2,
			},
			0xC02,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.ppu.incrementY()
			assert.Equal(t, tt.want, tt.ppu.v)
		})
	}
}

func Test_PPU_WriteScroll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		ppu          *PPU
		instructions func(ppu *PPU)
		wantt        uint16
		wantx        byte
		wantw        bool
	}{
		{
			"1",
			&PPU{
				t: 0x21c0,
			},
			func(ppu *PPU) {
				ppu.WriteScroll(0x00)
				ppu.WriteScroll(0x00)
			},
			0x0000,
			0x0,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.instructions(tt.ppu)
			assert.Equal(t, tt.wantt, tt.ppu.t)
			assert.Equal(t, tt.wantx, tt.ppu.x)
			assert.Equal(t, tt.wantw, tt.ppu.w)
		})
	}
}

func Test_PPU_InternalRegisters(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		ppu          *PPU
		instructions func(ppu *PPU)
		wantt        uint16
		wantv        uint16
		wantx        byte
		wantw        bool
	}{
		{
			"1",
			&PPU{},
			func(ppu *PPU) {
				ppu.WriteController(0x00)
				ppu.ReadStatus()
				ppu.WriteScroll(0x7D)
				ppu.WriteScroll(0x5E)
				ppu.WritePPUAddr(0x3D)
				ppu.WritePPUAddr(0xF0)
			},
			0x3DF0,
			0x3DF0,
			0x5,
			false,
		},
		{
			"2",
			&PPU{},
			func(ppu *PPU) {
				ppu.WritePPUAddr(0x04)
				ppu.WriteScroll(0x3E)
				ppu.WriteScroll(0x7D)
				ppu.WritePPUAddr(0xEF)
			},
			0x64EF,
			0x64EF,
			0x5,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.instructions(tt.ppu)
			assert.Equal(t, tt.wantt, tt.ppu.t)
			assert.Equal(t, tt.wantv, tt.ppu.v)
			assert.Equal(t, tt.wantx, tt.ppu.x)
			assert.Equal(t, tt.wantw, tt.ppu.w)
		})
	}
}

func Test_PeekWriteOnlyRegister(t *testing.T) {
	t.Parallel()
	ppu := &PPU{
		openbus: PPUDecayRegister{
			val: 0x30,
		},
	}
	got := ppu.PeekController()
	want := ppu.ReadController()
	assert.Equal(t, want, got)

	got = ppu.PeekMask()
	want = ppu.ReadMask()
	assert.Equal(t, want, got)

	got = ppu.PeekOAMAddr()
	want = ppu.ReadOAMAddr()
	assert.Equal(t, want, got)

	got = ppu.PeekScroll()
	want = ppu.ReadScroll()
	assert.Equal(t, want, got)

	got = ppu.PeekPPUAddr()
	want = ppu.ReadPPUAddr()
	assert.Equal(t, want, got)
}

// todo
// func Test_ReadPPUData(t *testing.T) {
// 	t.Parallel()
// 	tests := []struct {
// 		name        string
// 		ppu         *PPU
// 		want        byte
// 		wantv       uint16
// 		wantOpenbus byte
// 	}{
// 		{
// 			"1",
// 			&PPU{
// 				ctrl:    ControlRegister(0),
// 				buf:     0x10,
// 				openbus: 0x00,
// 				v:       0x0000,
// 			},
// 			0x10,
// 			0x10,
// 			0x10,
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			peek := tt.ppu.PeekPPUData()
// 			got := tt.ppu.ReadPPUData()

// 			assert.Equal(t, peek, got)
// 			assert.Equal(t, tt.want, got)
// 			assert.Equal(t, tt.wantv, tt.ppu.v)
// 			assert.Equal(t, tt.wantOpenbus, tt.ppu.openbus)
// 		})
// 	}

// }

func Test_ReadOAMData(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		ppu         *PPU
		want        byte
		wantOpenbus byte
	}{
		{
			"1",
			&PPU{
				primaryOAM: [256]byte{0x10},
				oamAddr:    0x00,
			},
			0x10,
			0x10,
		},
		{
			"2",
			&PPU{
				primaryOAM: [256]byte{0x10, 0, 0, 0, 0x20},
				oamAddr:    0x04,
			},
			0x20,
			0x20,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			peek := tt.ppu.PeekOAMData()
			got := tt.ppu.ReadOAMData()

			assert.Equal(t, peek, got)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOpenbus, tt.ppu.getOpenBus())
		})
	}

}

func Test_ReadStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		ppu                    *PPU
		want                   byte
		wantw                  bool
		wantSuppressVBlankFlag bool
		wantNMIDelay           int
		wantOpenbus            byte
	}{
		{
			"1",
			&PPU{
				status:  PPUStatusRegister(0x88),
				openbus: PPUDecayRegister{val: 0x31},
				cpu:     &Interrupter{},
				w:       true,
			},
			0x99,
			false,
			false,
			0,
			0x31,
		},
		{
			"2",
			&PPU{
				status:   PPUStatusRegister(0x00),
				openbus:  PPUDecayRegister{val: 0x01},
				cpu:      &Interrupter{},
				Scanline: 241,
				Cycle:    0,
				nmiDelay: 3,
			},
			0x01,
			false,
			true,
			3,
			0x01,
		},
		{
			"3",
			&PPU{
				status:   PPUStatusRegister(0x00),
				openbus:  PPUDecayRegister{val: 0x01},
				cpu:      &Interrupter{},
				Scanline: 241,
				Cycle:    1,
				nmiDelay: 3,
			},
			0x01,
			false,
			false,
			0,
			0x01,
		},
		{
			"4",
			&PPU{
				status:   PPUStatusRegister(0x00),
				openbus:  PPUDecayRegister{val: 0x01},
				cpu:      &Interrupter{},
				Scanline: 241,
				Cycle:    2,
				nmiDelay: 3,
			},
			0x01,
			false,
			false,
			0,
			0x01,
		},
		{
			"5",
			&PPU{
				status:   PPUStatusRegister(0x00),
				openbus:  PPUDecayRegister{val: 0x01},
				cpu:      &Interrupter{},
				Scanline: 241,
				Cycle:    3,
				nmiDelay: 3,
			},
			0x01,
			false,
			false,
			3,
			0x01,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			peek := tt.ppu.PeekStatus()
			got := tt.ppu.ReadStatus()

			assert.Equal(t, peek, got)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantw, tt.ppu.w)
			assert.Equal(t, tt.wantSuppressVBlankFlag, tt.ppu.suppressVBlankFlag)
			assert.Equal(t, tt.wantNMIDelay, tt.ppu.nmiDelay)
			assert.Equal(t, tt.wantOpenbus, tt.ppu.getOpenBus())
		})
	}
}
