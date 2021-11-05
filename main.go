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

var globalList []string

func main() {

	globalList = append(globalList, "zero")
	Start()

	var xxx IntQueue = []int{1,2,3,4,5}
	xxx.enqueue(6)

	for _, x := range xxx {
		fmt.Println(x)
	}
	
}

func Start() {
	globalList = append(globalList, "um")
}