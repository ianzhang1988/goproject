package main

import (
	"fmt"
	"time"
	// . "golang.org/x/time/rate"
)

func Produce(ch chan uint64) {
	var counter uint64 = 0
	for {
		ch <- counter
		counter += 1
	}
}

func ConsumWithLimitDelay(ch chan uint64, chout chan uint64, lim *Limiter) {
	// below 100000 is good, 200000 double
	for {
		n := lim.Reserve()
		if !n.OK() {
			continue
		}

		time.Sleep(n.Delay())
		chout <- <-ch
	}
}

func Count(ch chan uint64) {
	var counter uint64 = 0

	lastTime := time.Now()

	for {
		_ = <-ch
		counter += 1
		du := time.Since(lastTime)
		if du > 1*time.Second {
			fmt.Printf("%d %v\n", counter, du)
			counter = 0
			lastTime = lastTime.Add(du)
		}
	}
}

func main() {
	ch := make(chan uint64, 1000)
	ch2 := make(chan uint64, 100)
	go Produce(ch)
	lim := NewLimiter(200000.0, 1000)
	for i := 0; i < 100; i++ {
		go ConsumWithLimitDelay(ch, ch2, lim)
	}
	Count(ch2)
}
