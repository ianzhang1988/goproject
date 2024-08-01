package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer func(i int) {
				fmt.Println("end: ", i)
			}(i)

			fmt.Println("start: ", i)
			time.Sleep(5 * time.Second)

		}(i)
	}

	wg.Wait()
}
