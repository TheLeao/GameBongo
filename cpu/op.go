package cpu

import "github.com/theleao/gamebongo/gameboy"

type Op interface {
	ReadsMemory() bool
	WritesMemory() bool
	CausesOemBug(reg Registers, opCntxt int) (bool, int)
	Execute(reg Registers, addr gameboy.AddressSpace, args [2]int, cntxt int) int
	SwitchInterrupts(intrpt Interrupter)
	Proceed(reg Registers) bool
	ForceFinishCycle() bool
	OperandLength() int
	InOamArea(addr int)
}
