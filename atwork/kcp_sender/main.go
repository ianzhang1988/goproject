package main

import (
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
	"golang.org/x/crypto/pbkdf2"
)

var (
	server    string
	worker    int
	interval  int
	localPort int
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func addChan(ch chan []byte, data []byte, sync bool) bool {
	if sync {
		ch <- data
		return true
	} else {
		select {
		case ch <- data:
			return true
		default:
			return false
		}
	}
}

func readChan(ch chan []byte, data *[]byte, sync bool) bool {
	if sync {
		*data = <-ch
		return true
	} else {
		select {
		case *data = <-ch:
			return true
		default:
			return false
		}
	}
}

func genData(input chan []byte) {
	counter := 0
	dummyData := StringWithCharset(900, charset)
	format := `{"counter":%d, "id":"%s", "data":"%s"}`

	for {
		counter += 1
		data := fmt.Sprintf(format, counter, "id", dummyData)
		dataByte := []byte(data)
		addChan(input, dataByte, true)
	}

}

func client() {
	// dummyData := StringWithCharset(900, charset)
	// format := `{"counter":%d, "id":"%s", "data":"%s"}`

	key := pbkdf2.Key([]byte("abc"), []byte("abc salt"), 1024, 16, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	var sess *kcp.UDPSession
	var err error
	var conn *net.UDPConn

	input := make(chan []byte, 1000)
	retrans := make(chan []byte, worker)

	go genData(input)

	for {

		if localPort < 1 {
			sess, err = kcp.DialWithOptions(server, block, 10, 3)
			if err != nil {
				fmt.Println("dial err:", err)
				time.Sleep(1 * time.Second)
				continue
			}
		} else {
			localaddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", localPort))
			if err != nil {
				fmt.Println("dial err:", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// sess would not close conn, if conn is passed by argment
			if conn != nil {
				conn.Close()
			}
			conn, err = net.ListenUDP("udp", localaddr)
			if err != nil {
				fmt.Println("dial err:", err)
				time.Sleep(1 * time.Second)
				continue
			}

			sess, err = kcp.NewConn(server, block, 10, 3, conn)
			if err != nil {
				fmt.Println("dial err:", err)
				time.Sleep(1 * time.Second)
				continue
			}
		}

		fmt.Println("local addr: ", sess.LocalAddr())

		sess.SetNoDelay(1, 50, 2, 1)
		sess.SetWindowSize(100000, 100000)
		// Get a TCP connection

		// Setup client side of smux
		// conf := smux.DefaultConfig()
		session, err := smux.Client(sess, nil)
		if err != nil {
			panic(err)
		}

		var counter int32
		t := time.NewTicker(time.Duration(interval) * time.Second)
		go func() {
			for _ = range t.C {
				tmp := atomic.LoadInt32(&counter)
				fmt.Printf("time: %v, counter: %d\n", time.Now(), tmp)
				atomic.StoreInt32(&counter, 0)
			}
		}()

		wg := sync.WaitGroup{}

		for i := 0; i < worker; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for {
					// Open a new stream
					stream, err := session.OpenStream()
					if err != nil {
						fmt.Println("session.OpenStream: ", err)
						break
					}

					for {

						var data []byte
						if ok := readChan(retrans, &data, false); !ok {
							readChan(input, &data, true)
						}

						msg := map[string]interface{}{}
						json.Unmarshal(data, &msg)
						cnt := int(msg["counter"].(float64))

						// localCnt := atomic.LoadInt32(&counter)
						stream.SetDeadline(time.Now().Add(10 * time.Second))
						// Stream implements io.ReadWriteCloser
						// data := fmt.Sprintf(format, localCnt, "id", dummyData)

						_, err = stream.Write([]byte(data))
						if err != nil {
							fmt.Println("stream.Write: ", err)
							addChan(retrans, data, true)
							break
						}
						buf := make([]byte, 100)
						n, err := stream.Read(buf)
						if err != nil {
							fmt.Println("stream.Read: ", err)
							addChan(retrans, data, true)
							break
						}

						if string(buf[:n]) != fmt.Sprintf("%d", cnt) {
							fmt.Printf("counter miss expect:%s, actual:%s\n", fmt.Sprintf("%d", cnt), string(buf[:n]))
						}

						atomic.AddInt32(&counter, 1)

						time.Sleep(1 * time.Millisecond)
						// fmt.Println(string(buf[:n]))
					}

					fmt.Println("stream break")
					stream.Close()
				}

			}()
		}

		wg.Wait()
		t.Stop()
		session.Close()

		fmt.Println("session close: ", session.IsClosed())
	}
}

func main() {
	flag.StringVar(&server, "s", "127.0.0.1:9203", "server")
	flag.IntVar(&worker, "w", 10, "worker num")
	flag.IntVar(&interval, "n", 300, "show interval")
	flag.IntVar(&localPort, "p", 0, "local port")
	flag.Parse()
	client()
}
