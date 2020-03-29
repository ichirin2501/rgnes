package cassette

import (
	"fmt"

	"github.com/ichirin2501/rgnes/nes/memory"
)

type Mapper interface {
	memory.Memory
	Reset()
}

func NewMapper(c *Cassette) Mapper {
	switch c.Mapper {
	case 0:
		return &Mapper0{c}
	}
	panic(fmt.Sprintf("Unsupported mapper: %0x", c.Mapper))
}
