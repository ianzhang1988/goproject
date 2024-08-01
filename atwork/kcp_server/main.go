package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	"sync/atomic"

	"github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
	"golang.org/x/crypto/pbkdf2"
)

type MyMsg struct {
	Counter int    `json:"counter"`
	Id      string `json:"id"`
}

func handle(conn *kcp.UDPSession) {
	// defer conn.Close() // smux would close conn

	// buf := make([]byte, 4096)

	// counter := 0

	// Setup server side of smux
	session, err := smux.Server(conn, nil)
	if err != nil {
		fmt.Println("smux.Server: ", err)
		return
	}

	defer session.Close()

	var counter int32

	t := time.NewTicker(time.Duration(5) * time.Second)
	go func() {
		for _ = range t.C {
			tmp := atomic.LoadInt32(&counter)
			fmt.Printf("time: %v, counter: %d\n", time.Now(), tmp)
			atomic.StoreInt32(&counter, 0)
		}
	}()

	defer t.Stop()

	go func() {
		ch := session.CloseChan()
		<-ch
		fmt.Println("session closed")
		session.Close()
	}()

	for {
		// Accept a stream
		stream, err := session.AcceptStream()
		if err != nil {
			fmt.Println("session.AcceptStream: ", err)
			break
		}

		go func() {
			defer stream.Close()

			for {
				stream.SetDeadline(time.Now().Add(3 * time.Second))
				// Listen for a message
				buf := make([]byte, 4096)
				n, err := stream.Read(buf)
				if err != nil {
					fmt.Println("stream.Read: ", err)
					break
				}

				msg := MyMsg{}
				err = json.Unmarshal(buf[:n], &msg)
				if err != nil {
					// fmt.Printf("json err: %s\n", err)
					continue
				}

				// fmt.Println(string(buf))
				_, err = stream.Write([]byte(fmt.Sprintf("%d", msg.Counter)))
				if err != nil {
					fmt.Println("stream.Write: ", err)
					break
				}

				atomic.AddInt32(&counter, 1)
			}
		}()

	}

	fmt.Println("handle closed!")
}

func server() {

	key := pbkdf2.Key([]byte("abc"), []byte("abc salt"), 1024, 16, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	if listener, err := kcp.ListenWithOptions("0.0.0.0:9203", block, 10, 3); err == nil {
		for {
			s, err := listener.AcceptKCP()
			s.SetNoDelay(1, 50, 2, 1)
			s.SetWindowSize(100000, 100000)
			// s.SetNoDelay(-1,-1,-1,-1)
			if err != nil {
				panic(err)
			}
			fmt.Println("start conn: ", s.RemoteAddr().String())
			go handle(s)
		}
	} else {
		panic(err)
	}
}

func main() {
	server()
}
