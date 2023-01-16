package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	cancel()
	defer cancel()
	select {
	case <-ctx.Done():
		fmt.Println("timeout")
	}
	cancel()
	fmt.Println("exit ...")
	time.Sleep(3 * time.Second)
}
