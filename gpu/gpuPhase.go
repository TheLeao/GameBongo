package gpu

import "github.com/theleao/goingboy/gameboy"

type GpuPhase interface {
	Tick() bool
}

//HBLANK
type HBlankPhase struct {
	ticks int
}

func NewHBlankPhase(ticks int) HBlankPhase {
	return HBlankPhase{
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
	OemRam         gameboy.AddressSpace
	regs           gameboy.MemoryRegisters
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

func NewOamSearch(oemRam gameboy.AddressSpace, lcdc Lcdc, reg gameboy.MemoryRegisters) OamSearch {
	return OamSearch{
		OemRam: oemRam,
		regs: reg,
		Lcdc: lcdc,
	}
}

func (o *OamSearch) Start() {
	o.spritePosIndex = 0
	o.state = 'Y'
	o.spriteX = 0
	o.spriteY = 0
	o.i = 0

	for j := 0; j < len(o.Sprites); j++ {
		o.Sprites[j] = SpritePosition{}
	}
}

func (o *OamSearch) Tick() bool {
	spriteAddr := 0xfe00 + 4 * o.i

	switch o.state {
	case 'Y':
		o.spriteY = o.OemRam.GetByte(spriteAddr)
		o.state = 'X'
	case 'X':
		o.spriteX = o.OemRam.GetByte(spriteAddr + 1)

		if (o.spritePosIndex < len(o.Sprites)) && between(o.spriteY, o.regs.Get(LY) + 16, o.spriteY + o.Lcdc.GetSpriteHeight()) {
			o.Sprites[o.spritePosIndex] = NewSpritePosition(o.spriteX, o.spriteY, spriteAddr)
			o.spritePosIndex++
		}

		o.i++
		o.state = 'Y'
	}

	return o.i < 40
}

func between(from int, num int, to int) bool {
	return from <= num && num < to
}

//PIXEL TRANSFER
type PixelTransfer struct {
	fifo PixelQueue
	fetcher Fetcher
	lcdc Lcdc
	memRegs gameboy.MemoryRegisters
	gbc bool
	sprites []SpritePosition
	droppedPx int
	x int
	window bool
}

func NewPixelTransfer(vram0 gameboy.AddressSpace, vram1 gameboy.AddressSpace, oemRam gameboy.AddressSpace, lcdc Lcdc, regs gameboy.MemoryRegisters, 
	gbc bool, bgPalette ColorPalette, oamPalette ColorPalette) PixelTransfer {
		var pq PixelQueue
		if gbc {
			pq = nil //new ColorPixelFifo(lcdc, display, bgPalette, oamPalette);
		} else {
			pq = nil //new DmgPixelFifo(display, lcdc, r);
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
		p.FetchBackground()
	} else {
		p.fetcher.fetchDisabled = true
	}
}

func (p *PixelTransfer) FetchBackground() {
	bgX := p.memRegs.Get(SCX) / 0x08
	bgY := (p.memRegs.Get(SCY) + p.memRegs.Get(LY)) % 0x100

	p.fetcher.Fetch(p.lcdc.GetBgTileMapDisplay() + ((bgY/0x08) * 0x20), p.lcdc.GetBgWindowTileData(),
	bgX, p.lcdc.IsBgWindowTileDataSigned(), bgY % 0x08)
}

func (p *PixelTransfer) FetchWindow() {
	winX := (p.x - p.memRegs.Get(WX) + 7) / 0x08
	winY := p.memRegs.Get(LY) - p.memRegs.Get(WY)

	p.fetcher.Fetch(p.lcdc.GetWindowTileMapDisplay() + ((winY/0x08) * 0x20), p.lcdc.GetBgWindowTileData(),
	winX, p.lcdc.IsBgWindowTileDataSigned(), winY % 0x08)
}
