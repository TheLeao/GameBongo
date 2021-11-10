package core

import "fmt"

type AddressSpace interface {
	Accepts(addr int) bool
	SetByte(addr int, value int)
	GetByte(addr int) int
}

type VoidAddressSpace struct {
}

func (v VoidAddressSpace) Accepts(addr int) bool {
	return true
}

func (VoidAddressSpace) SetByte(addr int, value int) {
	if addr < 0 || addr > 0xffff {
		panic(fmt.Sprintf("Invalid address: %x", addr))
	}
}

func (VoidAddressSpace) GetByte(addr int) int {
	if addr < 0 || addr > 0xffff {
		panic(fmt.Sprintf("Invalid address: %x", addr))
	}

	return 0xff
}

type ShadowAddressSpace struct {
	addrSpace AddressSpace
	echoStart int
	targetStart int
	length int
}

//interface

func (s *ShadowAddressSpace) Accepts(addr int) bool {
	return addr >= s.echoStart && addr < (s.echoStart + s.length)
}

func (s *ShadowAddressSpace) SetByte(addr int, value int) {
	s.addrSpace.SetByte(addr, value)
}

func (s *ShadowAddressSpace) GetByte(addr int) int {
	return s.addrSpace.GetByte(s.translate(addr))
}

//
func (s *ShadowAddressSpace) translate(addr int) int {
	i := addr - s.echoStart
	if i < 0 || i >= s.length {
		panic("ShadowAddressSpace - translate - illegal argument")
	}
	return i + s.targetStart
}