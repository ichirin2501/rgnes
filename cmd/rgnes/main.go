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

	fmt.Println(c)

	for i := 0; i < 10; i++ {
		fmt.Printf("%0x\n", c.ProgramROM.Read(uint16(i)))
	}

	cpuBus := nes.NewCPUBus()
	irp := nes.NewInterrupt()

	cpu := nes.NewCPU(cpuBus, irp)
	fmt.Println(cpu)
}
