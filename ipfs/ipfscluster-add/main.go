package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type IpfsCluster struct {
	Name string `json:"name"`
	Cid  string `json:"cid"`
	Size int    `json:"size"`
}

func ipfsAdd(client *http.Client, data string) string {
	reader := strings.NewReader(data)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test")
	io.Copy(part, reader)
	writer.Close()

	r, _ := http.NewRequest("POST", "http://127.0.0.1:9094/add?local=true&stream-channels=false", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(r)
	if err == nil {
		data, _ := io.ReadAll(resp.Body)
		ret := &[]IpfsCluster{}
		err := json.Unmarshal(data, ret)
		if err != nil {
			fmt.Println("json failed: ", err.Error(), " data: ", string(data))
			return ""
		}

		return (*ret)[0].Cid
	} else {
		fmt.Println("err: ", err.Error())
	}

	return ""
}

func ipfsAddThread(data chan string, ret chan string) {
	client := &http.Client{}
L:
	for {
		select {
		case fileData, more := <-data:
			if more {
				cid := ipfsAdd(client, fileData)
				// fmt.Println("cid: ", cid)
				ret <- cid
			} else {
				// finish
				break L
			}
		}
	}
}

var (
	workerNum     = flag.Int("worker", 20, "worker for http send")
	fileNum       = flag.Int("num", 1, "file num for send")
	startNum      = flag.Int("start", 0, "start data at")
	outputCidFile = flag.String("o", "cid.txt", "output cid file")
)

func main() {
	flag.Parse()

	input := make(chan string, *workerNum*2)
	done := make(chan string)

	for i := 0; i < *workerNum; i++ {
		go ipfsAddThread(input, done)
	}

	go func() {
		file, err := os.Create(*outputCidFile)
		if err != nil {
			fmt.Println("open file failed: ", err.Error())
			return
		}
		defer file.Close()

		for {
			l, more := <-done
			if more {
				file.WriteString(fmt.Sprintf("%s\n", l))
			} else {
				break
			}
		}
	}()

	start := time.Now()

	for i := *startNum; i < *fileNum; i++ {
		input <- fmt.Sprintf("%032d", i)
		if i%100 == 0 {
			fmt.Printf("\r %d/%d", i, *fileNum)
		}
	}

	time.Sleep(10 * time.Second)

	close(input)
	close(done)

	fmt.Printf("time used: %v\n", time.Since(start).Seconds())
}
