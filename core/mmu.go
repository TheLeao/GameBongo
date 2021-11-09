package core

type Mmu struct {
	void   AddressSpace
	spaces []AddressSpace
}

//interface
func (m Mmu) Accepts(addr int) bool {
	return true
}

func (m Mmu) SetByte(addr int, value int) {
	for i, s := range m.spaces {
		if s.Accepts(addr) {
			m.spaces[i].SetByte(addr, value)
			return
		}
	}
}

func (m Mmu) GetByte(addr int) int {
	for i, s := range m.spaces {
		if s.Accepts(addr) {
			return m.spaces[i].GetByte(addr)
		}
	}

	return VoidAddressSpace{}.GetByte(addr)
}

func (m Mmu) Space(addr int) AddressSpace {
	for _, s := range m.spaces {
		if s.Accepts(addr) {
			return s
		}
	}

	return VoidAddressSpace{}
}
