package main

import (
	"context"
	"fmt"
	"time"
	// "golang.org/x/time/rate"
)

func Produce(ch chan uint64) {
	var counter uint64 = 0
	// lastTime := time.Now()
	// var lastCounter uint64 = 0
	for {
		ch <- counter
		counter += 1

		// du := time.Since(lastTime)
		// if du > 1*time.Second {
		// 	fmt.Printf("%d %v\r", counter-lastCounter, du)
		// 	lastCounter = counter
		// 	// lastTime = time.Now()
		// 	lastTime = lastTime.Add(du)
		// }
	}
}

func ConsumNoLimit(ch chan uint64, chout chan uint64, lim *Limiter) {
	// 2.3 mil
	for {
		chout <- <-ch
	}
}

func ConsumWithLimitWait(ch chan uint64, chout chan uint64, lim *Limiter) {
	num := 1
	// 100k ok 200k -> 800k
	for {
		lim.WaitN(context.Background(), num)
		// if err != nil { // never
		// 	fmt.Println(err)
		// 	continue
		// }
		for i := 0; i < num; i++ {
			chout <- <-ch
		}
	}
}

func ConsumWithLimitDelay(ch chan uint64, chout chan uint64, lim *Limiter) {
	// below 100000 is good, 200000 double
	num := 1
	for {
		n := lim.ReserveN(time.Now(), num)
		if !n.OK() {
			fmt.Println("Reserve failed")
			// time.Sleep(n.Delay())
			continue
		}

		time.Sleep(n.Delay())
		for i := 0; i < num; i++ {
			chout <- <-ch
		}
	}
}

func ConsumWithLimit(ch chan uint64, chout chan uint64, lim *Limiter) {

	for {
		if lim.Allow() {
			chout <- <-ch
			// _ = <-ch
		}
		// } else {
		// 	_ = lim.Wait(context.Background())
		// }
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
	// 也许把rate limiter实现在Produce才是对的。现在看到的limiter是线程安全的，但是并不适用并行限制。
	go Produce(ch)
	lim := NewLimiter(500000.0, 10000)
	for i := 0; i < 100; i++ {
		// go ConsumNoLimit(ch, ch2, lim)
		// go ConsumWithLimit(ch, ch2, lim)
		go ConsumWithLimitDelay(ch, ch2, lim)
		// go ConsumWithLimitWait(ch, ch2, lim)
	}
	Count(ch2)
}
