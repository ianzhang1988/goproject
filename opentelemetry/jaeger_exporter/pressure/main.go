package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func request(c *http.Client, input_url string) error {
	req, err := http.NewRequest("GET", input_url, nil)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("http code: %d", resp.StatusCode)
	}

	io.Copy(ioutil.Discard, resp.Body)

	return err
}

func request_thread(urlqueue chan string) {
	client := &http.Client{}
	for {

		url, ok := <-urlqueue
		if !ok {
			break
		}
		err := request(client, url)
		if err != nil {
			fmt.Printf("err: %s\n", err)
		}
	}

	fmt.Println("exit thread")
}

func main() {
	thread_num := 10
	urlqueue := make(chan string, thread_num*10)
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = thread_num

	for i := 0; i < thread_num; i++ {
		go request_thread(urlqueue)
	}

	for i := 0; i < 3000000; i++ {
		if i%1000 == 0 {
			fmt.Printf("\r%d", i)
		}
		if i%100 == 0 {
			time.Sleep(50 * time.Millisecond)
		}
		urlqueue <- "http://127.0.0.1:7777/hello"
	}

	close(urlqueue)
	fmt.Println("exit")
	time.Sleep(5 * time.Second)
}
