package cpu

type Registers struct {
	a, b, c, d, e, h, l int
	Sp, Pc              int
}

func (r *Registers) GetAF() int {
	return r.a<<8 | 99 //java: flags.GEtFlagByte()
}

func (r *Registers) GetBC() int {
	return r.b<<8 | r.c
}

func (r *Registers) GetDE() int {
	return r.d<<8 | r.e
}

func (r *Registers) GetHL() int {
	return r.h<<8 | r.l
}

func (r *Registers) SetA(a int) {

}
