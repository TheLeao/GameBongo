package cpu

type Registers struct {
	a, b, c, d, e, h, l int
	Sp, Pc              int
}

func (r *Registers) getAF() int {
	return r.a<<8 | 99 //java: flags.GEtFlagByte()
}

func (r *Registers) getBC() int {
	return r.b<<8 | r.c
}

func (r *Registers) getDE() int {
	return r.d<<8 | r.e
}

func (r *Registers) getHL() int {
	return r.h<<8 | r.l
}

func (r *Registers) setA(a int) {

}

func (r *Registers) incrementPC() {
	r.Pc += 1 & 0xffff
}

func (r *Registers) incrementSP() {
	r.Sp += 1 & 0xffff
}

func (r *Registers) decrementSP() {
	r.Sp -= 1 & 0xffff
}