package main

import (
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/natefinch/lumberjack.v2"
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

var (
	server     string
	messageNum = 7000
	// messageNum = 7
	id              string
	httpServer      string
	interval        int
	kcpTimeout      int
	data            int
	parity          int
	nodelay         string
	nodelayInt      int
	nodelayInterval int
	nodelayResend   int
	nodelayNc       int
)

func httpSend() {
	dummyData := StringWithCharset(900, charset)
	format := `{"counter":%d, "id":"%s", "data":"%s"}`

	counter := 0

	client := &http.Client{}

	for {

		batchCounter := 0
		PBatchCounter := &batchCounter

		start := time.Now()

		for i := 0; i < messageNum; i++ {
			pI := &i

			func() {
				data := fmt.Sprintf(format, counter, id, dummyData)
				req, err := http.NewRequest("POST", httpServer, bytes.NewBuffer([]byte(data)))
				if err != nil {
					// fmt.Println("Error creating request: ", err)
					logrus.WithFields(logrus.Fields{
						"err": err,
					}).Info("http creating request")
					return
				}

				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					// fmt.Println("Error making request: ", err)
					(*pI)--
					logrus.WithFields(logrus.Fields{
						"err": err,
					}).Info("http request")
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					(*pI)--
					logrus.WithFields(logrus.Fields{
						"err": err,
					}).Info("http request")
				}

				_, _ = ioutil.ReadAll(resp.Body)

				(*PBatchCounter)++
				counter++

			}()
		}

		logrus.WithFields(logrus.Fields{
			"num": batchCounter,
		}).Info("http write")

		timeuse := time.Since(start)
		wait := time.Duration(interval) * time.Second
		if timeuse < wait {
			sleep := wait - timeuse
			time.Sleep(sleep)
		}
	}
}

func initLog() {

	// log := logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	// logrus.SetFormatter(&logrus.JSONFormatter{})

	// Keep log files to maximum size 5 MB
	logrus.SetOutput(io.MultiWriter(
		&lumberjack.Logger{
			Filename:   "my.log",
			MaxSize:    5, // megabytes
			MaxBackups: 5,
			MaxAge:     28, //days
		},
		os.Stdout))
}

func main() {
	flag.StringVar(&server, "s", "", "server")
	flag.StringVar(&id, "id", "test", "")
	flag.StringVar(&httpServer, "http", "", "http server")
	flag.IntVar(&interval, "n", 300, "interval")
	flag.IntVar(&kcpTimeout, "kt", 3, "kcp timeout")
	flag.IntVar(&messageNum, "msg", 7000, "msg num")
	flag.IntVar(&data, "data", 10, "data num")
	flag.IntVar(&parity, "parity", 3, "parity num")
	flag.StringVar(&nodelay, "nodelay", "-1:-1:-1:-1", "parity num")

	flag.Parse()

	if nodelay != "" {
		parts := strings.Split(nodelay, ":")
		if len(parts) == 4 {
			nodelayInt, _ = strconv.Atoi(parts[0])
			nodelayInterval, _ = strconv.Atoi(parts[1])
			nodelayResend, _ = strconv.Atoi(parts[2])
			nodelayNc, _ = strconv.Atoi(parts[3])
			logrus.WithFields(logrus.Fields{
				"nodelayInt":      nodelayInt,
				"nodelayInterval": nodelayInterval,
				"nodelayResend":   nodelayResend,
				"nodelayNc":       nodelayNc,
			}).Info("nodelay set")
		} else {
			logrus.Error("parse nodelay err")
		}
	}

	initLog()

	if httpServer != "" {
		go httpSend()
	}

	dummyData := StringWithCharset(900, charset)
	format := `{"counter":%d, "id":"%s", "data":"%s"}`

	counter := 0

	key := pbkdf2.Key([]byte("abc"), []byte("abc salt"), 1024, 16, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	if sess, err := kcp.DialWithOptions(server, block, data, parity); err == nil {
		sess.SetNoDelay(nodelayInt, nodelayInterval, nodelayResend, nodelayNc)
		sess.SetWindowSize(10000, 10000)

		// fmt.Println("local: ", sess.LocalAddr())
		logrus.WithFields(logrus.Fields{
			"addr": sess.LocalAddr(),
		}).Info("local")
		defer sess.Close()
		for {

			var err error
			failed := false

			start := time.Now()

			batchCounter := 0

			for i := 0; i < messageNum; i++ {

				data := fmt.Sprintf(format, counter, id, dummyData)

				// buf := make([]byte, 10)
				//fmt.Println("sent:", data)
				sess.SetDeadline(time.Now().Add(time.Duration(kcpTimeout) * time.Second))
				if _, err = sess.Write([]byte(data)); err == nil {
					// if n, err := sess.Read(buf); err == nil {
					// 	// fmt.Println("recv:", string(buf))
					// 	if string(buf[:n]) != "OK" {
					// 		logrus.WithFields(logrus.Fields{
					// 			"buf": string(buf[:n]),
					// 		}).Info("read not OK")
					// 	}
					// } else {
					// 	// fmt.Println("read err: ", err)
					// 	logrus.WithFields(logrus.Fields{
					// 		"err": err,
					// 	}).Info("read")
					// 	failed = true

					// }
				} else {
					// fmt.Println("write err:", err)
					logrus.WithFields(logrus.Fields{
						"err": err,
					}).Info("write")
					failed = true
				}

				if failed {
					sess.Close()

					for {
						sess, err = kcp.DialWithOptions(server, block, 10, 3)
						sess.SetNoDelay(nodelayInt, nodelayInterval, nodelayResend, nodelayNc)
						sess.SetWindowSize(10000, 10000)
						if err != nil {
							// fmt.Println("redial err: ", err)
							logrus.WithFields(logrus.Fields{
								"err": err,
							}).Info("redial")
							time.Sleep(5 * time.Second)
						} else {
							// fmt.Println("redial success local: ", sess.LocalAddr())
							logrus.WithFields(logrus.Fields{
								"addr": sess.LocalAddr(),
							}).Info("redial success")
							failed = false

							break
						}
					}

					i--
					counter--
				}

				counter += 1
				batchCounter += 1
			}

			logrus.WithFields(logrus.Fields{
				"num": batchCounter,
			}).Info("write")

			timeuse := time.Since(start)
			wait := time.Duration(interval) * time.Second
			if timeuse < wait {
				sleep := wait - timeuse
				time.Sleep(sleep)
			}
		}
	} else {
		panic(err)
	}
}
