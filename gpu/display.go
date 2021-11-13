package gpu

type Display interface {
	EnableLcd()
	DisableLcd()
	PutDmgPixel(color int)
	PutColorPixel(rgb int)
	RequestRefresh()
	WaitForRefresh()
}

type NullDisplay struct {
}

func (n *NullDisplay) PutDmgPixel(color int) {
}

func (n *NullDisplay) PutColorPixel(gbcRgb int) {
}

func (n *NullDisplay) RequestRefresh() {
}

func (n *NullDisplay) WaitForRefresh() {
}

func (n *NullDisplay) EnableLcd() {
}

func (n *NullDisplay) DisableLcd() {
}
