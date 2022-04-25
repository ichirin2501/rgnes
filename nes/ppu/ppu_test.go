package ppu

import (
	"testing"

	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/stretchr/testify/assert"
)

// TODO: `H` の文字のテスト

func Test_PPU_MirrorVRAMAddr(t *testing.T) {
	tests := []struct {
		name string
		ppu  *PPU
		addr uint16
		want uint16
	}{
		{
			"1",
			&PPU{mirroring: cassette.MirroringHorizontal},
			0x2003,
			0x0003,
		},
		{
			"2",
			&PPU{mirroring: cassette.MirroringHorizontal},
			0x2403,
			0x0003,
		},
		{
			"3",
			&PPU{mirroring: cassette.MirroringHorizontal},
			0x2800,
			0x0400,
		},
		{
			"4",
			&PPU{mirroring: cassette.MirroringHorizontal},
			0x2C00,
			0x0400,
		},
		{
			"5",
			&PPU{mirroring: cassette.MirroringVertical},
			0x2000,
			0x0000,
		},
		{
			"6",
			&PPU{mirroring: cassette.MirroringVertical},
			0x2801,
			0x0001,
		},
		{
			"7",
			&PPU{mirroring: cassette.MirroringVertical},
			0x2400,
			0x0400,
		},
		{
			"8",
			&PPU{mirroring: cassette.MirroringVertical},
			0x2C01,
			0x0401,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ppu.mirrorVRAMAddr(tt.addr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_PPU_IncrementY(t *testing.T) {
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
			tt.ppu.incrementY()
			assert.Equal(t, tt.want, tt.ppu.v)
		})
	}
}

func Test_PPU_WriteScroll(t *testing.T) {
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
			tt.instructions(tt.ppu)
			assert.Equal(t, tt.wantt, tt.ppu.t)
			assert.Equal(t, tt.wantx, tt.ppu.x)
			assert.Equal(t, tt.wantw, tt.ppu.w)
		})
	}
}

func Test_PPU_InternalRegisters(t *testing.T) {
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
			tt.instructions(tt.ppu)
			assert.Equal(t, tt.wantt, tt.ppu.t)
			assert.Equal(t, tt.wantv, tt.ppu.v)
			assert.Equal(t, tt.wantx, tt.ppu.x)
			assert.Equal(t, tt.wantw, tt.ppu.w)
		})
	}
}
