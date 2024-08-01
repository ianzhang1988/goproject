package main

import (
	"fmt"
	"net/url"
)

func main() {
	baseUrl, err := url.Parse("http://example.com/")

	if err != nil {
		fmt.Println("Error parsing URL", err)
		return
	}

	// 创建URL参数
	params := url.Values{}
	params.Add("param1", "value1")
	params.Add("param2", "value2")

	// 添加参数到URL
	baseUrl.RawQuery = params.Encode()

	// 发送HTTP Get请求
	fmt.Println(baseUrl.String())
}
