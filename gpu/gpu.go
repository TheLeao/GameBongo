package gpu

import (
	"github.com/theleao/goingboy/gameboy"
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
	intrptr        gameboy.Interrupter
	gbc            bool
	memRegs        gameboy.MemoryRegisters
	dma            gameboy.Dma
	bgPalette      ColorPalette
	oamPalette     ColorPalette
	lcdEnableDelay int
	display        Display
	phase          GpuPhase
	hBlank         HBlankPhase
	oamSearch      OamSearch
	pixelTransfer  PixelTransfer
	vBlank         VBlankPhase
}

func NewGpu(display Display, intrptr gameboy.Interrupter, dma gameboy.Dma, oamRam gameboy.Ram, gbc bool) Gpu {
	InitializeTileAttributes()

	gpu := Gpu{}

	var mr []gameboy.MemRegisterType
	for _, r := range GpuRegisters() {
		mr = append(mr, gameboy.MemRegisterType{
			Addr: r,
		})
	}

	gpu.memRegs = gameboy.NewMemRegisters(mr...)
	gpu.Lcdc = Lcdc{}
	gpu.intrptr = gameboy.Interrupter{}

	gpu.gbc = gbc
	gpu.vRam0 = gameboy.Ram{
		Offset: 0x8000,
		Length: 0x2000,
	}
	
	if gbc {
		gpu.vRam1 = gameboy.Ram{
			Offset: 0x8000,
			Length: 0x2000,
		}
	}

	gpu.bgPalette = NewColorPallete(0xff68)
	gpu.oamPalette = NewColorPallete(0xff6a)
	gpu.oamPalette.FillWithFF()

	gpu.oamSearch = NewOamSearch(gpu.oamRam, gpu.Lcdc, gpu.memRegs)
	gpu.pixelTransfer = NewPixelTransfer(gpu.vRam0, gpu.vRam1, gpu.oamRam, gpu.Lcdc, gpu.memRegs, gpu.gbc, gpu.bgPalette, gpu.oamPalette)

	return gpu
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
	g.phase = NewHBlankPhase(250) //hBlankPhase.Start()
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

func (g *Gpu) Tick() int {
	if g.EnabledLcd {
		if g.lcdEnableDelay != -1 {
			g.lcdEnableDelay = g.lcdEnableDelay - 1
			if g.lcdEnableDelay == 0 {
				g.display.EnableLcd()
				g.EnabledLcd = true
			}
		}
	}
	if !g.EnabledLcd {
		return -1
	}

	oldMode := g.mode
	g.TicksInLine++
	if g.phase.Tick() {
		if g.TicksInLine == 4 && g.mode == VBLANK && g.memRegs.Get(LY) == 153 {
			g.memRegs.Put(LY, 0)
			g.requestLycEqualsLyInterrupt()
		}
	} else {
		switch oldMode {
		case OAMSEARCH:
			g.mode = PIXELTRANSFER
			g.phase = nil //pixelTransferPhase.start(oamSearchPhase.getSprites());
		case PIXELTRANSFER:
			g.mode = HBLANK
			g.phase = NewHBlankPhase(g.TicksInLine)
			g.requestLcdcInterrupt(3)
		case HBLANK:
			g.TicksInLine = 0
			rAddr, _ := GetGpuRegister(LY)
			if g.memRegs.PreIncrement(rAddr) == 144 {
				g.mode = VBLANK
				g.phase = &VBlankPhase{}
				g.intrptr.RequestInterrupt(VBLANK)
				g.requestLcdcInterrupt(4)
			} else {
				g.mode = OAMSEARCH
				g.oamSearch.Start()
				g.phase = &g.oamSearch
			}
			g.requestLcdcInterrupt(5)
			g.requestLycEqualsLyInterrupt()
		case VBLANK:
			g.TicksInLine = 0
			rAddr, _ := GetGpuRegister(LY)
			if g.memRegs.PreIncrement(rAddr) == 1 {
				g.mode = OAMSEARCH
				g.memRegs.Put(LY, 0)
				g.oamSearch.Start()
				g.phase = &g.oamSearch
				g.requestLcdcInterrupt(5)
			} else {
				g.phase = &VBlankPhase{}
			}
			g.requestLycEqualsLyInterrupt()
		}
	}

	if oldMode == g.mode {
		return -1
	} else {
		return g.mode
	}
}
