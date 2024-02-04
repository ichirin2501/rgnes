package nes

type NES struct {
	cpu    *CPU
	apu    *APU
	ppu    *PPU
	bus    *Bus
	joypad *Joypad
}

func New(mapper Mapper, renderer Renderer, player Player) *NES {
	irp := &interruptLines{}
	m := mapper.MirroingType()
	dma := &DMA{}

	ppu := NewPPU(renderer, mapper, m, irp)
	joypad := NewJoypad()
	apu := NewAPU(irp, player, dma)
	bus := NewBus(ppu, apu, mapper, joypad, dma)

	cpu := NewCPU(bus, irp)

	return &NES{
		cpu:    cpu,
		apu:    apu,
		ppu:    ppu,
		bus:    bus,
		joypad: joypad,
	}
}

func (n *NES) PowerUp() {
	n.cpu.PowerUp()
	n.apu.PowerUp()
}

func (n *NES) Reset() {
	n.cpu.Reset()
	n.apu.Reset()
}

func (n *NES) Step() {
	n.cpu.Step()
}

func (n *NES) SetButtonStatus(b byte, pressed bool) {
	n.joypad.SetButtonStatus(b, pressed)
}

func (n *NES) SetCPUPC(pc uint16) {
	n.cpu.PC = pc
}

func (n *NES) PeekMemory(addr uint16) byte {
	return n.bus.Peek(addr)
}
