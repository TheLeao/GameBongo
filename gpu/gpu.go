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

type Lcdc struct {
	Enabled bool
}
