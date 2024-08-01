package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan string, 1000)
	ch_max := make(chan float64, 100000)

	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {

			timemax := 0.0

			for j := 0; j < 1000; j++ {
				start := time.Now()
				ch <- "test"
				timeused := time.Since(start)
				if timeused.Seconds() > timemax {
					timemax = timeused.Seconds()
				}
			}

			ch_max <- timemax

			wg.Done()
		}()
	}

	wg2 := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func() {
			start := time.Now()
			counter := 0
			for t := range ch {
				_ = t
				counter += 1
			}
			fmt.Println(float32(counter) / float32(time.Since(start).Seconds()))
			wg2.Done()
		}()
	}

	wg.Wait()
	close(ch)

	wg2.Wait()

	close(ch_max)

	time_max_10 := [10]float64{}

	for m := range ch_max {

		for idx := range time_max_10 {
			if m > time_max_10[idx] {
				time_max_10[idx] = m
			}
		}
	}

	fmt.Printf("channel write max wait time: %v\n", time_max_10)
	// 100 goroutine [0.018830945 0.018830945 0.018830945 0.018830945 0.018830945 0.018830945 0.018830945 0.018830945 0.018830945 0.018830945]
	// 100000 goroutine [1.565394344 1.565394344 1.565394344 1.565394344 1.565394344 1.565394344 1.565394344 1.565394344 1.565394344 1.565394344]
}
