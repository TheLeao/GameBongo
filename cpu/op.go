package cpu

type Op interface {
	ReadsMemory() bool
	WritesMemory() bool
	CausesOemBug(reg Registers, opCntxt int) (bool, int)
}
