package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 解析域名
	addrs, err := net.LookupHost("www.baidu.com")
	if err != nil {
		fmt.Println("LookupHost error:", err)
		os.Exit(1)
	}

	// 打印解析结果
	for _, addr := range addrs {
		fmt.Println(addr)
	}
}
