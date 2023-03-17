package main

import (
	"fmt"

	"github.com/armon/go-radix"
)

func main() {
	r := radix.New()
	r.Insert("ct", "ct")
	r.Insert("cmnet", 2)
	r.Insert("cnc", 2)
	r.Insert("ct/bj", 1)
	r.Insert("cmnet/bj", 2)
	r.Insert("cnc/bj", 2)
	r.Insert("ct/sh", 1)
	r.Insert("cmnet/sh", 2)
	r.Insert("cnc/sh", 2)
	r.Insert("ct/bj/1.0", 2)
	r.Insert("ct/bj/2.0", 2)
	r.Insert("ct/sh/1.0", 2)
	r.Insert("ct/sh/2.0", 2)

	// Find the longest prefix match
	m, v, ok := r.LongestPrefix("foozip")
	if ok {
		fmt.Println("key:", m)
		fmt.Println("value:", v)
	}

	v, ok = r.Get("ct")
	if ok {
		fmt.Println("value:", v)
	}

	// 深度优先遍历
	r.WalkPrefix("ct/", func(key string, v interface{}) bool {
		fmt.Println("key:", key)
		fmt.Println("value:", v)
		return false
	})
}
