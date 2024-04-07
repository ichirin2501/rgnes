package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: `H` の文字のテスト

func Test_PPU_MirrorVRAMAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ram  *ppuRAM
		addr uint16
		want uint16
	}{
		{
			"1",
			newPPURAM(MirroringHorizontal),
			0x2003,
			0x0003,
		},
		{
			"2",
			newPPURAM(MirroringHorizontal),
			0x2403,
			0x0003,
		},
		{
			"3",
			newPPURAM(MirroringHorizontal),
			0x2800,
			0x0400,
		},
		{
			"4",
			newPPURAM(MirroringHorizontal),
			0x2C00,
			0x0400,
		},
		{
			"5",
			newPPURAM(MirroringVertical),
			0x2000,
			0x0000,
		},
		{
			"6",
			newPPURAM(MirroringVertical),
			0x2801,
			0x0001,
		},
		{
			"7",
			newPPURAM(MirroringVertical),
			0x2400,
			0x0400,
		},
		{
			"8",
			newPPURAM(MirroringVertical),
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
				ppu.writeScroll(0x00)
				ppu.writeScroll(0x00)
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

// func Test_PPU_InternalRegisters(t *testing.T) {
// 	t.Parallel()
// 	tests := []struct {
// 		name         string
// 		ppu          *PPU
// 		instructions func(ppu *PPU)
// 		wantt        uint16
// 		wantv        uint16
// 		wantx        byte
// 		wantw        bool
// 	}{
// 		{
// 			"1",
// 			&PPU{},
// 			func(ppu *PPU) {
// 				ppu.writeController(0x00)
// 				ppu.readStatus()
// 				ppu.writeScroll(0x7D)
// 				ppu.writeScroll(0x5E)
// 				ppu.writePPUAddr(0x3D)
// 				ppu.writePPUAddr(0xF0)
// 			},
// 			0x3DF0,
// 			0x3DF0,
// 			0x5,
// 			false,
// 		},
// 		{
// 			"2",
// 			&PPU{},
// 			func(ppu *PPU) {
// 				ppu.writePPUAddr(0x04)
// 				ppu.writeScroll(0x3E)
// 				ppu.writeScroll(0x7D)
// 				ppu.writePPUAddr(0xEF)
// 			},
// 			0x64EF,
// 			0x64EF,
// 			0x5,
// 			false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			tt.instructions(tt.ppu)
// 			assert.Equal(t, tt.wantt, tt.ppu.t)
// 			assert.Equal(t, tt.wantv, tt.ppu.v)
// 			assert.Equal(t, tt.wantx, tt.ppu.x)
// 			assert.Equal(t, tt.wantw, tt.ppu.w)
// 		})
// 	}
// }

func Test_PeekWriteOnlyRegister(t *testing.T) {
	t.Parallel()
	ppu := &PPU{
		iobus: ppuDecayRegister{
			val: 0x30,
		},
	}
	// PPUCTRL
	got := ppu.PeekRegister(0x2000)
	want := ppu.ReadRegister(0x2000)
	assert.Equal(t, want, got)

	// PPUMASK
	got = ppu.PeekRegister(0x2001)
	want = ppu.ReadRegister(0x2001)
	assert.Equal(t, want, got)

	// PPUOAMADDR
	got = ppu.PeekRegister(0x2003)
	want = ppu.ReadRegister(0x2003)
	assert.Equal(t, want, got)

	// PPUSCROLL
	got = ppu.PeekRegister(0x2005)
	want = ppu.ReadRegister(0x2005)
	assert.Equal(t, want, got)

	// PPUADDR
	got = ppu.PeekRegister(0x2006)
	want = ppu.ReadRegister(0x2006)
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

			// OAMDATA
			peek := tt.ppu.PeekRegister(0x2004)
			got := tt.ppu.ReadRegister(0x2004)

			assert.Equal(t, peek, got)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOpenbus, tt.ppu.iobus.get(tt.ppu.Clock))
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
		wantOpenbus            byte
	}{
		{
			"1",
			&PPU{
				status: ppuStatusRegister(0x88),
				iobus:  ppuDecayRegister{val: 0x31},
				//    0x ---D DDDD
				// -: 0x 1000 1000
				// D: 0x 0011 0001
				//    0x 1001 0001
				w: true,
			},
			0x91,
			false,
			false,
			0x91,
		},
		{
			"2",
			&PPU{
				status:   ppuStatusRegister(0x00),
				iobus:    ppuDecayRegister{val: 0x01},
				Scanline: 241,
				Cycle:    0,
			},
			0x01,
			false,
			true,
			0x01,
		},
		{
			"3",
			&PPU{
				status:   ppuStatusRegister(0x00),
				iobus:    ppuDecayRegister{val: 0x01},
				Scanline: 241,
				Cycle:    1,
			},
			0x01,
			false,
			false,
			0x01,
		},
		{
			"4",
			&PPU{
				status:   ppuStatusRegister(0x00),
				iobus:    ppuDecayRegister{val: 0x01},
				Scanline: 241,
				Cycle:    2,
			},
			0x01,
			false,
			false,
			0x01,
		},
		{
			"5",
			&PPU{
				status:   ppuStatusRegister(0x00),
				iobus:    ppuDecayRegister{val: 0x01},
				Scanline: 241,
				Cycle:    3,
			},
			0x01,
			false,
			false,
			0x01,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			nmiLine := NMIInterruptLine(0)
			tt.ppu.nmiLine = &nmiLine

			peek := tt.ppu.PeekRegister(0x2002)
			got := tt.ppu.ReadRegister(0x2002)

			assert.Equal(t, peek, got)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantw, tt.ppu.w)
			assert.Equal(t, tt.wantSuppressVBlankFlag, tt.ppu.suppressVBlankFlag)
			assert.Equal(t, tt.wantOpenbus, tt.ppu.iobus.get(tt.ppu.Clock))
		})
	}
}
