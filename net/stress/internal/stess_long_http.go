package stress

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type LongHTTPStressCase struct {
	Data       []byte
	TargetIP   string
	TargetPort string
	Client     *http.Client
	Url        string
	HoldTime   time.Duration
}

func NewLongHttpStressCase(data []byte, target_ip string, target_port string, url string, hold_time int) (*LongHTTPStressCase, error) {

	t := http.DefaultTransport.(*http.Transport).Clone()
	// t.DisableKeepAlives = true // 使用长连接
	t.MaxIdleConnsPerHost = -1 // 关闭连接池
	t.ResponseHeaderTimeout = 300 * time.Second
	client := &http.Client{
		Timeout:   300 * time.Second,
		Transport: t,
	}

	Case := &LongHTTPStressCase{
		Data:       data,
		TargetIP:   target_ip,
		TargetPort: target_port,
		Url:        url,
		Client:     client,
		HoldTime:   time.Duration(hold_time) * time.Second,
	}

	return Case, nil
}

// func (c *LongHTTPStressCase) Do() error {
// 	pr, pw := io.Pipe()

// 	interval := c.HoldTime / 2 / time.Duration(len(c.Data))

// 	go func() {
// 		for _, b := range c.Data {
// 			pw.Write([]byte{b})
// 			time.Sleep(interval)
// 		}
// 		pw.Close()
// 	}()

// 	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/%s", c.TargetIP, c.TargetPort, c.Url), pr)
// 	if err != nil {
// 		return err
// 	}
// 	// req.Header.Set("Content-Type", contentType)

// 	// datatmp, err := httputil.DumpRequestOut(req, true)
// 	// if err != nil {
// 	// 	fmt.Printf("dump err %s\n", err)
// 	// }
// 	// fmt.Printf("dump data: %s\n", string(datatmp))

// 	resp, err := c.Client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	length := resp.ContentLength

// 	if length < 1 {
// 		return fmt.Errorf("content length is 0")
// 	}

// 	interval = c.HoldTime / 2 / time.Duration(length)
// 	data := make([]byte, 1)
// 	for i := 0; i < int(length); i++ {
// 		_, err := resp.Body.Read(data)
// 		if err == io.EOF {
// 			return nil
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		time.Sleep(interval)
// 	}

// 	return nil
// }

func (c *LongHTTPStressCase) Do() error {
	body := bytes.NewReader(c.Data)
	resp, err := c.Client.Post(fmt.Sprintf("http://%s:%s/%s", c.TargetIP, c.TargetPort, c.Url), "", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	time.Sleep(c.HoldTime)

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
