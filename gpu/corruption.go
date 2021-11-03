package gpu

import (
	"github.com/theleao/goingboy/gameboy"
)

const ( //corruption types
	INC_DEC = iota
	POP_1
	POP_2
	PUSH_1
	PUSH_2
	LD_HL
)

func copyValuesSpriteCorruption(addr *gameboy.AddressSpace, from int, to int, len int) {
	for i := len - 1; i >= 0; i-- {
		b := (*addr).GetByte(0xfe00+from+i) % 0xff
		(*addr).SetByte(0xfe00+to+i, b)
	}
}

func CorruptOam(addr *gameboy.AddressSpace, corruptionType int, ticksInLine int) {
	var cpuCycle int = (ticksInLine+1)/4 + 1

	switch corruptionType {
	case INC_DEC:
		if cpuCycle >= 2 {
			copyValuesSpriteCorruption(addr, (cpuCycle-2)*8+2, (cpuCycle-1)*8+2, 6)
		}
	case POP_1:
	case LD_HL:
		if cpuCycle >= 4 {
			copyValuesSpriteCorruption(addr, (cpuCycle-3)*8+2, (cpuCycle-4)*8+2, 8)
			copyValuesSpriteCorruption(addr, (cpuCycle-3)*8+8, (cpuCycle-4)*8+0, 2)
			copyValuesSpriteCorruption(addr, (cpuCycle-4)*8+2, (cpuCycle-2)*8+2, 6)
		}
	case POP_2:
		if cpuCycle >= 5 {
			copyValuesSpriteCorruption(addr, (cpuCycle-5)*8+0, (cpuCycle-2)*8+0, 8)
		}
	case PUSH_1:
		if cpuCycle >= 4 {
			copyValuesSpriteCorruption(addr, (cpuCycle-4)*8+2, (cpuCycle-3)*8+2, 8)
			copyValuesSpriteCorruption(addr, (cpuCycle-3)*8+2, (cpuCycle-1)*8+2, 6)
		}
	case PUSH_2:
		if cpuCycle >= 5 {
			copyValuesSpriteCorruption(addr, (cpuCycle-4)*8+2, (cpuCycle-3)*8+2, 8)
		}
	default:
		panic(nil) //corruption error?
	}

}
