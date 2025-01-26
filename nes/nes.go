package nes

import "time"

type NES struct {
	cpu    *cpu
	apu    *apu
	ppu    *ppu
	bus    *cpuBus
	joypad *joypad

	done chan struct{}
}

func New(mapper Mapper, renderer Renderer, player Player) *NES {
	irqLine := irqInterruptLine(0)
	nmiLine := nmiInterruptLine(0)
	m := mapper.MirroingType()
	dma := &dma{}

	ppu := newPPU(renderer, mapper, m, &nmiLine)
	joypad := newJoypad()
	apu := newAPU(&irqLine, player, dma)
	bus := newCPUBus(ppu, apu, mapper, joypad, dma)

	cpu := newCPU(bus, &nmiLine, &irqLine)

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
	n.cpu.powerUp()
	n.apu.powerUp()
	n.ppu.powerUp()
}

func (n *NES) Reset() {
	n.cpu.reset()
	n.apu.reset()
	n.ppu.reset()
}

func (n *NES) Step() {
	n.cpu.step()
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
				n.cpu.step()
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
	n.joypad.setButtonStatus(b, pressed)
}

func (n *NES) SetCPUPC(pc uint16) {
	n.cpu.PC = pc
}

func (n *NES) PeekMemory(addr uint16) byte {
	return n.bus.peek(addr)
}
