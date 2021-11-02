package gpu

type Lcdc struct {
	value int
}

func NewLcdc() Lcdc {
	return Lcdc{
		value: 0x91,
	}
}

func (l *Lcdc) Accepts(addr int) bool {
	return addr == 0xff408
}

func (l *Lcdc) SetByte(addr int, value int) {
	l.value = value
}

func (l *Lcdc) GetByte(addr int) int {
	return l.value
}

func (l *Lcdc) IsBgAndWindowDisplay() bool {
	return (l.value & 0x01) != 0
}

func (l *Lcdc) IsObjDisplay() bool {
	return l.value&0x02 != 0
}

func (l *Lcdc) GetSpriteHeight() int {
	if l.value&0x04 == 0 {
		return 8
	} else {
		return 16
	}
}

func (l *Lcdc) GetBgTileMapDisplay() int {
	if l.value&0x08 == 0 {
		return 0x9800
	} else {
		return 0x9c00
	}
}

func (l *Lcdc) GetBgWindowTileData() int {
	if l.value&0x10 == 0 {
		return 0x9000
	} else {
		return 0x8000
	}
}

func (l *Lcdc) IsBgWindowTileDataSigned() bool {
	return l.value&0x10 == 0
}

func (l *Lcdc) IsWindowDisplay() bool {
	return l.value&0x20 != 0
}

func (l *Lcdc) GetWindowTileMapDisplay() int {
	if l.value&0x40 == 0 {
		return 0x9800
	} else {
		return 0x9c00
	}
}

func (l *Lcdc) IsLcdEnabled() bool {
	return l.value&0x80 != 0
}

func (l *Lcdc) SetValue(value int) {
	l.value = value
}

func (l *Lcdc) GetValue() int {
	return l.value
}
