package gpu

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

type ColorPixelQueue struct {
	lcdc       Lcdc
	display    Display
	bgPalette  ColorPalette
	oamPalette ColorPalette
	pixels     IntQueue
	palettes   IntQueue
	priorities IntQueue
}

func NewColorPixelQueue(l Lcdc, d Display, bgP ColorPalette, oamP ColorPalette) ColorPixelQueue {
	return ColorPixelQueue{
		lcdc:       l,
		display:    d,
		bgPalette:  bgP,
		oamPalette: oamP,
	}
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