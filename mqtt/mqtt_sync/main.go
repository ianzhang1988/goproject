package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type message struct {
	topic   string
	payload []byte
}

var connected = func(host string) func(mqtt.Client) {
	var f mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Printf("Connected<%s>: %v\n", host, client.IsConnectionOpen())
	}
	return f
}

var workerNum = flag.Int("num", 10, "worker num")

var host = flag.String("host", "", "host")
var caFile = flag.String("ca", "", "ca file")
var password = flag.String("password", "", "The password (optional)")
var user = flag.String("user", "", "The User (optional)")

var dstHost = flag.String("dst_host", "", "dst host")
var dstCaFile = flag.String("dst_ca", "", "dst ca file")
var dstPassword = flag.String("dst_password", "", "The dst password (optional)")
var dstUser = flag.String("dst_user", "", "The dst User (optional)")

func NewTLSConfig(filename string) *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(filename)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	// cert, err := tls.LoadX509KeyPair("samplecerts/client-crt.pem", "samplecerts/client-key.pem")
	// if err != nil {
	// 	panic(err)
	// }

	// Just to print out the client certificate..
	// cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(cert.Leaf)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		// Certificates: []tls.Certificate{cert},
	}
}

func NewMQTTClient(host, user, password, caFile string, msgHandle *mqtt.MessageHandler) mqtt.Client {
	// url := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Println("conneting to:", host)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(host)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(*msgHandle)
	opts.SetOnConnectHandler(connected(host))
	opts.SetPingTimeout(10 * time.Second)
	opts.SetCleanSession(true)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetAutoReconnect(true)
	if caFile != "" {
		tlsconfig := NewTLSConfig(caFile)
		opts.SetTLSConfig(tlsconfig)
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("connect to ", host, " failed.")
		panic(token.Error())
	}

	return c
}

func main() {

	flag.Parse()

	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.WARN = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	msgsChannel := make(chan message, 10000)
	dstClients := make([]mqtt.Client, *workerNum)

	var srcMsgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())

		newMsg := message{topic: msg.Topic(), payload: msg.Payload()}

		msgsChannel <- newMsg
	}

	srcClient := NewMQTTClient(*host, *user, *password, *caFile, &srcMsgHandler)

	// sub all src topic
	if token := srcClient.Subscribe("#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println("sub error: ", token.Error())
		os.Exit(1)
	}

	filter := func(payload []byte) bool {
		// string(payload)
		num, err := strconv.Atoi(string(payload))
		if err == nil && num%2 != 0 {
			return false
		}

		return true
	}

	worker := func(c mqtt.Client) {
		for {
			msg := <-msgsChannel

			// filter
			if !filter(msg.payload) {
				continue
			}

			fmt.Println("pub ", msg.topic, " ", msg.payload)
			token := c.Publish(msg.topic, 0, false, msg.payload)
			token.Wait()
		}
	}

	for i := 0; i < *workerNum; i++ {
		var dummyMsgHandler mqtt.MessageHandler

		dstC := NewMQTTClient(*dstHost, *dstUser, *dstPassword, *dstCaFile, &dummyMsgHandler)
		dstClients[i] = dstC

		go worker(dstC)
	}

	exit := false
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听k8s信号 SIGTERM
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGTERM:
				fmt.Println("SIGTERM")
				exit = true
			default:
				fmt.Println("other signal:", s)
			}
		}
	}()

	// clean up
	for !exit {
		time.Sleep(1 * time.Second)
	}

	if token := srcClient.Unsubscribe("#"); token.Wait() && token.Error() != nil {
		fmt.Println("unsub error: ", token.Error())
	}

	srcClient.Disconnect(250)

	for _, c := range dstClients {
		c.Disconnect(250)
	}

	time.Sleep(1 * time.Second)
}
