package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/armon/go-radix"
)

type Box struct {
	SN       string
	Isp      string
	Province string
	City     string
}

func main() {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("----- begin ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	r := radix.New()

	fmt.Println("start insert")
	ispList := []string{"cmnet", "ct", "cnc"}
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		b := &Box{
			SN:       fmt.Sprintf("my_sn_%d", i),
			Isp:      ispList[i%len(ispList)],
			Province: "shangdong",
			City:     "weifang",
		}

		r.Insert(fmt.Sprintf("%s/%s/%s/%s", b.Isp, b.Province, b.City, b.SN), &b)
	}

	v, ok := r.Get("ct/shangdong/weifang/my_sn_7")
	if ok {
		fmt.Println("value:", v)
	}

	fmt.Println("time used:", time.Since(start))

	runtime.ReadMemStats(&m)
	fmt.Println("----- after ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Println("----- after GC ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	counter := 0
	start = time.Now()
	r.WalkPrefix("ct/", func(key string, v interface{}) bool {
		// fmt.Println("key:", key)
		// fmt.Println("value:", v)
		counter += 1
		return false
	})
	fmt.Println("walk time used:", time.Since(start), "counter:", counter)
}
