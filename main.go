package main

import (
	"fmt"

	"github.com/theleao/goingboy/core"
)

// import (
// 	"github.com/theleao/goingboy/cpu"
// )

type user struct {
	name string
}

type IntQueue []int

func (q IntQueue) enqueue(v int) {
	q = append(q, v)
}

func (q IntQueue) dequeue() {

}

type Ram struct {
	space  []int
	length int
	offset int
}

func (r Ram) Accepts(addr int) bool {
	return addr >= r.offset && addr < r.offset+r.length
}

func (r Ram) SetByte(addr int, value int) {
	ptr := &r
	ptr.space[addr-r.offset] = value
}

func (r Ram) GetByte(addr int) int {
	return 999
}

func (r Ram) SetByteNoPointer(addr int, value int) {
	r.space[addr-r.offset] = value
}

var globalList []string

func main() {
	r := Ram{}
	h := NewHdma(r)
	h.transfInProgress = true
	h.SetByte(0,0)

	fmt.Println(h.transfInProgress)
}

func Start() {
	globalList = append(globalList, "um")
}

func (r Ram) TryPointer() {	
	r.length = 22
}

func NewHdma(addr core.AddressSpace) Hdma {
	return Hdma{
		hdma1234: Ram{
			length: 4,
			offset: 5,
		},
		addrSpace: addr,
	}
}

//interface

func (h *Hdma) Accepts(addr int) bool {
	return true
}

func (h *Hdma) SetByte(addr int, value int) {
	h.transfInProgress = false
}

func (h *Hdma) GetByte(addr int) int {
	panic("HDMA illegal argument on GetByte")
}

type Hdma struct {
	addrSpace        core.AddressSpace
	hdma1234         Ram
	mode             int
	transfInProgress bool
	hBlankTransfer   bool
	lcdEnabled       bool
	length           int
	src              int
	dst              int
	tick             int
}