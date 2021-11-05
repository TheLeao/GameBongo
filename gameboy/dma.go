package gameboy

type Dma struct {
	transfInProgress bool
	restarted        bool
	from             int
	ticks            int
	regValue         int
	addrSpace        DmaAddressSpace
	oam              AddressSpace
	speedMode        SpeedMode
}

type DmaAddressSpace struct {
	addrSpace AddressSpace
}

func (d *DmaAddressSpace) GetByte(addr int) int {
	if addr < 0xe000 {
		return d.addrSpace.GetByte(addr)
	} else {
		return d.addrSpace.GetByte(addr - 0x2000)
	}
}

func (d *DmaAddressSpace) SetByte(addr int, value int) {
	panic("Not supported")
}

func (d *DmaAddressSpace) Accepts(addr int) bool {
	return true
}

func NewDma(addr AddressSpace, oam AddressSpace, spd SpeedMode) Dma {
	dmaAddr := DmaAddressSpace{
		addrSpace: addr,
	}

	return Dma{
		addrSpace: dmaAddr,
		speedMode: spd,
		oam:       oam,
		regValue:  0xff,
	}
}

func (d *Dma) Tick() {
	if d.transfInProgress {
		d.ticks += 1
		if d.ticks >= 648/d.speedMode.GetSpeedMode() {
			d.transfInProgress = false
			d.restarted = false
			d.ticks = 0
			for i := 0; i < 0xa0; i++ {
				d.oam.SetByte(0xfe00+i, d.addrSpace.GetByte(d.from+i))
			}
		}
	}
}

func (d *Dma) Accepts(addr int) bool {
	return addr == 0xff46
}

func (d *Dma) SetByte(addr int, value int) {
	d.from = value * 0x100
	d.restarted = d.IsOamBlocked()
	d.ticks = 0
	d.transfInProgress = true
	d.regValue = value
}

func (d *Dma) GetByte(addr int) int {
	return d.regValue
}

func (d *Dma) IsOamBlocked() bool {
	return d.restarted || (d.transfInProgress && d.ticks >= 5)
}
