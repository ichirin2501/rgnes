package nes

import "time"

type NES struct {
	cpu    *CPU
	apu    *APU
	ppu    *PPU
	bus    *CPUBus
	joypad *Joypad

	done chan struct{}
}

func New(mapper Mapper, renderer Renderer, player Player) *NES {
	irqLine := defaultInterruptLineState
	nmiLine := defaultInterruptLineState
	m := mapper.MirroingType()
	dma := &DMA{}

	ppu := NewPPU(renderer, mapper, m, &nmiLine)
	joypad := NewJoypad()
	apu := NewAPU(&irqLine, player, dma)
	bus := NewCPUBus(ppu, apu, mapper, joypad, dma)

	cpu := NewCPU(bus, &nmiLine, &irqLine)

	return &NES{
		cpu:    cpu,
		apu:    apu,
		ppu:    ppu,
		bus:    bus,
		joypad: joypad,

		done: make(chan struct{}),
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

func (n *NES) Run() {
	beforeTime := time.Now()
	steps := float64(0)
	for {
		select {
		case <-n.done:
			return
		default:
			now := time.Now()
			// du / (1sec/CPUClockFrequency)
			du := float64(now.Sub(beforeTime)*CPUClockFrequency) / float64(time.Second)
			if steps+du >= 1.0 {
				beforeClock := n.cpu.bus.realClock()
				n.cpu.Step()
				afterClock := n.cpu.bus.realClock()
				steps = steps + du - float64(afterClock-beforeClock)
				beforeTime = now
			}
		}
	}
}

func (n *NES) Close() {
	close(n.done)
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
