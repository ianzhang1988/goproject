package main

import (
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

type MyMsg struct {
	Counter int    `json:"counter"`
	Id      string `json:"id"`
}

type LogstashMsg struct {
	Time time.Time `json:"@timestamp"`
	Msg  string    `json:"message"`
}

var (
	id       string
	interval int
)

func httpServer() {

	counterChan := make(chan MyMsg, 10000)

	go countMsg(counterChan, "http")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("http read err: ", err)
		}
		defer r.Body.Close()

		msg := MyMsg{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			// fmt.Printf("json err: %s\n", err)
			w.WriteHeader(http.StatusTeapot)
			return
		}

		counterChan <- msg
	})

	http.ListenAndServe(":9204", nil)
}

func handle(conn *kcp.UDPSession, out chan MyMsg) {
	defer conn.Close()
	buf := make([]byte, 4096)

	// counter := 0

	for {
		conn.SetDeadline(time.Now().Add(11 * time.Minute))

		// fmt.Println("d 1")
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}
		// fmt.Println("d 2")

		msg := MyMsg{}
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			fmt.Printf("json err: %s\n", err)
			continue
		}

		// fmt.Println("client counter: ", msg.Counter)

		// _, err = conn.Write([]byte(fmt.Sprintf("OK %d", msg.Counter)))
		// if err != nil {
		// 	fmt.Println(err)
		// 	break
		// }

		out <- msg
	}

	fmt.Println("handle close: ", conn.RemoteAddr().String())
}

func countMsg(input chan MyMsg, proto string) {
	start := time.Now()
	counter := 0

	lock := sync.Mutex{}

	t := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for _ = range t.C {
			start = time.Now()
			lock.Lock()
			fmt.Printf("%s time: %v, counter: %d\n", proto, start, counter)
			counter = 0
			lock.Unlock()
		}
	}()

	msgCounterExpect := 0

	for msg := range input {

		if msg.Counter != msgCounterExpect {
			fmt.Printf("%s counter miss, expect %d, actual: %d\n", proto, msgCounterExpect, msg.Counter)
			msgCounterExpect = msg.Counter + 1
		} else {
			msgCounterExpect += 1
		}

		lock.Lock()
		counter += 1
		lock.Unlock()
	}
}

func main() {

	flag.IntVar(&interval, "n", 300, "interval")
	flag.StringVar(&id, "id", "test", "")
	flag.Parse()

	// go httpServer()

	key := pbkdf2.Key([]byte("abc"), []byte("abc salt"), 1024, 16, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	counterChan := make(chan MyMsg, 10000)
	defer func() {
		close(counterChan)
	}()

	go countMsg(counterChan, "udp")

	if listener, err := kcp.ListenWithOptions("0.0.0.0:9203", block, 10, 3); err == nil {
		for {
			s, err := listener.AcceptKCP()
			s.SetNoDelay(1, 50, 2, 1)
			s.SetWindowSize(10000, 10000)
			if err != nil {
				panic(err)
			}
			fmt.Println("start conn: ", s.RemoteAddr().String())
			go handle(s, counterChan)
		}
	} else {
		panic(err)
	}

}
