package cpu

import (
	"fmt"

	"github.com/theleao/gamebongo/gameboy"
)

var Arguments []Argument

type Argument struct {
	Label    string
	OprndLen int
	IsMemory bool
	DataType int
	readFn   func(*Registers, gameboy.AddressSpace, []int, int) int
	writeFn  func(*Registers, gameboy.AddressSpace, []int, int)
}

func InitializeArguments() {
	Arguments = append(Arguments, NewArgument("A", 0, false, D8))
	Arguments = append(Arguments, NewArgument("B", 0, false, D8))
	Arguments = append(Arguments, NewArgument("C", 0, false, D8))
	Arguments = append(Arguments, NewArgument("D", 0, false, D8))
	Arguments = append(Arguments, NewArgument("E", 0, false, D8))
	Arguments = append(Arguments, NewArgument("H", 0, false, D8))
	Arguments = append(Arguments, NewArgument("L", 0, false, D8))
	Arguments = append(Arguments, NewArgument("AF", 0, false, D16))
	Arguments = append(Arguments, NewArgument("BC", 0, false, D16))
	Arguments = append(Arguments, NewArgument("HL", 0, false, D16))
	Arguments = append(Arguments, NewArgument("SP", 0, false, D16))
	Arguments = append(Arguments, NewArgument("PC", 0, false, D16))
	Arguments = append(Arguments, NewArgument("d8", 1, false, D8))
	Arguments = append(Arguments, NewArgument("d16", 2, false, D16))
	Arguments = append(Arguments, NewArgument("r8", 1, false, R8))
	Arguments = append(Arguments, NewArgument("a16", 2, false, D16))
	Arguments = append(Arguments, NewArgument("(BC)", 0, true, D8))
	Arguments = append(Arguments, NewArgument("(DE)", 0, true, D8))
	Arguments = append(Arguments, NewArgument("(HL)", 0, true, D8))
	Arguments = append(Arguments, NewArgument("(a8)", 1, true, D8))
	Arguments = append(Arguments, NewArgument("(a16)", 2, true, D8))
	Arguments = append(Arguments, NewArgument("(C)", 0, true, D8))
}

func NewArgument(label string, oprndLen int, isMemory bool, dataType int) Argument {
	a := Argument{
		Label:    label,
		OprndLen: oprndLen,
		IsMemory: isMemory,
		DataType: dataType,
	}

	switch label {
	case "A":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.A
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.A = c
		}
	case "B":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.B
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.B = c
		}
	case "C":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.C
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.C = c
		}
	case "D":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.D
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.D = c
		}
	case "E":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.E
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.E = c
		}
	case "H":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.H
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			*&r.H = c
		}
	case "L":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.L
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			*&r.L = c
		}
	case "AF":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.getAF()
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.setAF(c)
		}
	case "BC":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.getBC()
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.setBC(c)
		}
	case "DE":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.getDE()
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.setDE(c)
		}
	case "HL":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.getHL()
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.setHL(c)
		}
	case "SP":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.Sp
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.Sp = c
		}
	case "PC":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return r.Pc
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			r.Pc = c
		}
	case "d8":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return args[0]
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			panic(fmt.Sprintf("Invalid operation: %s", label))
		}
	case "d16":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return ToWordBytes(args)
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			panic(fmt.Sprintf("Invalid operation: %s", label))
		}
	case "r8":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return ToSigned(args[0])
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			panic(fmt.Sprintf("Invalid operation: %s", label))
		}
	case "a16":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return ToWordBytes(args)
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			panic(fmt.Sprintf("Invalid operation: %s", label))
		}
	case "(BC)":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return a.GetByte(r.getBC())
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			a.SetByte(r.getBC(), c)
		}
	case "(DE)":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return a.GetByte(r.getDE())
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			a.SetByte(r.getDE(), c)
		}
	case "(HL)":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return a.GetByte(r.getHL())
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			a.SetByte(r.getHL(), c)
		}
	case "(a8)":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return a.GetByte(0xff00 | args[0])
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			a.SetByte(r.getBC(), c)
		}
	case "(a16)":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return ToWordBytes(args)
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			a.SetByte(ToWordBytes(args), c)
		}
	case "(C)":
		a.readFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) int {
			return a.GetByte(0xff00 | r.C)
		}
		a.writeFn = func(r *Registers, a gameboy.AddressSpace, args []int, c int) {
			a.SetByte(0xff00|r.C, c)
		}
	}

	return a
}

func (a *Argument) Read(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) int {
	return a.readFn(reg, addr, args, cntxt)
}

func (a *Argument) Write(reg *Registers, addr gameboy.AddressSpace, args []int, cntxt int) {
	a.writeFn(reg, addr, args, cntxt)
}

func GetArgument(src string) Argument {
	for _, a := range Arguments {
		if a.Label == src {
			return a
		}
	}

	panic(fmt.Sprintf("Argument not found: %s", src))
}
