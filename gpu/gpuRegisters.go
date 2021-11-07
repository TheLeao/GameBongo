package gpu

const (
	STAT = iota
	SCY
	SCX
	LY
	LYC
	BGP
	OBP0
	OBP1
	WY
	WX
	VBK
)

const ( //memory register types reference
	R = iota
	W
	RW
)

func GpuRegisters() []int {
	return []int{STAT, SCY, SCX, LY, LYC, BGP, OBP0, OBP1, WY, WX, VBK}
}

func GetGpuRegister(reg int) (int, int) {
	switch reg {
	case STAT:
		return 0xff41, RW
	case SCY:
		return 0xff42, RW
	case SCX:
		return 0xff43, RW
	case LY:
		return 0xff44, R
	case LYC:
		return 0xff45, RW
	case BGP:
		return 0xff47, RW
	case OBP0:
		return 0xff48, RW
	case OBP1:
		return 0xff49, RW
	case WY:
		return 0xff4a, RW
	case WX:
		return 0xff4b, RW
	case VBK:
		return 0xff4f, W
	default:
		panic("Wrong register")
	}
}
