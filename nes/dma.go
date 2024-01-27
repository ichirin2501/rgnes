package nes

type DMA struct {
	oamDMAOccurred bool
	oamTargetAddr  uint16

	dcmDMAOccurred bool
	dcmTargetAddr  uint16
	dcmDelay       byte
}

func (d *DMA) TriggerOnOAM(addr uint16) {
	d.oamDMAOccurred = true
	d.oamTargetAddr = addr
}

func (d *DMA) TriggerOnDCMLoad(addr uint16) {
	d.dcmTargetAddr = addr
	d.dcmDelay = 3
}

func (d *DMA) TriggerOnDCMReload(addr uint16) {
	d.dcmDMAOccurred = true
	d.dcmTargetAddr = addr
}
