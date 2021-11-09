package gpu

import "github.com/theleao/goingboy/core"

type GpuPhase interface {
	Tick() bool
}

//HBLANK
type HBlankPhase struct {
	ticks int
}

func NewHBlankPhase(ticks int) GpuPhase {
	return &HBlankPhase{
		ticks: ticks,
	}
}

func (h *HBlankPhase) Tick() bool {
	h.ticks++
	return h.ticks < 456
}

//VBLANK
type VBlankPhase struct {
	ticks int
}

func NewVBlankPhase() VBlankPhase {
	return VBlankPhase{}
}

func (v *VBlankPhase) Tick() bool {
	v.ticks++
	return v.ticks < 456
}

//OAM SEARCH
type OamSearch struct {
	state          rune
	OemRam         core.AddressSpace
	regs           core.MemoryRegisters
	Lcdc           Lcdc
	spritePosIndex int
	spriteX        int
	spriteY        int
	i              int
	Sprites        [10]SpritePosition
}

type SpritePosition struct {
	X       int
	Y       int
	Address int
}

func NewSpritePosition(x int, y int, addr int) SpritePosition {
	return SpritePosition{
		X:       x,
		Y:       y,
		Address: addr,
	}
}

func NewOamSearch(oemRam core.AddressSpace, lcdc Lcdc, reg core.MemoryRegisters) OamSearch {
	return OamSearch{
		OemRam: oemRam,
		regs: reg,
		Lcdc: lcdc,
	}
}

func (o OamSearch) Start() {
	ptr := &o
	ptr.spritePosIndex = 0
	ptr.state = 'Y'
	ptr.spriteX = 0
	ptr.spriteY = 0
	ptr.i = 0

	for j := 0; j < len(o.Sprites); j++ {
		ptr.Sprites[j] = SpritePosition{}
	}
}

func (o OamSearch) Tick() bool {
	ptr := &o
	spriteAddr := 0xfe00 + 4 * ptr.i

	switch ptr.state {
	case 'Y':
		ptr.spriteY = ptr.OemRam.GetByte(spriteAddr)
		ptr.state = 'X'
	case 'X':
		ptr.spriteX = ptr.OemRam.GetByte(spriteAddr + 1)

		if (ptr.spritePosIndex < len(ptr.Sprites)) && between(ptr.spriteY, ptr.regs.Get(LY) + 16, ptr.spriteY + ptr.Lcdc.GetSpriteHeight()) {
			ptr.Sprites[ptr.spritePosIndex] = NewSpritePosition(ptr.spriteX, ptr.spriteY, spriteAddr)
			ptr.spritePosIndex++
		}

		ptr.i++
		ptr.state = 'Y'
	}

	return ptr.i < 40
}

func between(from int, num int, to int) bool {
	return from <= num && num < to
}

//PIXEL TRANSFER
type PixelTransfer struct {
	fifo PixelQueue
	fetcher Fetcher
	lcdc Lcdc
	memRegs core.MemoryRegisters
	gbc bool
	sprites []SpritePosition
	droppedPx int
	x int
	window bool
}

func NewPixelTransfer(vram0 core.AddressSpace, vram1 core.AddressSpace, oemRam core.AddressSpace, lcdc Lcdc, regs core.MemoryRegisters, 
	gbc bool, bgPalette ColorPalette, oamPalette ColorPalette, display Display) PixelTransfer {
		var pq PixelQueue
		if gbc {
			pq = ColorPixelQueue{
				lcdc:       lcdc,
				display:    display,
				bgPalette:  bgPalette,
				oamPalette: oamPalette,
			}
		} else {
			pq = DmgPixelQueue{
				display: display,
				regs: regs,
			}
		}

		f := NewFetcher(pq, vram0, vram1, oemRam, lcdc, regs, gbc)

		return PixelTransfer{
			fifo: pq,
			fetcher: f,			
		}
}

func (p *PixelTransfer) Start(sprites []SpritePosition) {
	p.sprites = sprites
	p.droppedPx = 0
	p.x = 0
	p.window = false

	p.fetcher.Init()
	if p.gbc || p.lcdc.IsBgAndWindowDisplay() {
		p.fetchBackground()
	} else {
		p.fetcher.fetchDisabled = true
	}
}

func (p *PixelTransfer) fetchBackground() {
	bgX := p.memRegs.Get(SCX) / 0x08
	bgY := (p.memRegs.Get(SCY) + p.memRegs.Get(LY)) % 0x100

	p.fetcher.Fetch(p.lcdc.GetBgTileMapDisplay() + ((bgY/0x08) * 0x20), p.lcdc.GetBgWindowTileData(),
	bgX, p.lcdc.IsBgWindowTileDataSigned(), bgY % 0x08)
}

func (p *PixelTransfer) fetchWindow() {
	winX := (p.x - p.memRegs.Get(WX) + 7) / 0x08
	winY := p.memRegs.Get(LY) - p.memRegs.Get(WY)

	p.fetcher.Fetch(p.lcdc.GetWindowTileMapDisplay() + ((winY/0x08) * 0x20), p.lcdc.GetBgWindowTileData(),
	winX, p.lcdc.IsBgWindowTileDataSigned(), winY % 0x08)
}

func (p PixelTransfer) Tick() bool {
	ptr := &p
	p.fetcher.Tick()

	if p.lcdc.IsBgAndWindowDisplay() || p.gbc {
		if p.fifo.Length() <= 8 {
			return true
		}

		scxAddr, _ := GetGpuRegister(SCX)
		if p.droppedPx < p.memRegs.Get(scxAddr) % 8 {
			p.fifo.DropPixel()
			ptr.droppedPx += 1
			return true
		}

		lyAddr, _ := GetGpuRegister(LY)
		wyAddr, _ := GetGpuRegister(WY)
		wxAddr, _ := GetGpuRegister(WX)
		if !p.window && p.lcdc.IsWindowDisplay() && p.memRegs.Get(lyAddr) >= p.memRegs.Get(wyAddr) && p.x == p.memRegs.Get(wxAddr) - 7 {
			ptr.window = true
			p.fetchWindow()
			return true
		}
	}

	if p.lcdc.IsObjDisplay() {
		if p.fetcher.SpriteInProgress() {
			return true
		}

		spriteAdd := false

		for i := 0; i < len(p.sprites); i++ {
			s := p.sprites[i]
			//checking "nil"
			if (s.Address == 0 && s.X == 0 && s.Y == 0) {
				continue
			}

			if p.x == 0 && s.X < 8 {
				if !spriteAdd {
					p.fetcher.AddSprite(s, 8 - s.X, i)
					spriteAdd = true
				}
				p.sprites[i] = SpritePosition{}
			} else if s.X - 8 == p.x {
				if !spriteAdd {
					p.fetcher.AddSprite(s, 0, i)
					spriteAdd = true
				}
				p.sprites[i] = SpritePosition{}
			}

			if spriteAdd {
				return true
			}
		}
	}

	p.fifo.PutPixelToScreen()
	p.x++
	if p.x == 160 {
		return false
	}
	return true
}