package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type A struct {
	Time time.Time
}

func main() {
	fmt.Println(time.Now().Format(time.RFC3339))
	now := time.Now()
	now_str := now.UTC().Format(time.RFC3339Nano)
	fmt.Println(now_str)
	a := &A{Time: now.UTC()}
	data, _ := json.Marshal(a)
	fmt.Println(string(data))

	data = []byte("{\"Time\":\"2022-10-27T09:22:03.167666Z\"}")
	a = &A{}
	json.Unmarshal(data, a)
	fmt.Println(a.Time.Format(time.RFC3339Nano))

	time1 := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	time2 := time.Date(2023, 1, 1, 9, 0, 0, 0, time.UTC)
	fmt.Println("time diff: ", time2.Sub(time1).Hours())
}
