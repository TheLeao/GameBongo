package gpu

import (
	"github.com/theleao/gamebongo/cpu"
	"github.com/theleao/gamebongo/gameboy"
	"github.com/theleao/gamebongo/memory"
)

const ( //GPU Mode
	HBLANK = iota
	VBLANK
	OAMSEARCH
	PIXELTRANSFER
)

type Gpu struct {
	Lcdc        Lcdc
	mode        int
	TicksInLine int
	vRam0       gameboy.AddressSpace
	vRam1       gameboy.AddressSpace
	oamRam      gameboy.AddressSpace
	intrptr     cpu.Interrupter
	gbc         bool
	memRegs     memory.MemRegisters
	dma         memory.Dma
	bgPalette   ColorPalette
	oamPalette  ColorPalette
}

//Implementing interface
func (g *Gpu) Accepts(addr int) bool {
	return g.GetAddressSpace(addr) != nil
}

func (g *Gpu) SetByte(addr int, value int) {
	a, _ := GetGpuRegister(STAT)
	if addr == a {
		set
	}
}

func (g *Gpu) GetByte(addr int) int {
}

func (g *Gpu) GetAddressSpace(addr int) gameboy.AddressSpace {
	if g.vRam0.Accepts(addr) {
		return g.getVideoRam()
	} else if g.oamRam.Accepts(addr) && !g.dma.IsOamBlocked() {
		return g.oamRam
	} else if g.Lcdc.Accepts(addr) {
		return &g.Lcdc
	} else if g.memRegs.Accepts(addr) {
		return &g.memRegs
	} else if g.gbc && g.bgPalette.Accepts(addr) {
		return &g.bgPalette
	} else if g.gbc && g.oamPalette.Accepts(addr) {
		return &g.oamPalette
	} else {
		return nil
	}
}

func (g *Gpu) getVideoRam() gameboy.AddressSpace {
	gpuRegAddr, _ := GetGpuRegister(VBK)
	if g.gbc && (gpuRegAddr&1) == 1 {
		return g.vRam1
	} else {
		return g.vRam0
	}
}

func (g *Gpu) setStat(value int) {
	a, _ := GetGpuRegister(STAT)
	g.memRegs.Put(a, value&0b11111000)
}

func (g *Gpu) getStat() int {
	statAddr, _ := GetGpuRegister(STAT)
	lycAddr, _ := GetGpuRegister(LYC)
	lyAddr, _ := GetGpuRegister(LY)
	l := 0
	if g.memRegs.Get(lycAddr) == g.memRegs.Get(lyAddr) {
		l = 1 << 2
	}
	return g.memRegs.Get(statAddr) | g.mode | l | 0x80
}
