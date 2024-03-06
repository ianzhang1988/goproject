// practice for https://go.dev/blog/pipelines

package main

import (
	"fmt"
	"sync"
	"time"
)

func gen(num int) chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < num; i++ {
			fmt.Println("gen:", i)
			out <- i
			time.Sleep(10 * time.Millisecond)
		}
		close(out)
	}()
	return out
}

func double(in chan int) chan int {
	out := make(chan int)
	go func() {
		for num := range in {
			out <- num * 2
		}
		close(out)
	}()
	return out
}

func merge(ins ...(chan int)) chan int {
	wg := sync.WaitGroup{}

	out := make(chan int)

	wg.Add(len(ins))
	for _, in := range ins {
		go func(in chan int) {
			for num := range in {
				out <- num
			}
			wg.Done()
		}(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	g := gen(10)

	d1 := double(g)
	d2 := double(g)

	for num := range merge(d1, d2) {
		fmt.Println(num)
	}
}
