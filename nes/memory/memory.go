package memory

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

type memory []byte

func NewMemory(size int) Memory {
	buf := make([]byte, size)
	return memory(buf)
}

func (m memory) Read(addr uint16) byte {
	return m[addr]
}
func (m memory) Write(addr uint16, val byte) {
	m[addr] = val
}

func Read16(m MemoryReader, addr uint16) uint16 {
	l := m.Read(addr)
	h := m.Read(addr + 1)
	return (uint16(h) << 8) | uint16(l)
}
