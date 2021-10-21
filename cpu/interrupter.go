package cpu

import "math/bits"

type Interrupter struct {
	ime                      bool
	gbc                      bool
	interruptFlag            int
	interruptEnabled         int
	pendingEnableInterrupts  int
	pendingDisableInterrupts int
}

//Interrupter type constants
const (
	VBLANK = 0x0040 //0
	LCDC   = 0x0048 //1
	TIMER  = 0x0050 //2
	SERIAL = 0x0058 //3
	P1013  = 0x0060 //4
)

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

func (i *Interrupter) enableInterrupts(delay bool) {

}

func (i *Interrupter) disableInterrupts(delay bool) {

}

func (i *Interrupter) isHaltBug() bool {
	return (i.interruptFlag&i.interruptEnabled) != 0 && !i.ime
}

func (i *Interrupter) clearInterrupt(intrptType int) {
	b := bits.Reverse(uint(1 << intrptType))
	i.interruptFlag = i.interruptFlag & int(b)
}

func (i *Interrupter) OnInstructionFinished() {
	if i.pendingEnableInterrupts != -1 {
		i.pendingEnableInterrupts--
		if i.pendingEnableInterrupts == 0 {
			i.enableInterrupts(false)
		}
	}

	if i.pendingDisableInterrupts != -1 {
		i.pendingDisableInterrupts--
		if i.pendingDisableInterrupts == 0 {
			i.disableInterrupts(false)
		}
	}
}
