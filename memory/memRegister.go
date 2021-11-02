package memory

import "fmt"

type MemRegisterType struct {
	regType int
	addr    int
}

type MemRegisters struct {
	registers map[int]MemRegisterType
	values    map[int]int
}

func NewMemRegisters(regs ...MemRegisterType) MemRegisters {

	mr := MemRegisters{}

	m := make(map[int]MemRegisterType)
	for _, r := range regs {
		if a, ok := m[r.addr]; ok {
			panic(fmt.Sprintf("Address repeated between registers: %d", a.addr))
		}

		m[r.addr] = r
		mr.values[r.addr] = 0
	}

	mr.registers = m
	return mr
}

func (m *MemRegisters) Get(regType MemRegisterType) {
	//todo
}

func (m *MemRegisters) Accepts(addr int) bool {
	if _, ok := m.registers[addr]; ok {
		return true
	} else {
		return false
	}
}

func (m *MemRegisters) SetByte(addr int, value int) {
	if m.registers[addr].regType == W || m.registers[addr].regType == RW {
		m.values[addr] = value
	}
}

func (m *MemRegisters) GetByte(addr int) int {
	if m.registers[addr].regType == R || m.registers[addr].regType == RW {
		return m.values[addr]
	} else {
		return 0xff
	}
}

func (m *MemRegisters) PreIncrement(reg MemRegisterType) int {
	if _, ok := m.registers[reg.addr]; ok {
		v := m.values[reg.addr] + 1
		m.values[reg.addr] = v
		return v
	} else {
		panic(fmt.Sprintf("Unvalid register: type %d address %x", reg.regType, reg.addr))
	}
}
