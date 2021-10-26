package cpu

type IntRegistryFunc = func(f *Flags, arg int) int
type BiIntRegistryFunc = func(f *Flags, arg1 int, arg2 int) int

type Alu struct {
	funcs   map[AluFunctionKey]IntRegistryFunc
	biFuncs map[AluFunctionKey]BiIntRegistryFunc
}

type AluFunctionKey struct {
	name    string
	dtType1 int
	dtType2 int
}

var ALU Alu
var OEMBUG [2]AluFunctionKey

type OemBug struct {
	fn IntRegistryFunc
}

func InitializeAlu() {
	ALU = NewAlu()
	// OEMBUG = [2]AluFunctionKey{
	// 	NewAluFunctionKey("INC", D16),
	// 	NewAluFunctionKey("DEC", D16),
	// }
}

func NewAluFunctionKey(n string, t1 int) AluFunctionKey {
	return AluFunctionKey{
		name:    n,
		dtType1: t1,
		dtType2: UNSET,
	}
}

func NewAluBiIntFunctionKey(n string, t1 int, t2 int) AluFunctionKey {
	return AluFunctionKey{
		name:    n,
		dtType1: t1,
		dtType2: t2,
	}
}

func (f *AluFunctionKey) Equals(a AluFunctionKey) bool {
	if f.name == a.name {
		return true
	} else if f.dtType1 == a.dtType1 {
		return true
	} else {
		if f.dtType2 != UNSET {
			return f.dtType2 == a.dtType2
		} else {
			return a.dtType2 == UNSET
		}
	}
}

