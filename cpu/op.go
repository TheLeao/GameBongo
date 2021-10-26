package cpu

import (
	"fmt"

	"github.com/theleao/gamebongo/gameboy"
	"github.com/theleao/gamebongo/gpu"
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
	return l.arg.Read(reg, addr, args)
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

//PUSH Op 1

type PushOp1 struct {
	fn IntRegistryFunc
	Op
}

func NewPushOp1(f IntRegistryFunc) Op {
	return &PushOp1{
		fn: f,
		Op: NewOp(),
	}
}

func (p *PushOp1) WritesMemory() bool {
	return true
}

func (p *PushOp1) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	reg.Sp = p.fn(&reg.Flags, reg.Sp)
	addr.SetByte(reg.Sp, (cntxt & 0xff00) >> 8)
	return cntxt
}

func (p *PushOp1) CausesOemBug(reg *Registers, opCntxt int) (bool, int) {
	if p.InOamArea(reg.Sp) {
		return true, gpu.PUSH_1
	} else {
		return false, -1
	}
}

func (p *PushOp1) GetString() string {
	return "[_ ] → (SP--)"
}

//PUSH Op 2

type PushOp2 struct {
	fn IntRegistryFunc
	Op
}

func NewPushOp2(f IntRegistryFunc) Op {
	return &PushOp1{
		fn: f,
		Op: NewOp(),
	}
}

func (p *PushOp2) WritesMemory() bool {
	return true
}

func (p *PushOp2) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	reg.Sp = p.fn(&reg.Flags, reg.Sp)
	addr.SetByte(reg.Sp, cntxt & 0xff00)
	return cntxt
}

func (p *PushOp2) CausesOemBug(reg *Registers, opCntxt int) (bool, int) {
	if p.InOamArea(reg.Sp) {
		return true, gpu.PUSH_2
	} else {
		return false, -1
	}
}

func (p *PushOp2) GetString() string {
	return "[ _] → (SP--)"
}


//POP Op 1

type PopOp1 struct {
	fn IntRegistryFunc
	Op
}

func NewPopOp1(f IntRegistryFunc) Op {
	return &PopOp1{
		fn: f,
		Op: NewOp(),
	}
}

func (p *PopOp1) ReadsMemory() bool {
	return true
}

func (p *PopOp1) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	lsb := addr.GetByte(reg.Sp)
	reg.Sp = p.fn(&reg.Flags, reg.Sp)
	return lsb
}

func (p *PopOp1) CausesOemBug(reg *Registers, opCntxt int) (bool, int) {
	if p.InOamArea(reg.Sp) {
		return true, gpu.POP_1
	} else {
		return false, -1
	}
}

func (p *PopOp1) GetString() string {
	return "(SP++) → [ _]"
}

//POP Op 2

type PopOp2 struct {
	fn IntRegistryFunc
	Op
}

func NewPopOp2(f IntRegistryFunc) Op {
	return &PopOp1{
		fn: f,
		Op: NewOp(),
	}
}

func (p *PopOp2) ReadsMemory() bool {
	return true
}

func (p *PopOp2) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	msb := addr.GetByte(reg.Sp)
	reg.Sp = p.fn(&reg.Flags, reg.Sp)
	return cntxt | (msb << 8)
}

func (p *PopOp2) CausesOemBug(reg *Registers, opCntxt int) (bool, int) {
	if p.InOamArea(reg.Sp) {
		return true, gpu.POP_2
	} else {
		return false, -1
	}
}

func (p *PopOp2) GetString() string {
	return "(SP++) → [_ ]"
}

//ALU Op 1

type AluOp1 struct {
	fn BiIntRegistryFunc
	arg Argument
	operation string
	lastDataType int
	Op
}

func NewAluOp1(f BiIntRegistryFunc, a Argument, o string, dt int) Op {
	return &AluOp1{
		fn: f,
		arg: a,
		operation: o,
		lastDataType: dt,
		Op: NewOp(),
	}
}

func (a *AluOp1) ReadsMemory() bool {
	return a.arg.IsMemory
}

func (a *AluOp1) OperandLength() int {
	return a.arg.OprndLen
}

func (a *AluOp1) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	v2 := a.arg.Read(reg, addr, args)
	return a.fn(&reg.Flags, cntxt, v2)
}

func (a *AluOp1) GetString() string {
	if a.lastDataType == D16 {
		return fmt.Sprintf("%s([__],%s) → [__]", a.operation, a.arg)
	} else {
		return fmt.Sprintf("%s([_],%s) → [_]", a.operation, a.arg);
	}
}

//ALU Op 2

type AluOp2 struct {
	fn BiIntRegistryFunc	
	operation string
	d8Val int
	Op
}

func NewAluOp2(f BiIntRegistryFunc, o string, d8v int) Op {
	return &AluOp2{
		fn: f,
		operation: o,
		d8Val: d8v,
		Op: NewOp(),
	}
}

func (a *AluOp2) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	return a.fn(&reg.Flags, cntxt, a.d8Val)
}

func (a *AluOp2) GetString() string {
	return fmt.Sprintf("%s(%d,[_]) → [_]", a.operation, a.d8Val)
}

//ALU Op 3

type AluOp3 struct {
	fn IntRegistryFunc
	operation string
	lastDataType int
	Op
}

func NewAluOp3(f IntRegistryFunc, o string, dt int) Op {
	return &AluOp3{
		fn: f,
		operation: o,
		lastDataType: dt,
		Op: NewOp(),
	}
}

func (a *AluOp3) Execute(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	return a.fn(&reg.Flags, cntxt)
}

func (a *AluOp3) CausesOemBug(reg *Registers, opCntxt int) (bool, int) {
	//check if AluOp's function is 'INC' or 'DEC' with a data type value of D16
	hasFunc := false
	if (a.operation == "INC" || a.operation == "DEC") && a.lastDataType == D16{
		hasFunc = true
	}
	//also checking if address (cntx variable) is in range
	if hasFunc && (opCntxt >= 0xfe00 && opCntxt <= 0xfeff) {
		return true, gpu.INC_DEC
	} else {
		return false, -1
	}
}

func (a *AluOp3) GetString() string {
	if a.lastDataType == D16 {
		return a.operation + "([__]) → [__]"
	} else {
		return a.operation + "([_]) → [_]"
	}
}