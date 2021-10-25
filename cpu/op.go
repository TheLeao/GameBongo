package cpu

import "github.com/theleao/gamebongo/gameboy"

//Base Op
type Op struct {
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

type LoadOp struct {
	arg Argument
}

func (l *LoadOp) ReadsMemory() {

}
