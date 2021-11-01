package gpu

import (
	"github.com/theleao/gamebongo/cpu"
	"github.com/theleao/gamebongo/gameboy"
)

const ( //GPU Mode
	HBLANK = iota
	VBLANK
	OAMSEARCH
	PIXELTRANSFER
)

type Gpu struct {
	Lcdc        Lcdc
	Mode        int
	TicksInLine int
	vRam0 gameboy.AddressSpace
	vRam1 gameboy.AddressSpace
	oamRam gameboy.AddressSpace
	intrptr cpu.Interrupter
	
}

//Implementing interface
func (g *Gpu) Accepts(addr int) bool {
	return true
}

func (g *Gpu) SetByte(addr int, value int) {
}

func (g *Gpu) GetByte(addr int) int {
}

func (g *Gpu) GetAddressSpace() gameboy.AddressSpace {

}

type Lcdc struct {
	Enabled bool
}