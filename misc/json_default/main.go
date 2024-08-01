package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type NetProbe struct {
	Id           string            `json:"id"`
	Type         string            `json:"type"`
	Meta         map[string]string `json:"meta,omitempty"`
	Statistics   bool              `json:"statistics,omitempty"`
	TargetIpPort bool              `json:"target_ipport,omitempty"`
	Retry        int               `json:"retry,omitempty"`
}

func main() {
	var probe = NetProbe{
		Statistics: true, // 将Statistics字段预设为true
		// 其他字段可以根据你的需要进行预设
	}

	data := []byte(`{"id":"123", "type":"probe", "meta": {"address": "localhost", "port": "80"}}`) // 这是一个没有 "statistics" 字段的JSON数据
	err := json.Unmarshal(data, &probe)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", probe) // 打印解析后的结构体，你会看到 Statistics 的值是 true

	data = []byte(`{"statistics": false,"id":"123", "type":"probe", "meta": {"address": "localhost", "port": "80"}}`)
	err = json.Unmarshal(data, &probe)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", probe) // 打印解析后的结构体，你会看到 Statistics 的值是 false

}
