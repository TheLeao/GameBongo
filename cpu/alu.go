package cpu

type Alu struct {
	funcs   map[AluFunctionKey]func(f *Flags, arg int) int
	biFuncs map[AluFunctionKey]func(f *Flags, arg1 int, arg2 int) int
}

// type IntRegistryFunc interface {
// 	Apply(f Flags, arg int) int
// }

// type BiIntRegistryFunc interface {
// 	Apply(f Flags, arg1 int, arg2 int) int
// }

type AluFunctionKey struct {
	name    string
	dtType1 int
	dtType2 int
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
	f := make(map[AluFunctionKey]func(f *Flags, arg int) int)
	bf := make(map[AluFunctionKey]func(f *Flags, arg1 int, arg2 int) int)

	//1
	addFunction(f, "INC", D8, func(fl *Flags, arg int) int {
		res := (arg + 1) & 0xff
		fl.SetZ(res == 0)
		fl.SetN(false)
		fl.SetH((arg & 0x0f) == 0x0f)
		return res
	})

	//Todo - make the 'middle' ones - just to remove comp error
	//C# Alufunctions.cs

	//5
	addBiIntFunction(bf, "ADD", D16, D16, func(f *Flags, arg1 int, arg2 int) int {
		f.SetN(false)
		f.SetH((arg1&0x0fff)+(arg2+0x0fff) > 0x0fff)
		f.SetC(arg1+arg2 > 0xffff)
		return (arg1 + arg2) & 0xffff
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
