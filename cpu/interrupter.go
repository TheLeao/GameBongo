package cpu

type Interrupter struct {
	ime              bool
	gbc              bool
	interruptFlag    int
	interruptEnabled int
}

func NewInterrupter(gbc bool) Interrupter {
	return Interrupter{
		gbc:           gbc,
		interruptFlag: 0xe1,
	}
}

func (*Interrupter) Accepts(addr int) bool {
	return addr == 0xff0f || addr == 0xffff
}

func (i *Interrupter) SetByte(addr int, value int) {
	switch addr {
	case 0xff0f:
		i.interruptFlag = value | 0xe0
	case 0xffff:
		i.interruptEnabled = value
	}
}

func (i *Interrupter) GetByte(addr int) int {
	switch addr {
	case 0xff0f:
		return i.interruptFlag
	case 0xffff:
		return i.interruptEnabled
	default:
		return 0xff
	}
}