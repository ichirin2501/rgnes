package nes

type CPUBus struct {
	cycle  *CPUCycle
	noCopy noCopy
}

func NewCPUBus(cycle *CPUCycle) *CPUBus {
	return &CPUBus{
		cycle: cycle,
	}
}

func (bus *CPUBus) Read(addr uint16) byte {
	return 0
}

func (bus *CPUBus) Write(addr uint16, val byte) {

}
