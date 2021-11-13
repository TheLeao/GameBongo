package core

import "fmt"

type Ram struct {
	space  []int
	Length int
	Offset int
}

func NewRam(offset int, ln int) Ram {
	return Ram{
		space:  make([]int, ln),
		Offset: offset,
		Length: ln,
	}
}

func (r *Ram) Accepts(addr int) bool {
	return addr >= r.Offset && addr < r.Offset+r.Length
}

func (r *Ram) SetByte(addr int, value int) {
	r.space[addr-r.Offset] = value
}

func (r *Ram) GetByte(addr int) int {
	index := addr - r.Offset
	if index < 0 || index >= len(r.space) {
		panic(fmt.Sprintf("Ram: Index out of bounds. Address %d", addr))
	}
	return r.space[index]
}

type GbcRam struct {
	ram  [7 * 0x1000]int
	svbk int
}

//interface
func (g *GbcRam) Accepts(addr int) bool {
	return addr == 0xff70 || (addr >= 0xd000 && addr < 0xe000)
}

func (g *GbcRam) SetByte(addr int, value int) {
	if addr == 0xff70 {
		g.svbk = value
	} else {
		g.ram[g.translate(addr)] = value
	}
}

func (g *GbcRam) GetByte(addr int) int {
	if addr == 0xff70 {
		return g.svbk
	} else {
		return g.ram[g.translate(addr)]
	}
}

//

func (g *GbcRam) translate(addr int) int {
	ramBank := g.svbk & 0x7
	if ramBank == 0 {
		ramBank = 1
	}

	res := addr - 0xd000 + (ramBank-1)*0x1000
	if res < 0 || res >= len(g.ram) {
		panic("GBC Ram translate: illegal argument")
	}

	return res
}
