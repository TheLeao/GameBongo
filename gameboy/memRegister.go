package gameboy

import "fmt"

type MemRegisterType struct {
	regType int
	addr    int
}

type MemoryRegisters struct {
	registers map[int]MemRegisterType
	values    map[int]int
}

func NewMemRegisters(regs ...MemRegisterType) MemoryRegisters {

	mr := MemoryRegisters{}

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

func (m *MemoryRegisters) Get(regAddr int) int {
	if _, ok := m.registers[regAddr]; ok {
		return m.values[regAddr]
	} else {
		panic(fmt.Sprintf("Invalid register: address %x", regAddr))
	}
}

func (m *MemoryRegisters) Put(regAddress int, value int) {
	if _, ok := m.registers[regAddress]; ok {
		m.values[regAddress] = value
	} else {
		panic(fmt.Sprintf("Invalid register: address %x", regAddress))
	}
}

func (m *MemoryRegisters) Accepts(addr int) bool {
	if _, ok := m.registers[addr]; ok {
		return true
	} else {
		return false
	}
}

func (m *MemoryRegisters) SetByte(addr int, value int) {
	if m.registers[addr].regType == W || m.registers[addr].regType == RW {
		m.values[addr] = value
	}
}

func (m *MemoryRegisters) GetByte(addr int) int {
	if m.registers[addr].regType == R || m.registers[addr].regType == RW {
		return m.values[addr]
	} else {
		return 0xff
	}
}

func (m *MemoryRegisters) PreIncrement(regAddr int) int {
	if _, ok := m.registers[regAddr]; ok {
		v := m.values[regAddr] + 1
		m.values[regAddr] = v
		return v
	} else {
		panic(fmt.Sprintf("Invalid register: address %x", regAddr))
	}
}
