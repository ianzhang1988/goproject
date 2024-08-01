package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	dialerFunc := func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, time.Second*10)
		if err != nil {
			return nil, err
		}
		// 打印ip地址
		fmt.Printf("Connected to %v\n", conn.RemoteAddr())
		return conn, nil
	}

	transport := &http.Transport{
		DialContext:       dialerFunc,
		ForceAttemptHTTP2: true,
	}
	client := &http.Client{
		Transport: transport,
	}

	// 测试获取ip地址
	_, _ = client.Get("https://www.baidu.com")
	fmt.Println("---------------")
	// 测试获取ip地址
	_, _ = client.Get("http://127.0.0.1:8080")
}
