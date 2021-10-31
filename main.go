package main

import "fmt"

// import (
// 	"github.com/theleao/goingboy/cpu"
// )

type user struct {
	name string
}

var globalList []string

func main() {

	globalList = append(globalList, "zero")
	Start()

	for _, s := range globalList {
		fmt.Println(s)
	}
}

func Start() {
	globalList = append(globalList, "um")
}