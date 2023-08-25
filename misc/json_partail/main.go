package main

import (
	"encoding/json"
	"fmt"
)

// 似乎没有部分解析的功能，目前只想到用interface来处理

type PartialJSON struct {
	FieldA        string          `json:"fieldA"`
	RemainingJSON json.RawMessage `json:"pass"`
}

func main() {
	jsonData := []byte(`{
		"fieldA": "valueA",
		"pass":{
			"fieldB": "valueB",
			"fieldC": "valueC"
		}
	}`)

	var partial PartialJSON
	err := json.Unmarshal(jsonData, &partial)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("FieldA:", partial.FieldA)
	fmt.Println("RestOfFields:", string(partial.RemainingJSON))

	var j map[string]interface{}
	err = json.Unmarshal(jsonData, &j)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for k, v := range j {
		fmt.Println("k", k, "v", v)
	}

	if v, ok := j["fieldA"]; ok {
		if value, ok := v.(string); ok {
			fmt.Println("value", value)
		}

	}

	d, err := json.Marshal(j)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(d))

	j["add"] = map[string]string{
		"a": "1",
		"b": "2",
	}
	d, err = json.Marshal(j)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(d))
}
