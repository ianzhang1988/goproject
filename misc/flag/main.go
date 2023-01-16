package main

import (
	"flag"
	"fmt"
)

var (
	Bool = flag.Bool("B", false, "bool test")
)

func main() {
	flag.Parse()

	fmt.Println("Bool:", *Bool)
}
