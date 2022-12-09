package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle ...")
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	// time.Sleep(10 * time.Second)
}

func Request() {
	for i := 0; i < 10000; i++ {
		go func() {
			rp, wp := io.Pipe()
			defer rp.Close()
			go func() {
				defer wp.Close()
				for i := 0; i < 100; i++ {
					_, err := wp.Write([]byte("a"))
					if err != nil {
						// fmt.Println("write pip err:", err)
						break
					}
					time.Sleep(time.Second)
				}
			}()
			req, err := http.NewRequest("POST", "http://127.0.0.1:8003/test", rp)
			if err != nil {
				fmt.Println("make request err:", err)
			}
			_, err = http.DefaultClient.Do(req)
			fmt.Println("request")
			if err != nil {
				fmt.Println("request err:", err)
			}
		}()
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	new_mux := http.NewServeMux()
	new_mux.HandleFunc("/test", Handler)

	server := http.Server{
		Addr:        ":8003",
		ReadTimeout: time.Second * 5,
		Handler:     new_mux,
	}

	go Request()

	err := server.ListenAndServe()
	panic(err)
}
