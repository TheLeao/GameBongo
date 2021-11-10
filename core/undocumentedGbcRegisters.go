package core

type UndocumentedGbcRegisters struct {
	ram   Ram
	xff6c int
}

func NewUndocumentedGbcRegisters() UndocumentedGbcRegisters {
	r := Ram{
		Offset: 0xff72,
		Length: 6,
	}

	r.SetByte(0xff74, 0xff)
	r.SetByte(0xff75, 0x8f)

	return UndocumentedGbcRegisters{
		ram:   r,
		xff6c: 0xfe,
	}
}

//interface

func (u *UndocumentedGbcRegisters) Accepts(addr int) bool {
	return addr == 0xff6c || u.ram.Accepts(addr)
}

func (u *UndocumentedGbcRegisters) SetByte(addr int, value int) {
	switch addr {
	case 0xff6c:
		u.xff6c = 0xfe | (value & 1)
	case 0xff72:
	case 0xff73:
	case 0xff74:
		u.ram.SetByte(addr, value)
	case 0xff75:
		u.ram.SetByte(addr, 0x8f|(value&0b01110000))
	}
}

func (u *UndocumentedGbcRegisters) GetByte(addr int) int {
	if addr == 0xff6c {
		return u.xff6c
	} else if u.ram.Accepts(addr) {
		return u.ram.GetByte(addr)
	} else {
		panic("UndocumentedGbcRegisters - GetByte - illegal argument")
	}
}