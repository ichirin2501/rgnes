package cpu

// func Test_AddressingMode(t *testing.T) {
// 	t.Parallel()
// 	tests := []struct {
// 		name            string
// 		op              *opcode
// 		cpu             *CPU
// 		m               []byte // MemoryReader
// 		wantAddr        uint16
// 		wantPageCrossed bool
// 	}{
// 		{
// 			"absolute",
// 			&opcode{Mode: absolute},
// 			&CPU{},
// 			[]byte{0x01, 0x20},
// 			0x2001,
// 			false,
// 		},
// 		{
// 			"zeroPage",
// 			&opcode{Mode: zeroPage},
// 			&CPU{},
// 			[]byte{0x01, 0x20},
// 			0x0001,
// 			false,
// 		},
// 		{
// 			"zeroPageX/1",
// 			&opcode{Mode: zeroPageX},
// 			&CPU{X: 0xFF},
// 			[]byte{0x01, 0x20},
// 			0x0000,
// 			false,
// 		},
// 		{
// 			"zeroPageX/2",
// 			&opcode{Mode: zeroPageX},
// 			&CPU{X: 0x10},
// 			[]byte{0x01, 0x20},
// 			0x0011,
// 			false,
// 		},
// 		{
// 			"zeroPageY/1",
// 			&opcode{Mode: zeroPageY},
// 			&CPU{Y: 0x20},
// 			[]byte{0x01, 0x20},
// 			0x0021,
// 			false,
// 		},
// 		{
// 			"relative/1",
// 			&opcode{Mode: relative},
// 			&CPU{PC: 0x5},
// 			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00},
// 			0x0007,
// 			false,
// 		},
// 		{
// 			"relative/2",
// 			&opcode{Mode: relative},
// 			&CPU{PC: 0x5},
// 			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00},
// 			0x0005,
// 			false,
// 		},
// 		{
// 			"absoluteX/1",
// 			&opcode{Mode: absoluteX},
// 			&CPU{PC: 0x0, X: 0x1},
// 			[]byte{0xFF, 0x00},
// 			0x0100,
// 			true,
// 		},
// 		{
// 			"absoluteX/2",
// 			&opcode{Mode: absoluteX},
// 			&CPU{PC: 0x0, X: 0x1},
// 			[]byte{0xFE, 0x01},
// 			0x01FF,
// 			false,
// 		},
// 		{
// 			"absoluteY/2",
// 			&opcode{Mode: absoluteY},
// 			&CPU{PC: 0x0, Y: 0x1},
// 			[]byte{0xFE, 0x02},
// 			0x02FF,
// 			false,
// 		},
// 		{
// 			"indirect",
// 			&opcode{Mode: indirect},
// 			&CPU{PC: 0x00},
// 			[]byte{0x03, 0x00, 0x00, 0x22, 0x33},
// 			0x3322,
// 			false,
// 		},
// 		// これの追加テストしたい
// 		{
// 			"indexedIndirect/1",
// 			&opcode{Mode: indexedIndirect},
// 			&CPU{PC: 0x01, X: 0xFF},
// 			[]byte{0x00, 0x04, 0x00, 0x40, 0x01},
// 			0x0140,
// 			false,
// 		},
// 		{
// 			"indirectIndexed/1",
// 			&opcode{Mode: indirectIndexed},
// 			&CPU{PC: 0x00, Y: 0x20},
// 			[]byte{0x02, 0x00, 0x01, 0x50, 0x00},
// 			0x5021,
// 			false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			bus := &Bus{CPURAM: tt.m}
// 			cpu := tt.cpu
// 			cpu.m = bus
// 			gotAddr, gotPageCrossed := cpu.fetchOperand(tt.op)
// 			assert.Equal(t, tt.wantAddr, gotAddr)
// 			assert.Equal(t, tt.wantPageCrossed, gotPageCrossed)
// 		})
// 	}
// }
