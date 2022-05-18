package ppu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: `H` の文字のテスト

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

func getFakeMirroringVertical() *fakeMirroringType {
	m := fakeMirroringVertical
	return &m
}
func getFakeMirroringHorizontal() *fakeMirroringType {
	m := fakeMirroringHorizontal
	return &m
}
func getFakeMirroringFourScreen() *fakeMirroringType {
	m := fakeMirroringFourScreen
	return &m
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
			newVRAM(getFakeMirroringHorizontal()),
			0x2003,
			0x0003,
		},
		{
			"2",
			newVRAM(getFakeMirroringHorizontal()),
			0x2403,
			0x0003,
		},
		{
			"3",
			newVRAM(getFakeMirroringHorizontal()),
			0x2800,
			0x0400,
		},
		{
			"4",
			newVRAM(getFakeMirroringHorizontal()),
			0x2C00,
			0x0400,
		},
		{
			"5",
			newVRAM(getFakeMirroringVertical()),
			0x2000,
			0x0000,
		},
		{
			"6",
			newVRAM(getFakeMirroringVertical()),
			0x2801,
			0x0001,
		},
		{
			"7",
			newVRAM(getFakeMirroringVertical()),
			0x2400,
			0x0400,
		},
		{
			"8",
			newVRAM(getFakeMirroringVertical()),
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
