package gpu

import "github.com/theleao/goingboy/core"

type IntQueue []int

func (q IntQueue) enqueue(v int) []int {
	return append(q, v)
}

func (q IntQueue) dequeue() ([]int, int) {
	n := q[0]
	return q[1:], n
}

func (q IntQueue) clear() []int {
	return make([]int, len(q))
}

type PixelQueue interface {
	Length() int
	PutPixelToScreen()
	DropPixel()
	Enqueue8Pixels(pixels []int, ta TileAttribute)
	SetOverlay(pixelLine []int, offset int, ta TileAttribute, oamIndex int)
	Clear()
}

//COLOR PIXEL QUEUE

type ColorPixelQueue struct {
	lcdc       Lcdc
	display    Display
	bgPalette  ColorPalette
	oamPalette ColorPalette
	pixels     IntQueue
	palettes   IntQueue
	priorities IntQueue
}

func (c *ColorPixelQueue) Length() int {
	return len(c.pixels)
}

func (c *ColorPixelQueue) dequeuePixel() int {
	var x, y, z int
	c.priorities, x = c.priorities.dequeue()
	c.palettes, y = c.palettes.dequeue()
	c.pixels, z = c.pixels.dequeue()
	return c.getColor(x, y, z)
}

func (c *ColorPixelQueue) getColor(priority int, palette int, color int) int {
	if priority >= 0 && priority < 10 {
		return c.oamPalette.GetPallete(palette)[color]
	} else {
		return c.bgPalette.GetPallete(palette)[color]
	}
}

func (c *ColorPixelQueue) PutPixelToScreen() {
	c.display.PutColorPixel(c.dequeuePixel())
}

func (c *ColorPixelQueue) DropPixel() {
	c.dequeuePixel()
}

func (c *ColorPixelQueue) Enqueue8Pixels(pxLine []int, tileAttr TileAttribute) {

	for _, p := range pxLine {
		c.pixels.enqueue(p)
		c.palettes.enqueue(tileAttr.ColorPaletteIndex())
		if tileAttr.IsPriority() {
			c.priorities.enqueue(100)
		} else {
			c.priorities.enqueue(-1)
		}
	}
}

func (c *ColorPixelQueue) SetOverlay(pxLine []int, offset int, tileAttr TileAttribute, oamIndex int) {

	for i := 0; i < len(pxLine); i++ {

		p := pxLine[i]
		j := i - offset
		if p == 0 {
			continue //color 0 is transparent
		}
		oldPriority := c.priorities[j]
		put := false

		if (oldPriority == -1 || oldPriority == 100) && !c.lcdc.IsBgAndWindowDisplay() {
			put = true
		} else if oldPriority == 100 {
			put = c.pixels[j] == 0
		} else if oldPriority == -1 && !tileAttr.IsPriority() {
			put = true
		} else if oldPriority == -1 && !tileAttr.IsPriority() && c.pixels[j] == 0 {
			put = true
		} else if oldPriority >= 0 && oldPriority < 10 { //other sprite than bg
			put = oldPriority > oamIndex
		}

		if put {
			c.pixels[j] = p
			c.palettes[j] = tileAttr.ColorPaletteIndex()
			c.priorities[j] = oamIndex
		}
	}
}

func (c *ColorPixelQueue) Clear() {
	c.pixels = c.pixels.clear()
	c.palettes = c.palettes.clear()
	c.priorities = c.priorities.clear()
}

//DMG PIXEL QUEUE

type DmgPixelQueue struct {
	Display   Display
	Regs      core.MemoryRegisters
	Pixels    IntQueue
	Palettes  IntQueue
	PixelType IntQueue
}

func (d *DmgPixelQueue) Length() int {
	return len(d.Pixels)
}

func (d *DmgPixelQueue) PutPixelToScreen() {
	var p int
	d.PixelType, p = d.PixelType.dequeue()
	d.Display.PutDmgPixel(p)
}

func (d *DmgPixelQueue) DropPixel() {
	d.DequeuePixel()
}

//not interface
func (d *DmgPixelQueue) getColor(palette int, colorIndex int) int {
	return 0b11 & (palette >> (colorIndex * 2))
}

//not interface
func (d *DmgPixelQueue) DequeuePixel() int {
	d.PixelType, _ = d.PixelType.dequeue()
	var pal, pix int
	d.Palettes, pal = d.Palettes.dequeue()
	d.Pixels, pix = d.Pixels.dequeue()
	return d.getColor(pal, pix)
}

func (d *DmgPixelQueue) Enqueue8Pixels(pixels []int, ta TileAttribute) {
	for _, p := range pixels {
		d.Pixels = d.Pixels.enqueue(p)
		gpuReg, _ := GetGpuRegister(BGP)
		reg := d.Regs.Get(gpuReg)
		d.Palettes = d.Palettes.enqueue(reg)
		d.PixelType = d.PixelType.enqueue(0)
	}
}

func (d *DmgPixelQueue) SetOverlay(pixelLine []int, offset int, ta TileAttribute, oamIndex int) {
	priority := ta.IsPriority()
	dmgAddr, _ := ta.DmgPalette()
	overlayPalette := d.Regs.Get(dmgAddr)

	for i := 0; i < len(pixelLine); i++ {
		pl := pixelLine[i]
		j := i - offset
		if d.PixelType[j] == 1 {
			continue
		}
		if (priority && d.Pixels[j] == 0) || !priority && pl != 0 {
			d.Pixels[j] = pl
			d.Palettes[j] = overlayPalette
			d.PixelType[j] = 1
		}
	}
}

func (d *DmgPixelQueue) Clear() {
	d.Pixels = d.Pixels.clear()
	d.Palettes = d.Palettes.clear()
	d.PixelType = d.PixelType.clear()
}
