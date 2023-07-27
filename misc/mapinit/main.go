package main

import "fmt"

type A struct {
	B map[string]string
	C map[string]string
}

func main() {
	a := A{}
	fmt.Println(a)
	// a.B["a"] = "b" // panic: assignment to entry in nil map
	aa := A{B: map[string]string{}}
	fmt.Println(aa)
	aa.B["a"] = "b"
}
