package nes

type CPUBus struct {
	noCopy noCopy
}

func NewCPUBus() *CPUBus {
	return &CPUBus{}
}

func (bus *CPUBus) Read(addr uint16) byte {
	return 0
}

func (bus *CPUBus) Write(addr uint16, val byte) {

}
