package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"time"
)

func getRealSizeOf(v interface{}) (int, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}

func main() {
	testmap := make(map[[22]byte]uint8, 20000000)
	size, err := getRealSizeOf(&testmap)
	if err != nil {
		println(err)
		return
	}
	println("size", size)

	for i := 0; i < 2000000; i++ {
		var ipPort_key [22]byte
		key := ipPort_key[18:22]
		binary.LittleEndian.PutUint32(key, uint32(i))

		testmap[ipPort_key] = uint8(1)
	}

	size, err = getRealSizeOf(&testmap)
	if err != nil {
		println(err)
		return
	}

	println("size", size)

	time.Sleep(100 * time.Second)
}
