package cpu

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/bus"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/interrupt"
	"github.com/ichirin2501/rgnes/nes/memory"
	"github.com/ichirin2501/rgnes/nes/ppu"
)

func TestCPU(t *testing.T) {
	os.Setenv("DEBUG", "1")
	buf := bytes.NewBuffer(make([]byte, 0))
	debugWriter = buf
	defer func() {
		os.Setenv("DEBUG", "0")
		debugWriter = os.Stdout
	}()

	c, err := cassette.NewCassette("testdata/nestest.nes")
	if err != nil {
		t.Fatal(err)
	}
	mapper := cassette.NewMapper(c)
	cycle := NewCPUCycle()
	ram := memory.NewMemory(0x8100)

	ppuRam := memory.NewMemory(0x2000)
	ppuBus := bus.NewPPUBus(ppuRam, mapper)
	ppu := ppu.NewPPU(ppuBus)

	apu := apu.NewAPU()
	cpuBus := bus.NewCPUBus(ram, ppu, apu, mapper)
	irp := interrupt.NewInterrupt()
	cpu := NewCPU(cpuBus, cycle, irp)

	cpu.r.PC = 0xC000
	for i := 0; i < 8991; i++ {
		cpu.Step()
	}

	f, err := os.Open("testdata/nestest-formatted2.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	expectedReader := bufio.NewReader(f)

	for {
		got, err1 := buf.ReadString('\n')
		want, err2 := expectedReader.ReadString('\n')

		if err1 == io.EOF && err2 == io.EOF {
			break
		}
		if err1 != nil || err2 != nil {
			t.Fatal("don't eq")
		}
		assert.Equal(t, want, got)
	}
}
