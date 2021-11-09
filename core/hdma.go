package core

const (
	HDMA1 = 0xff51
	HDMA2 = 0xff52
	HDMA3 = 0xff53
	HDMA4 = 0xff54
	HDMA5 = 0xff55
)

type Hdma struct {
	addrSpace        AddressSpace
	hdma1234         Ram
	mode             int
	transfInProgress bool
	hBlankTransfer   bool
	lcdEnabled       bool
	length           int
	src              int
	dst              int
	tick             int
}

func NewHdma(addr AddressSpace) AddressSpace {
	return &Hdma{
		hdma1234: Ram{
			Length: 4,
			Offset: HDMA1,
		},
		addrSpace: addr,
	}
}

//interface

func (h *Hdma) Accepts(addr int) bool {
	return addr >= HDMA1 && addr <= HDMA5 //between 1 and 5
}

func (h *Hdma) SetByte(addr int, value int) {
	if h.hdma1234.Accepts(addr) {
		h.hdma1234.SetByte(addr, value)
	} else if addr == HDMA5 {
		if h.transfInProgress && (addr&(1<<7)) == 0 {
			h.transfInProgress = false
		}
	} else {
		h.startTransfer(value)
	}
}

func (h *Hdma) GetByte(addr int) int {
	if h.hdma1234.Accepts(addr) {
		return 0xff
	} else if addr == HDMA5 {
		if h.transfInProgress {
			return 0
		} else {
			return (1 << 7) | h.length
		}
	} else {
		panic("HDMA illegal argument on GetByte")
	}
}

//

func (h *Hdma) startTransfer(reg int) {
	h.hBlankTransfer = (reg & (1 << 7)) != 0
	h.length = reg & 0x7f

	h.src = (h.hdma1234.GetByte(HDMA1) << 8) | (h.hdma1234.GetByte(HDMA2) & 0xf0)
	h.dst = ((h.hdma1234.GetByte(HDMA3) & 0x1f) << 8) | (h.hdma1234.GetByte(HDMA4) & 0xf0)
	h.src = h.src & 0xfff0
	h.dst = (h.dst & 0x1fff) | 0x8000
	h.transfInProgress = true
}