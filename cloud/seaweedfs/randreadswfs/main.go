package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/goseaweedfs"
)

var (
	idsFile string
	worker  int
)

var Sum uint64

func download(idChan chan string, sw *goseaweedfs.Seaweed) {

	for id := range idChan {

		_, err := sw.Download(id, url.Values{}, func(r io.Reader) error {
			data, err := io.ReadAll(r)
			if err != nil {
				fmt.Println("download read err:", err)
				return err
			}
			atomic.AddUint64(&Sum, uint64(len(data)))
			return nil
		})
		if err != nil {
			fmt.Println("download err:", err)
		}
	}

}

func main() {

	flag.StringVar(&idsFile, "f", "benchmark.file", "file ids")
	flag.IntVar(&worker, "n", 10, "worker num")

	flag.Parse()

	data, err := os.ReadFile(idsFile)
	idArr := strings.Split(string(data), "\n")
	rand.Shuffle(len(idArr), func(i, j int) {
		idArr[i], idArr[j] = idArr[j], idArr[i]
	})

	idChan := make(chan string, worker)

	client := http.Client{}

	sw, err := goseaweedfs.NewSeaweed("http://localhost:9333", []string{}, 1024*1024, &client)
	if err != nil {
		fmt.Println("seaweed err:", err)
		return
	}

	wg := sync.WaitGroup{}

	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {

			download(idChan, sw)
			wg.Done()
		}()

	}

	start := time.Now()
	for _, id := range idArr {
		idChan <- id
	}

	close(idChan)

	wg.Wait()

	timeUsed := time.Since(start)

	fmt.Println("time used:", timeUsed)
	fmt.Println("bytes read:", Sum)

	fmt.Printf("speed %f MiB/s", float64(Sum)/1024/1024/timeUsed.Seconds())

}
