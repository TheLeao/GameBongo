package gpu

import (
	"github.com/theleao/goingboy/core"
)

const ( //don't mess with this order
	ReadTileId = iota
	ReadData1
	ReadData2
	Push
	ReadSpriteTileId
	ReadSpriteFlags
	ReadSpriteData1
	ReadSpriteData2
	PushSprite
)

var EmptyPixelLine [8]int

type Fetcher struct {
	state   int
	vRam0   core.AddressSpace
	vRam1   core.AddressSpace
	oemRam  core.AddressSpace
	memRegs core.MemoryRegisters
	lcdc Lcdc
	gbc bool
	fetchDisabled bool
	pixelLine []int
	fifo PixelQueue
	mapAddr int
	xOffset int
	tileDtAddr int
	tileIdSigned bool
	tileLine int
	tileId int
	tileAttr TileAttribute
	tileData1 int
	tileData2 int
	spriteTileLine int
	sprite SpritePosition
	spriteAttr TileAttribute
	spriteOffset int
	spriteOamIndex int
	divider int
}

func NewFetcher(pq PixelQueue, vram0 core.AddressSpace, vram1 core.AddressSpace, oem core.AddressSpace,
	l Lcdc, regs core.MemoryRegisters, gbc bool) Fetcher {
		return Fetcher{
			fifo: pq,
			vRam0: vram0,
			vRam1: vram1,
			oemRam: oem,
			lcdc: l,
			memRegs: regs,
			gbc: gbc,
		}
}

func (f *Fetcher) Init() {
	f.state = ReadTileId
	f.tileId = 0
	f.tileData1 = 0
	f.tileData2 = 0
	f.divider = 2
	f.fetchDisabled = false
}

func (f *Fetcher) Fetch(mapAddr int, tileDtAddr int, xOffset int, tileIdSigned bool, tileLine int) {
	f.mapAddr = mapAddr
	f.tileDtAddr = tileDtAddr
	f.xOffset = xOffset
	f.tileIdSigned = tileIdSigned
	f.tileLine = tileLine
	
	f.fifo.Clear()
	f.state = ReadTileId
	f.tileId = 0
	f.tileData1 = 0
	f.tileData2 = 0
	f.divider = 2
}

func (f *Fetcher) AddSprite(spr SpritePosition, offset int, oamIndex int) {
	f.sprite = spr
	f.state = ReadSpriteTileId
	f.spriteTileLine = f.memRegs.Get(LY) + 16 - spr.Y
	f.spriteOffset = offset
	f.spriteOamIndex = oamIndex
}

func (f *Fetcher) Tick() {
	if f.fetchDisabled && f.state == ReadTileId {
		if f.fifo.Length() <= 8 {
			f.fifo.Enqueue8Pixels(EmptyPixelLine[:], f.tileAttr)
		}

		return
	}

	f.divider -= 1
	if f.divider == 0 {
		f.divider = 2
	} else {
		return
	}

	switch f.state {
	case ReadTileId:
		f.tileId = f.vRam0.GetByte(f.mapAddr + f.xOffset)
		if f.gbc {
			f.tileAttr = TileAttributes[f.vRam1.GetByte(f.mapAddr + f.xOffset)]
		} else {
			f.tileAttr = EmptyTile
		}
		f.state = ReadData1
	case ReadData1:
		f.tileData1 = f.TileData(f.tileId, f.tileLine, 0, f.tileDtAddr, f.tileIdSigned, f.tileAttr, 8)
		f.state = ReadData2
	case ReadData2:
		f.tileData2 = f.TileData(f.tileId, f.tileLine, 1, f.tileDtAddr, f.tileIdSigned, f.tileAttr, 8)
		f.state = Push
		fallthrough
	case Push:
		if f.fifo.Length() <= 8 {
			z := f.ZipThis(f.tileData1, f.tileData2, f.tileAttr.IsXFlip())
			f.fifo.Enqueue8Pixels(z, f.tileAttr)
			f.xOffset = (f.xOffset + 1) % 0x20
			f.state = ReadTileId
		}
	case ReadSpriteTileId:
		f.tileId = f.oemRam.GetByte(f.sprite.Address + 2)
		f.state = ReadSpriteFlags
	case ReadSpriteFlags:
		f.spriteAttr = TileAttributes[f.oemRam.GetByte(f.sprite.Address + 3)]
		f.state = ReadSpriteData1
	case ReadSpriteData1:
		if f.lcdc.GetSpriteHeight() == 16 {
			f.tileId = f.tileId & 0xfe
		}
		f.tileData1 = f.TileData(f.tileId, f.spriteTileLine, 0, 0x8000, false, f.spriteAttr, f.lcdc.GetSpriteHeight())
		f.state = ReadSpriteData2
	case ReadSpriteData2:
		f.tileData2 = f.TileData(f.tileId, f.spriteTileLine, 1, 0x8000, false, f.spriteAttr, f.lcdc.GetSpriteHeight())
		f.state = PushSprite
	case PushSprite:
		z := f.ZipThis(f.tileData1, f.tileData2, f.spriteAttr.IsXFlip())
		f.fifo.SetOverlay(z, f.spriteOffset, f.spriteAttr, f.spriteOamIndex);
        f.state = ReadTileId
	}
}

func (f *Fetcher) TileData(tileId int, line int, byteNum int, tileDtAddr int, signed bool, tileAttr TileAttribute, tileHeight int) int {
	var effectiveLine, tileAddr int

	if tileAttr.IsYFlip() {
		effectiveLine = tileHeight - 1 - line
	} else {
		effectiveLine = line
	}

	if signed {
		tileAddr = tileDtAddr + toSigned(tileId) * 0x10
	} else {
		tileAddr = tileDtAddr + tileId * 0x10
	}

	var vRam core.AddressSpace
	if !f.gbc || tileAttr.Bank() == 0 {
		vRam = f.vRam0
	} else {
		vRam = f.vRam1
	}

	return vRam.GetByte(tileAddr + effectiveLine*2 + byteNum)
}

func (f Fetcher) SpriteInProgress() bool {
	return f.state >= ReadSpriteTileId && f.state <= PushSprite
}

func toSigned(byteValue int) int {
	if byteValue & (1 << 7) == 0 {
		return byteValue
	} else {
		return byteValue - 0x100
	}
}

func (f *Fetcher) ZipThis(dt1 int, dt2 int, reverse bool) []int {
	for i := 7; i >= 0; i-- {
		mask := 1 << i
		
		p := 0
		if dt2 & mask != 0 {
			p = 2
		}
		if dt1 & mask != 0 {
			p += 1
		}

		if reverse {
			f.pixelLine[i] = p
		} else {
			f.pixelLine[7 - i] = p
		}
	}

	return f.pixelLine
}

func (f Fetcher) Zip(dt1 int, dt2 int, reverse bool, pxLine []int) []int {
	for i := 7; i >= 0; i-- {
		mask := 1 << i
		
		p := 0
		if dt2 & mask != 0 {
			p = 2
		}
		if dt1 & mask != 0 {
			p += 1
		}

		if reverse {
			pxLine[i] = p
		} else {
			pxLine[7 - i] = p
		}
	}

	return pxLine
}