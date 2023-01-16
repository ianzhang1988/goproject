package main

import "fmt"

func main() {
	a := []byte{'a', 'b', 'c'}
	b := []byte{'1', '2', '3'}
	t := a

	fmt.Println(string(a))
	fmt.Println(string(b))
	fmt.Printf("%x\n", a)
	fmt.Printf("%x\n", b)

	a = b
	b = t

	fmt.Println(string(a))
	fmt.Println(string(b))
	fmt.Printf("%x\n", a)
	fmt.Printf("%x\n", b)

}
