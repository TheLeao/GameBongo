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

func Test() {

}