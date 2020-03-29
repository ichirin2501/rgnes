package main

import (
	"fmt"
	"os"

	"github.com/ichirin2501/rgnes/nes"
)

func main() {

	f := os.Args[1]
	fmt.Println("f = ", f)

	c, err := nes.NewCassette(f)
	if err != nil {
		panic(err)
	}

	//fmt.Println(c)
	fmt.Printf("len(c.PRG) = 0x%04X\n", len(c.PRG))
	fmt.Printf("len(c.CHR) = 0x%04X\n", len(c.CHR))

	// cycle := nes.NewCPUCycle()
	// ram := nes.NewMemory(0x810)

	// cpuBus := nes.NewCPUBus(cycle, ram, c.ProgramROM)
	// irp := nes.NewInterrupt()

	// cpu := nes.NewCPU(cpuBus, cycle, irp)
	// cpu.Reset()
	// fmt.Println(cpu)

}
