package cpu

import (
	"fmt"

	"github.com/theleao/gamebongo/gameboy"
)

//Base Op
type Op interface {
	ReadsMemory() bool
	WritesMemory() bool
	CausesOemBug(reg *Registers, opCntxt int) (bool, int)
	Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int
	SwitchInterrupts(intrpt Interrupter)
	Proceed(reg Registers) bool
	ForceFinishCycle() bool
	OperandLength() int
	InOamArea(addr int) bool
	GetString() string
}

type op struct {
	Op
}

func NewOp() Op {
	return &op{}
}

func (o *op) ReadsMemory() bool {
	return false
}

func (o *op) WritesMemory() bool {
	return false
}

func (o *op) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	return cntxt
}

func (o *op) SwitchInterrupts(intrpt Interrupter) {
}

func (o *op) Proceed(reg Registers) bool {
	return true
}

func (o *op) ForceFinishCycle() bool {
	return false
}

func (o *op) OperandLength() int {
	return 0
}

func (o *op) InOamArea(addr int) bool {
	return addr >= 0xff00 && addr <= 0xfeff
}

func (o *op) GetString() string {
	panic("Wrong Op call of method GetString")
}

//OPs to implement/override the interface

//LOAD Op
type LoadOp struct {
	arg Argument
	Op
}

func NewLoadOp(a Argument) Op {
	return &LoadOp{
		arg: a,
		Op: NewOp(),
	}
}

//"overriding"
func (l *LoadOp) ReadsMemory() bool {
	return l.arg.IsMemory
}

func (l *LoadOp) OperandLength() int {
	return l.arg.OprndLen
}

func (l *LoadOp) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	return l.arg.Read(reg, addr, args, cntxt)
}

func (l *LoadOp) GetString() string{
	if l.arg.DataType == D16 {
		return fmt.Sprintf("%s → [__]", l.arg.Label)
	} else {
		return fmt.Sprintf("%s → [_]", l.arg.Label)
	}
}

//LOAD WORD Op
type LoadWordOp struct {
	value int
	Op
}

func NewLoadWordOp(val int) Op {
	return &LoadWordOp{
		value: val,
		Op: NewOp(),
	}
} 

func (w *LoadWordOp) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	return w.value
}

func (w *LoadWordOp) GetString() string {
	return fmt.Sprintf("0x%02X → [__]", w.value)
}

//PROCEED IF Op
type ProceedIfOp struct {
	condition string
	Op
}

func NewProceedIfOp(c string) Op {
	return &ProceedIfOp{
		condition: c,
		Op: NewOp(),
	}
} 

func (p *ProceedIfOp) Proceed(r Registers) bool {
	switch p.condition {
	case "NZ":
		return !r.Flags.IsZ()
	case "Z":
		return r.Flags.IsZ()
	case "NC":
		return !r.Flags.IsC()
	case "C":
		return r.Flags.IsC()
	default:
		return false
	}
}

func (p *ProceedIfOp) GetString() string {
	return fmt.Sprintf("? %s:", p.condition)
}

//PUSH Op

type PushOp1 struct {
	Op
}

func NewPushOp() Op {
	return &PushOp1{
		Op: NewOp(),
	}
}

func (n *PushOp1) WritesMemory() bool {
	return true
}