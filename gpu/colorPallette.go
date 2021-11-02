package gpu

import (
	"bytes"
	"fmt"
)

type ColorPalette struct {
	indexAddr int
	dataAddr  int
	palettes  [8][4]int
	index     int
	autoIncr  bool
}

func NewColorPallete(offset int) ColorPalette {
	return ColorPalette{
		indexAddr: offset,
		dataAddr:  offset + 1,
	}
}

func (c *ColorPalette) Accepts(addr int) bool {
	return addr == c.indexAddr || addr == c.dataAddr
}

func (c *ColorPalette) SetByte(addr int, value int) {
	if addr == c.indexAddr {
		c.index = value & 0x3f
		c.autoIncr = (value & (1 << 7)) != 0
	} else if addr == c.dataAddr {
		color := c.palettes[c.index/8][(c.index%8)/2]
		if c.index%2 == 0 {
			color = (color & 0xff00) | value
		} else {
			color = (color & 0x00ff) | (value << 8)
		}
		c.palettes[c.index/8][(c.index%8)/2] = color

		if c.autoIncr {
			c.index = (c.index + 1) & 0x3f
		}
	} else {
		panic("Error while setting color palette value")
	}
}

func (c *ColorPalette) GetByte(addr int) int {
	if addr == c.indexAddr {
		if c.autoIncr {
			return c.index | 0x80 | 0x40
		} else {
			return c.index | 0x00 | 0x40
		}
	} else if addr == c.dataAddr {
		color := c.palettes[c.index/8][(c.index%8)/2]
		if c.index%2 == 0 {
			return color & 0xff
		} else {
			return (color >> 8) & 0xff
		}
	} else {
		panic("Error to get color palette value")
	}
}

func (c *ColorPalette) GetPallete(i int) [4]int {
	return c.palettes[i]
}

func (c *ColorPalette) GetString() string {
	var s string
	for i := 0; i < 8; i++ {
		var b bytes.Buffer
		b.WriteString(string(i))
		b.WriteString(": ")

		palettes := c.GetPallete(i)

		for _, p := range palettes {
			b.WriteString(fmt.Sprintf("%04x ", p))
		}

		//changing last whitespace for \n
		s := b.String()
		s = s[:len(s)-1] + "\n"
	}

	return s
}

func (c *ColorPalette) FillWithFF() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 4; j++ {
			c.palettes[i][j] = 0x7fff
		}
	}
}
