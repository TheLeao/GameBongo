package gpu

//call inicialization
var TileAttributes []TileAttribute
var EmptyTile TileAttribute

type TileAttribute int

func InitializeTileAttributes() {
	for i := 0; i < 256; i++ {
		TileAttributes[i] = TileAttribute(i)
	}
	EmptyTile = TileAttributes[0]
}

func (t TileAttribute) IsPriority() bool {
	return t & (1 << 7) != 0
}

func (t TileAttribute) IsYFlip() bool {
	return t & (1 << 6) != 0
}

func (t TileAttribute) IsXFlip() bool {
	return t & (1 << 5) != 0
}

func (t TileAttribute) DmgPalette() (int, int) {
	if t & (1 << 4) == 0 {
		return GetGpuRegister(OBP0)
	} else {
		return GetGpuRegister(OBP1)
	}
}

func (t TileAttribute) Bank() int {
	if t & (1 << 3) == 0 {
		return 0
	} else {
		return 1
	}
}

func (t TileAttribute) ColorPaletteIndex() int {
	return int(t) & 0x07
}