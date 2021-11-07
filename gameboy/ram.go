package gameboy

import "fmt"

type Ram struct {
	space  []int
	length int
	offset int
}

func (r Ram) Accepts(addr int) bool {
	return addr >= r.offset && addr < r.offset+r.length
}

func (r Ram) SetByte(addr int, value int) {
	ptr := &r
	ptr.space[addr-r.offset] = value
}

func (r Ram) GetByte(addr int) int {
	index := addr - r.offset
	if index < 0 || index >= len(r.space) {
		panic(fmt.Sprintf("Ram: Index out of bounds. Address %d", addr))
	}
	return r.space[index]
}