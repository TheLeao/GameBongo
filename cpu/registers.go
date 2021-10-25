package cpu

import (
	"fmt"
)

type Registers struct {
	A, B, C, D, E, H, L int
	Sp, Pc              int
	Flags               Flags
}

func (r *Registers) getAF() int {
	return r.A<<8 | 99 //java: flags.GEtFlagByte()
}

func (r *Registers) setAF(af int) {
	r.A = GetMsb(af)
	r.Flags.SetFlagsByte(GetLsb(af))
}

func (r *Registers) getBC() int {
	return r.B<<8 | r.C
}

func (r *Registers) setBC(bc int) {
	r.B = GetMsb(bc)
	r.C = GetLsb(bc)
}

func (r *Registers) getDE() int {
	return r.D<<8 | r.E
}

func (r *Registers) setDE(de int) {
	r.D = GetMsb(de)
	r.E = GetLsb(de)
}

func (r *Registers) getHL() int {
	return r.H<<8 | r.L
}

func (r *Registers) setHL(hl int) {
	r.H = GetMsb(hl)
	r.L = GetLsb(hl)
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

func (r *Registers) GetString() string {
	return fmt.Sprintf("AF=%04x, BC=%04x, DE=%04x, HL=%04x, SP=%04x, PC=%04x, %s", r.getAF(), r.getBC(), r.getDE(), r.getHL(), r.Sp, r.Pc, r.Flags.GetString())
}
