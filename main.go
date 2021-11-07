package main

import "fmt"

// import (
// 	"github.com/theleao/goingboy/cpu"
// )

type user struct {
	name string
}

type IntQueue [] int

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

func (r Ram) SetByteNoPointer(addr int, value int) {	
	r.space[addr-r.offset] = value
}

var globalList []string

func main() {

	globalList = append(globalList, "zero")
	Start()

	xxx := IntQueue{1,2,3,4,5}
	xxx.enqueue(6)

	
	r := Ram{space: make([]int, 10)}
	fmt.Println(r.space)
	r.SetByte(5, 999)
	fmt.Println(r.space)
	r.SetByteNoPointer(4, 888)

}

func Start() {
	globalList = append(globalList, "um")
}