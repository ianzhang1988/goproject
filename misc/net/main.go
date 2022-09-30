package main

import (
	"bytes"
	"fmt"
	"net"
)

func main() {
	ip1 := net.ParseIP("127.0.0.1")
	ip2 := net.ParseIP("127.0.0.1")

	fmt.Println("ip1 == ip2:", ip1.Equal(ip2))
	bIP1 := []byte(ip1)
	bIP2 := []byte(ip2)
	fmt.Println("bIP1 == bIP2:", bytes.Compare(bIP1, bIP2) == 0)

	fmt.Println("bIP1 = ", net.IP(bIP1).String())
	fmt.Println("bIP1 = 127.0.0.1", net.IP(bIP1).String() == "127.0.0.1")
}
