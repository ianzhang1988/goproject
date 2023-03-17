package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	token := make([]byte, 4)
	rand.Seed(time.Now().UnixNano())
	rand.Read(token)
	fmt.Println(token)
	fmt.Println(rand.Uint32())
}
