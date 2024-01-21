package nes

type DMA struct {
	oamDMAOccurred bool
	oamTargetAddr  uint16

	dcmDMAOccurred bool
	dcmTargetAddr  uint16
}

func (d *DMA) SignalOAMDMA(addr uint16) {
	d.oamDMAOccurred = true
	d.oamTargetAddr = addr
}

func (d *DMA) SignalDCMDMA(addr uint16) {
	d.dcmDMAOccurred = true
	d.dcmTargetAddr = addr
}
