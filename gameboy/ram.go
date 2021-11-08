package gameboy

import "fmt"

type Ram struct {
	space  []int
	Length int
	Offset int
}

func (r Ram) Accepts(addr int) bool {
	return addr >= r.Offset && addr < r.Offset+r.Length
}

func (r Ram) SetByte(addr int, value int) {
	ptr := &r
	ptr.space[addr-r.Offset] = value
}

func (r Ram) GetByte(addr int) int {
	index := addr - r.Offset
	if index < 0 || index >= len(r.space) {
		panic(fmt.Sprintf("Ram: Index out of bounds. Address %d", addr))
	}
	return r.space[index]
}