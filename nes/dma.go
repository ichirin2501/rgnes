package nes

type oamDMAState int

const (
	oamDMANoneState oamDMAState = iota
	oamDMAHaltState
	oamDMAAlignmentState
	oamDMAReadState
	oamDMAWriteState
)

type dmcDMAState int

const (
	dmcDMANoneState dmcDMAState = iota
	dmcDMAHaltState
	dmcDMADummyState
	dmcDMAAlignmentState
	dmcDMARunState
)

type dma struct {
	oamTargetAddr uint16
	dmcTargetAddr uint16
	dmcDelay      byte

	oamState     oamDMAState
	oamSaveState oamDMAState
	oamTempByte  byte
	oamCount     uint16
	dmcState     dmcDMAState
}

func (d *dma) triggerOnOAM(addr uint16) {
	d.oamState = oamDMAHaltState
	d.oamTargetAddr = addr
}

func (d *dma) triggerOnDMCLoad(addr uint16) {
	d.dmcTargetAddr = addr
	d.dmcDelay = 3
}

func (d *dma) triggerOnDMCReload(addr uint16) {
	d.dmcState = dmcDMAHaltState
	d.dmcTargetAddr = addr
}
