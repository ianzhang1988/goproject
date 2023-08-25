package main

import (
	"encoding/json"
	"fmt"
)

type OldJson struct {
	Sanity int32
	// Isp_name_utf8      [128]uint8
	// Province_name_utf8 [64]byte
}

type NewJson struct {
	Sanity             int32
	Isp_name_utf8      string
	Province_name_utf8 string
}

func main() {
	old := OldJson{Sanity: 1}
	data, err := json.Marshal(&old)
	if err != nil {
		fmt.Println("marshal err:", err)
		return
	}

	fmt.Println("json:", string(data))

	new := NewJson{}
	err = json.Unmarshal(data, &new)
	if err != nil {
		fmt.Println("unmarshal err:", err)
		return
	}

	fmt.Println("new:", new)
}
