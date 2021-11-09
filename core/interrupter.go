package core

type Interrupter struct {
	Ime                      bool
	Gbc                      bool
	InterruptFlag            int
	InterruptEnabled         int
	PendingEnableInterrupts  int
	PendingDisableInterrupts int
}

//Interrupter type constants
const (
	VBLANK = 0x0040 //0
	LCDC   = 0x0048 //1
	TIMER  = 0x0050 //2
	SERIAL = 0x0058 //3
	P1013  = 0x0060 //4
)

func NewInterrupter(Gbc bool) Interrupter {
	return Interrupter{
		Gbc:           Gbc,
		InterruptFlag: 0xe1,
	}
}

func (*Interrupter) Accepts(addr int) bool {
	return addr == 0xff0f || addr == 0xffff
}

func (i *Interrupter) SetByte(addr int, value int) {
	switch addr {
	case 0xff0f:
		i.InterruptFlag = value | 0xe0
	case 0xffff:
		i.InterruptEnabled = value
	}
}

func (i *Interrupter) GetByte(addr int) int {
	switch addr {
	case 0xff0f:
		return i.InterruptFlag
	case 0xffff:
		return i.InterruptEnabled
	default:
		return 0xff
	}
}

func (i *Interrupter) EnableInterrupts(delay bool) {

}

func (i *Interrupter) DisableInterrupts(delay bool) {

}

func (i *Interrupter) IsHaltBug() bool {
	return (i.InterruptFlag&i.InterruptEnabled) != 0 && !i.Ime
}

func (i *Interrupter) ClearInterrupt(intrptType int) {
	b := ^(1 << intrptType)
	i.InterruptFlag = i.InterruptFlag & int(b)
}

func (i *Interrupter) OnInstructionFinished() {
	if i.PendingEnableInterrupts != -1 {
		i.PendingEnableInterrupts--
		if i.PendingEnableInterrupts == 0 {
			i.EnableInterrupts(false)
		}
	}

	if i.PendingDisableInterrupts != -1 {
		i.PendingDisableInterrupts--
		if i.PendingDisableInterrupts == 0 {
			i.DisableInterrupts(false)
		}
	}
}

func (i *Interrupter) RequestInterrupt(intrType int) {
	ord := 0
	switch intrType {
	case VBLANK:
		ord = 0
	case LCDC:
		ord = 1
	case TIMER:
		ord = 2
	case SERIAL:
		ord = 3
	case P1013:
		ord = 4
	}

	i.InterruptFlag = i.InterruptFlag | (1 << ord)
}
