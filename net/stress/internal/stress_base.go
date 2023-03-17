package stress

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	counter_ok  uint64 = 0
	counter_err uint64 = 0
)

type StressCase interface {
	Do() error
}

func doStress(Case StressCase, num int, interval time.Duration) {
	for i := 0; i < num; i++ {
		err := Case.Do()
		if err != nil {
			atomic.AddUint64(&counter_err, 1)
			log.Printf("failed: %s\n", err)
		}
		atomic.AddUint64(&counter_ok, 1)
		time.Sleep(interval)
	}
}

func ShowCounter() {
	last_ok_value := counter_ok
	last_err_value := counter_err
	t := time.NewTicker(60 * time.Second)

	for {
		<-t.C
		cok := atomic.LoadUint64(&counter_ok)
		cerr := atomic.LoadUint64(&counter_err)
		log.Printf("sent: %d/%d\n", cerr-last_err_value, cok-last_ok_value)
		last_ok_value = cok
		last_err_value = cerr
	}
}

func DoStress(Cases []StressCase, num int, rate int) error {
	go ShowCounter()

	wg := sync.WaitGroup{}
	num_each := num / len(Cases)
	rate_each := rate / len(Cases)
	sleep_time := time.Duration(1) * time.Second / time.Duration(rate_each)
	log.Println("sleep_time:", sleep_time)
	for _, Case := range Cases {
		wg.Add(1)
		go func(Case StressCase) {
			doStress(Case, num_each, sleep_time)
			wg.Done()
		}(Case)
	}
	wg.Wait()
	return nil
}
