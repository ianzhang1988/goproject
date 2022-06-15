package main

import "fmt"

// Sets the bit at pos in the integer n.
func setBit(n int, pos uint) int {
	n |= (1 << pos)
	return n
}

// Clears the bit at pos in n.
func clearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

func main() {
	fmt.Println(setBit(1, 2))   // 5
	fmt.Println(clearBit(3, 1)) // 1
}
