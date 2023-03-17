package stress

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type HTTPStressCase struct {
	Data       []byte
	TargetIP   string
	TargetPort string
	Client     *http.Client
	Url        string
}

var http_counter uint64 = 0
var http_counter_err uint64 = 0

func NewHttpStressCase(data []byte, target_ip string, target_port string, url string, keep_alive bool, conn_pool bool) (*HTTPStressCase, error) {

	t := http.DefaultTransport.(*http.Transport).Clone()

	// log.Printf("keep_alive: %v, conn_pool %v", keep_alive, conn_pool)
	if !keep_alive {
		t.DisableKeepAlives = true
	}
	if !conn_pool {
		t.MaxIdleConnsPerHost = -1 // 关闭连接池
	}
	t.ResponseHeaderTimeout = 5 * time.Second
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: t,
	}

	Case := &HTTPStressCase{
		Data:       data,
		TargetIP:   target_ip,
		TargetPort: target_port,
		Url:        url,
		Client:     client,
	}

	return Case, nil
}

func (c *HTTPStressCase) Do() error {
	body := bytes.NewReader(c.Data)
	resp, err := c.Client.Post(fmt.Sprintf("http://%s:%s/%s", c.TargetIP, c.TargetPort, c.Url), "", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func doHttpstress(Case *HTTPStressCase, num int, interval time.Duration) {
	for i := 0; i < num; i++ {
		err := Case.Do()
		if err != nil {
			atomic.AddUint64(&http_counter_err, 1)
			log.Printf("http faild: %s\n", err)
		}
		atomic.AddUint64(&http_counter, 1)
		time.Sleep(interval)
	}
}

func ShowHttpCounter() {
	last_value := http_counter
	last_err_value := http_counter_err
	t := time.NewTicker(60 * time.Second)

	for {
		<-t.C
		c := atomic.LoadUint64(&http_counter)
		c_err := atomic.LoadUint64(&http_counter_err)
		log.Printf("sent: %d/%d\n", c_err-last_err_value, c-last_value)
		last_value = c
		last_err_value = c_err
	}
}

func HttpStress(Cases []*HTTPStressCase, num int, rate int) error {
	go ShowHttpCounter()

	wg := sync.WaitGroup{}
	num_each := num / len(Cases)
	rate_each := rate / len(Cases)
	sleep_time := time.Duration(1) * time.Second / time.Duration(rate_each)
	log.Println("sleep_time:", sleep_time)
	for _, Case := range Cases {
		wg.Add(1)
		go func(Case *HTTPStressCase) {
			doHttpstress(Case, num_each, sleep_time)
			wg.Done()
		}(Case)
	}
	wg.Wait()
	return nil
}
