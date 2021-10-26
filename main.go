package main

import "fmt"

// import (
// 	"github.com/theleao/gamebongo/cpu"
// )

type user struct {
	name string
}

var globalList []string

func main() {
	//cpu.LittleTest()

	globalList = append(globalList, "zero")
	Start()

	for _, s := range globalList {
		fmt.Println(s)
	}
}

func Start() {
	globalList = append(globalList, "um")
}

func TestFuncComparison() {

}