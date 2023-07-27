package main

import (
	"encoding/json"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

type Person struct {
	Name  string
	Age   int
	Email string
}

func main() {
	// 将结构体序列化为 MessagePack
	data, err := msgpack.Marshal(&Person{Name: "Alice", Age: 30, Email: "alice@example.com"})
	if err != nil {
		fmt.Println(err)
	}

	data_json, err := json.Marshal(&Person{Name: "Alice", Age: 30, Email: "alice@example.com"})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("msgpack data len:", len(data))
	fmt.Println("json data len:", len(data_json))

	// 将 MessagePack 反序列化为结构体
	var person Person
	err = msgpack.Unmarshal(data, &person)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(person)

}
