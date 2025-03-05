package main

import (
	"context"
	"fmt"
	"time"

	"github.com/juju/ratelimit"
	"golang.org/x/time/rate"
)

func gorate() {
	lim := rate.NewLimiter(
		rate.Limit(1000),
		10000,
	)

	fmt.Println("available: ", lim.Tokens(), " rate: ", lim.Limit())
	lim.AllowN(time.Now(), 9000)
	fmt.Println("available: ", lim.Tokens())
	start := time.Now()
	// r := lim.ReserveN(time.Now(), 3000)
	// fmt.Println("time: ", r.Delay())
	// time.Sleep(r.Delay())
	lim.WaitN(context.Background(), 3000)
	fmt.Printf("time used: %v\n", time.Since(start))

}

func juju() {
	lim := ratelimit.NewBucket(10*time.Second/10000, 10000)

	fmt.Println("available: ", lim.Available(), " rate: ", lim.Rate())
	lim.TakeAvailable(9000)
	fmt.Println("available: ", lim.Available())

	start := time.Now()
	fmt.Println("wait for")
	// d := lim.Take(3000)
	// fmt.Println("time: ", d)
	lim.Wait(3000)
	fmt.Printf("time used: %v\n", time.Since(start))
}

func main() {
	gorate()
	juju()
}
