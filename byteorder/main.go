package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	////////////////////////////////// no need to explain anything
	b := []byte{2, 3, 5, 7, 11, 13} /// within this comment block.
	for _, e := range b {           //
		fmt.Printf("%d ", e) //
	} //
	fmt.Printf("\n") //
	//////////////////////////////
	num := binary.LittleEndian.Uint32(b) /// <<< Why this results in
	/// 117768962 is the question.
	/// 2*256^0 + 3*256^1 + 5*256^2 + 7*256^3 = 117768962
	fmt.Printf("customNum=%d\n", int(num))

	bb := make([]byte, 4)
	binary.LittleEndian.PutUint32(bb, num)
	for _, e := range bb { //
		fmt.Printf("%d ", e) //
	}
	fmt.Printf("\n") //

	binary.BigEndian.PutUint32(bb, num)
	for _, e := range bb { //
		fmt.Printf("%d ", e) //
	}
	fmt.Printf("\n") //
}
