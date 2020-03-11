package nes

type AddrReader interface {
	Read(addr uint16) byte
}
type AddrWriter interface {
	Write(addr uint16, val byte)
}
type Memory interface {
	AddrReader
	AddrWriter
}
