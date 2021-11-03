package gpu

import (
	"github.com/theleao/goingboy/gameboy"
	"github.com/theleao/goingboy/interrupter"
	"github.com/theleao/goingboy/memory"
)

const ( //GPU Mode
	HBLANK = iota
	VBLANK
	OAMSEARCH
	PIXELTRANSFER
)

type Gpu struct {
	EnabledLcd     bool
	Lcdc           Lcdc
	TicksInLine    int
	mode           int
	vRam0          gameboy.AddressSpace
	vRam1          gameboy.AddressSpace
	oamRam         gameboy.AddressSpace
	intrptr        interrupter.Interrupter
	gbc            bool
	memRegs        memory.MemRegisters
	dma            memory.Dma
	bgPalette      ColorPalette
	oamPalette     ColorPalette
	lcdEnableDelay int
	display        Display
	phase          int
}

type GpuPhase interface {
	Tick() bool
}

//Implementing interface
func (g *Gpu) Accepts(addr int) bool {
	return g.GetAddressSpace(addr) != nil
}

func (g *Gpu) SetByte(addr int, value int) {
	a, _ := GetGpuRegister(STAT)
	if addr == a {
		g.setStat(a)
	} else {
		addrSpace := g.GetAddressSpace(addr)
		if addrSpace == &g.Lcdc {
			g.setLcdc(value)
		} else if addrSpace != nil {
			//Is this necessary? 2
			addrSpace.SetByte(addr, value)
		}
	}
}

func (g *Gpu) GetByte(addr int) int {
	a, _ := GetGpuRegister(STAT)
	if addr == a {
		return g.getStat()
	} else {
		addrSpace := g.GetAddressSpace(addr)
		if addrSpace == nil {
			return 0xff
		}
		v, _ := GetGpuRegister(VBK)
		if addr == v {
			if g.gbc {
				return 0xfe
			} else {
				return 0xff
			}
		} else {
			return addrSpace.GetByte(addr)
		}
	}
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

func (g *Gpu) setLcdc(value int) {
	g.Lcdc.SetValue(value)
	if (value & (1 << 7)) == 0 {
		g.disableLcd()
	} else {
		g.enableLcd()
	}
}

func (g *Gpu) disableLcd() {
	g.memRegs.Put(LY, 0)
	g.TicksInLine = 0
	g.phase = 250
	g.mode = HBLANK
	g.EnabledLcd = false
	g.lcdEnableDelay = -1
	g.display.EnableLcd()
}

func (g *Gpu) enableLcd() {
	g.lcdEnableDelay = 244
}

func (g *Gpu) requestLcdcInterrupt(statBit int) {
	if (g.memRegs.Get(STAT) & (1 << statBit)) != 0 {
		g.intrptr.RequestInterrupt(0x0048)
	}
}

func (g *Gpu) requestLycEqualsLyInterrupt() {
	if g.memRegs.Get(LYC) == g.memRegs.Get(LY) {
		g.requestLcdcInterrupt(6)
	}
}
