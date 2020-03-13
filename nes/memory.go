package nes

type MemoryReader interface {
	Read(addr uint16) byte
}
type MemoryWriter interface {
	Write(addr uint16, val byte)
}
type Memory interface {
	MemoryReader
	MemoryWriter
}
