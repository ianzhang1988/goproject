package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var bytes uint64

func Handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("handle ...")
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	atomic.AddUint64(&bytes, uint64(len(data)))

	dataStr := string(data)
	dataLine := strings.Split(dataStr, "\n")

	dataList := []string{}

	for _, l := range dataLine {
		if l == "" {
			continue
		}
		dataList = append(dataList, l)
	}

	// time.Sleep(10 * time.Second)
}

func Request() {

	string_line := ""
	for i := 0; i < 512; i++ {
		string_line += "A"
	}

	string_data := ""
	for i := 0; i < 1000; i++ {
		string_data += string_line
		string_data += "\n"
	}

	fmt.Println("data len:", len(string_data))

	client := &http.Client{
		Timeout: 1 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
		},
	}

	repeat := 5000
	begin := time.Now()
	for i := 0; i < repeat; i++ {
		data := strings.NewReader(string_data)
		resp, err := client.Post("http://127.0.0.1:8005/test", "", data)
		if err != nil {
			fmt.Printf("post err: %s\n", err)
		}

		io.Copy(ioutil.Discard, resp.Body) // read the response body
		resp.Body.Close()                  // close the response body
	}

	time_used := time.Since(begin).Seconds()

	fmt.Printf("time: %f, speed: %f\n", time_used, float64(repeat*len(string_data))/1024.0/1024.0/time_used)
	fmt.Printf("send: %d, receive: %d\n", repeat*len(string_data), bytes)
}

func main() {
	new_mux := http.NewServeMux()
	new_mux.HandleFunc("/test", Handler)

	server := http.Server{
		Addr:        ":8005",
		ReadTimeout: time.Second * 5,
		Handler:     new_mux,
	}

	go Request()

	err := server.ListenAndServe()
	panic(err)
}
