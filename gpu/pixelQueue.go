package gpu

type IntQueue []int

func (q IntQueue) enqueue(v int) []int {
	return append(q, v)
}

func (q IntQueue) dequeue() ([]int, int) {
	n := q[0]
	return q[1:], n
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

}