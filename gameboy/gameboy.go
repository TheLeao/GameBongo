package gameboy

type AddressSpace interface {
	Accepts(addr int) bool
	SetByte(addr int, value int)
	GetByte(addr int) int
}

type Stone struct {
	name string
}

func Test() {

}