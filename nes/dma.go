package nes

type OAMDMAState int

const (
	OAMDMANoneState OAMDMAState = iota
	OAMDMAHaltState
	OAMDMAAlignmentState
	OAMDMAReadState
	OAMDMAWriteState
)

type DMCDMAState int

const (
	DMCDMANoneState DMCDMAState = iota
	DMCDMAHaltState
	DMCDMADummyState
	DMCDMAAlignmentState
	DMCDMARunState
)

type DMA struct {
	oamTargetAddr uint16
	dmcTargetAddr uint16
	dmcDelay      byte

	oamState     OAMDMAState
	oamSaveState OAMDMAState
	oamTempByte  byte
	oamCount     uint16
	dmcState     DMCDMAState
}

func (d *DMA) TriggerOnOAM(addr uint16) {
	d.oamState = OAMDMAHaltState
	d.oamTargetAddr = addr
}

func (d *DMA) TriggerOnDMCLoad(addr uint16) {
	d.dmcTargetAddr = addr
	d.dmcDelay = 3
}

func (d *DMA) TriggerOnDMCReload(addr uint16) {
	d.dmcState = DMCDMAHaltState
	d.dmcTargetAddr = addr
}
