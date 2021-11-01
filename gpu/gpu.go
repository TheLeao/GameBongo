package gpu

const ( //GPU Mode
	HBLANK = iota
	VBLANK
	OAMSEARCH
	PIXELTRANSFER
)

type Gpu struct {
	Lcdc        Lcdc
	Mode        int
	TicksInLine int
}

func (g *Gpu) Accepts(addr int) bool {
	return true
}

func (g *Gpu) SetByte(addr int, value int) {
}

func (g *Gpu) GetByte(addr int) int {
}

type Lcdc struct {
	Enabled bool
}
