package gameboy

import "fmt"

type MemRegisterType struct {
	RegType int
	Addr    int
}

type MemoryRegisters struct {
	Registers map[int]MemRegisterType
	Values    map[int]int
}

func NewMemRegisters(regs ...MemRegisterType) MemoryRegisters {

	mr := MemoryRegisters{}

	m := make(map[int]MemRegisterType)
	for _, r := range regs {
		if a, ok := m[r.Addr]; ok {
			panic(fmt.Sprintf("Address repeated between registers: %d", a.Addr))
		}

		m[r.Addr] = r
		mr.Values[r.Addr] = 0
	}

	mr.Registers = m
	return mr
}

func (m *MemoryRegisters) Get(regAddr int) int {
	if _, ok := m.Registers[regAddr]; ok {
		return m.Values[regAddr]
	} else {
		panic(fmt.Sprintf("Invalid register: address %x", regAddr))
	}
}

func (m *MemoryRegisters) Put(regAddress int, value int) {
	if _, ok := m.Registers[regAddress]; ok {
		m.Values[regAddress] = value
	} else {
		panic(fmt.Sprintf("Invalid register: address %x", regAddress))
	}
}

func (m *MemoryRegisters) Accepts(addr int) bool {
	if _, ok := m.Registers[addr]; ok {
		return true
	} else {
		return false
	}
}

func (m *MemoryRegisters) SetByte(addr int, value int) {
	if m.Registers[addr].RegType == W || m.Registers[addr].RegType == RW {
		m.Values[addr] = value
	}
}

func (m *MemoryRegisters) GetByte(addr int) int {
	if m.Registers[addr].RegType == R || m.Registers[addr].RegType == RW {
		return m.Values[addr]
	} else {
		return 0xff
	}
}

func (m *MemoryRegisters) PreIncrement(regAddr int) int {
	if _, ok := m.Registers[regAddr]; ok {
		v := m.Values[regAddr] + 1
		m.Values[regAddr] = v
		return v
	} else {
		panic(fmt.Sprintf("Invalid register: address %x", regAddr))
	}
}
