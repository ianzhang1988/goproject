package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var count uint64 = 10 // 初始化计数器值为 10

	// 减1操作
	atomic.AddUint64(&count, ^uint64(0))

	fmt.Println("New count:", count) // 输出: New count: 9
}
