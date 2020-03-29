package cpu

import (
	"path/filepath"
	"testing"

	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/bus"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/interrupt"
	"github.com/ichirin2501/rgnes/nes/memory"
	"github.com/ichirin2501/rgnes/nes/ppu"
)

func TestCPU(t *testing.T) {
	path := filepath.Join("testdata", "nestest.nes")
	c, err := cassette.NewCassette(path)
	if err != nil {
		t.Fatal(err)
	}
	mapper := cassette.NewMapper(c)
	cycle := NewCPUCycle()
	ram := memory.NewMemory(0x81000)
	ppu := ppu.NewPPU()
	apu := apu.NewAPU()
	cpuBus := bus.NewCPUBus(ram, ppu, apu, mapper)
	irp := interrupt.NewInterrupt()
	cpu := NewCPU(cpuBus, cycle, irp)

	cpu.r.PC = 0xC000
	for i := 0; i < 8991; i++ {
		cpu.Step()
	}
}
