package gpu

type Display interface {
	EnableLcd()
	DisableLcd()
	PutDmgPixel(color int)
	PutColorPixel(rgb int)
	RequestRefresh()
	WaitForRefresh()
}

//aaa

// func (d *Display) PutDmgPixel(color int) {

// }

// func (d *Display) PutColorPixel(gbcRgb int) {

// }

// func (d *Display) RequestRefresh() {

// }

// func (d *Display) WaitForRefresh() {

// }