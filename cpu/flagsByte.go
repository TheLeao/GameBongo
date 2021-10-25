package cpu

const (
	POS_C = 4
	POS_H = 5
	POS_N = 6
	POS_Z = 7
)

type Flags struct {
	flagsByte int
}

func (f *Flags) IsZ() bool {
	return GetBit(f.flagsByte, POS_Z)
}

func (f *Flags) IsN() bool {
	return GetBit(f.flagsByte, POS_N)
}

func (f *Flags) IsH() bool {
	return GetBit(f.flagsByte, POS_H)
}

func (f *Flags) IsC() bool {
	return GetBit(f.flagsByte, POS_C)
}

func (f *Flags) SetZ(z bool) {
	f.flagsByte = SetBitValue(f.flagsByte, POS_Z, z)
}

func (f *Flags) SetN(n bool) {
	f.flagsByte = SetBitValue(f.flagsByte, POS_N, n)
}

func (f *Flags) SetH(h bool) {
	f.flagsByte = SetBitValue(f.flagsByte, POS_H, h)
}

func (f *Flags) SetC(c bool) {
	f.flagsByte = SetBitValue(f.flagsByte, POS_C, c)
}

func (f *Flags) SetFlagsByte(value int) {
	f.flagsByte = value & 0xf0
}

func (f *Flags) GetFlagsByte() int {
	return f.flagsByte
}

func (f *Flags) GetString() string {
	var z, n, h, c string

	if f.IsZ() {
		z = "Z"
	} else {
		z = "-"
	}

	if f.IsN() {
		n = "N"
	} else {
		n = "-"
	}

	if f.IsH() {
		h = "H"
	} else {
		h = "-"
	}

	if f.IsC() {
		c = "c"
	} else {
		c = "-"
	}

	return z + n + h + c + "----"
}