//Create Alu and functions
func NewAlu() Alu {
	f := make(map[AluFunctionKey]IntRegistryFunc)
	bf := make(map[AluFunctionKey]BiIntRegistryFunc)

	//1
	addFunction(f, "INC", D8, func(fl *Flags, arg int) int {
		res := (arg + 1) & 0xff
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH((arg & 0x0f) == 0x0f)
		return res
	})

	//2
	addFunction(f, "INC", D16, func(fl *Flags, arg int) int {
		return (arg + 1) & 0xffff
	})

	//3
	addFunction(f, "DEC", D8, func(fl *Flags, arg int) int {
		res := (arg - 1) & 0xff
		fl.SetZ(res == 0)
		fl.SetN(true)
		fl.SetH((arg & 0x0f) == 0x0)
		return res
	})

	//4
	addFunction(f, "DEC", D16, func(fl *Flags, arg int) int {
		return (arg - 1) & 0xffff
	})

	//5
	addBiIntFunction(bf, "ADD", D16, D16, func(fl *Flags, arg1 int, arg2 int) int {
		fl.SetN(false)
		fl.SetH((arg1&0x0fff)+(arg2+0x0fff) > 0x0fff)
		fl.SetC(arg1+arg2 > 0xffff)
		return (arg1 + arg2) & 0xffff
	})

	//6
	addBiIntFunction(bf, "ADD", D16, R8, func(fl *Flags, arg1 int, arg2 int) int {
		return (arg1 + arg2) & 0xffff
	})

	//7
	addBiIntFunction(bf, "ADD_SP", D16, R8, func(fl *Flags, arg1 int, arg2 int) int {
		fl.SetZ(false)
		fl.SetN(false)
		fl.SetC((((arg1 & 0xff) + (arg2 & 0xff)) & 0x100) != 0)
		fl.SetH((((arg1 & 0x0f) + (arg2 & 0x0f)) & 0x10) != 0)
		return (arg1 + arg2) & 0xffff
	})

	//8
	addFunction(f, "DAA", D8, func(fl *Flags, arg int) int {
		res := arg
		if fl.IsN() {
			if fl.IsH() {
				res = (res - 6) & 0xff
			}

			if fl.IsC() {
				res = (res - 0x60) & 0xff
			}
		} else {
			if fl.IsH() || (res&0xf) > 9 {
				res += 0x06
			}

			if fl.IsC() || res > 0x9f {
				res += 0x60
			}
		}

		fl.SetH(false)
		if res > 0xff {
			fl.SetC(true)
		}
		res = res & 0xff
		fl.SetZ(res == 0)
		return res
	})

	//9
	addFunction(f, "CPL", D8, func(fl *Flags, arg int) int {
		fl.SetN(true)
		fl.SetH(true)
		return ^arg & 0xff
	})

	//10
	addFunction(f, "SCF", D8, func(fl *Flags, arg int) int {
		fl.SetN(false)
		fl.SetH(false)
		fl.SetC(true)
		return arg
	})

	//11
	addFunction(f, "CCF", D8, func(fl *Flags, arg int) int {
		fl.SetN(false)
		fl.SetH(false)
		fl.SetC(!fl.IsC())
		return arg
	})

	//12
	addBiIntFunction(bf, "ADD", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		fl.SetZ(((arg1 + arg2) & 0xff) == 0)
		fl.SetN(false)
		fl.SetH((arg1&0x0f)+(arg2&0x0f) > 0x0f)
		fl.SetC(arg1+arg2 > 0xff)
		return (arg1 + arg2) & 0xff
	})

	//13
	addBiIntFunction(bf, "ADC", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		carry := 0
		if fl.IsC() {
			carry = 1
		}
		fl.SetZ(((arg1 + arg2 + carry) & 0xff) == 0)
		fl.SetN(false)
		fl.SetH((arg1&0x0f)+(arg2&0x0f)+carry > 0x0f)
		fl.SetC(arg1+arg2+carry > 0xff)
		return (arg1 + arg2 + carry) & 0xff
	})

	//14
	addBiIntFunction(bf, "SUB", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		fl.SetZ(((arg1 - arg2) & 0xff) == 0)
		fl.SetN(true)
		fl.SetH((0x0f & arg2) > (0x0f & arg1))
		fl.SetC(arg2 > arg2)
		return (arg1 - arg2) & 0xff
	})

	//15
	addBiIntFunction(bf, "SBC", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		carry := 0
		if fl.IsC() {
			carry = 1
		}
		res := arg1 - arg2 - carry
		fl.SetZ((res & 0xff) == 0)
		fl.SetN(true)
		fl.SetH(((arg1 ^ arg2 ^ (res & 0xff)) & (1 << 4)) != 0)
		fl.SetC(res < 0)
		return res & 0xff
	})

	//16
	addBiIntFunction(bf, "AND", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		res := arg1 & arg2
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(true)
		fl.SetC(false)
		return res
	})

	//17
	addBiIntFunction(bf, "AND", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		res := arg1 | arg2
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(false)
		fl.SetC(false)
		return res
	})

	//18
	addBiIntFunction(bf, "XOR", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		var result = (arg1 ^ arg2) & 0xff
		fl.SetZ(result == 0)
		fl.SetN(false)
		fl.SetH(false)
		fl.SetC(false)
		return result
	})

	//19
	addBiIntFunction(bf, "CP", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		fl.SetZ(((arg1 - arg2) & 0xff) == 0)
		fl.SetN(true)
		fl.SetH((0x0f & arg2) > (0x0f & arg1))
		fl.SetC(arg2 > arg1)
		return arg1
	})

	//20
	addFunction(f, "RLC", D8, func(fl *Flags, arg int) int {
		res := (arg << 1) & 0xff
		if (arg & (1 << 7)) != 0 {
			res |= 1
			fl.SetC(true)
		} else {
			fl.SetC(false)
		}
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(false)
		return res
	})

	//21
	addFunction(f, "RRC", D8, func(fl *Flags, arg int) int {
		res := arg >> 1
		if (arg & 1) == 1 {
			res |= (1 << 7)
			fl.SetC(true)
		} else {
			fl.SetC(false)
		}
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(false)
		return res
	})

	//22
	addFunction(f, "RL", D8, func(fl *Flags, arg int) int {
		res := (arg << 1) & 0xff
		if fl.IsC() {
			res |= 1
		}
		fl.SetC((arg & (1 << 7)) != 0)
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(false)
		return res
	})

	//23
	addFunction(f, "RR", D8, func(fl *Flags, arg int) int {
		res := arg >> 1
		if fl.IsC() {
			res |= 1 << 7
		}
		fl.SetC((arg & 1) != 0)
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(false)
		return res
	})

	//24
	addFunction(f, "SLA", D8, func(fl *Flags, arg int) int {
		res := (arg << 1) & 0xff
		fl.SetC((arg & (1 << 7)) != 0)
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(false)
		return res
	})

	//25
	addFunction(f, "SRA", D8, func(fl *Flags, arg int) int {
		res := (arg >> 1) | (arg & (1 << 7))
		fl.SetC((arg & 1) != 0)
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(false)
		return res
	})

	//26
	addFunction(f, "SWAP", D8, func(fl *Flags, arg int) int {
		upper := arg & 0xf0
		lower := arg & 0x0f
		res := (lower << 4) | (upper >> 4)
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(false)
		fl.SetC(false)
		return res
	})

	//27
	addFunction(f, "SRL", D8, func(fl *Flags, arg int) int {
		res := (arg >> 1)
		fl.SetC((arg & 1) != 0)
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH(false)
		return res
	})

	//28
	addBiIntFunction(bf, "BIT", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		bit := arg2
		fl.SetN(false)
		fl.SetH(true)
		if bit < 8 {
			fl.SetZ(!GetBit(arg1, arg2))
		}
		return arg1
	})

	//29
	addBiIntFunction(bf, "RES", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		return ClearBit(arg1, arg2)
	})

	//30
	addBiIntFunction(bf, "SET", D8, D8, func(fl *Flags, arg1 int, arg2 int) int {
		return SetBit(arg1, arg2)
	})

	return Alu{
		funcs:   f,
		biFuncs: bf,
	}
}

func addFunction(m map[AluFunctionKey]func(fl *Flags, arg int) int, name string, dtType int, fn func(fl *Flags, arg int) int) {
	m[NewAluFunctionKey(name, dtType)] = fn
}

func addBiIntFunction(m map[AluFunctionKey]func(fl *Flags, arg1 int, arg2 int) int, name string, dtType1 int, dtType2 int, fn func(fl *Flags, arg1 int, arg2 int) int) {
	m[NewAluBiIntFunctionKey(name, dtType1, dtType2)] = fn
}

func (a *Alu) GetFunction(name string, dtType int) IntRegistryFunc {
	return a.funcs[NewAluFunctionKey(name, dtType)]
}