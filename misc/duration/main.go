package main

import (
	"fmt"
	"time"
)

func main() {
	interval := time.Second / 2
	fmt.Println(interval)
	interval = time.Second / 0.5
	fmt.Println(interval)
	num := 3
	interval = time.Second / time.Duration(num)
	fmt.Println(interval)
}
