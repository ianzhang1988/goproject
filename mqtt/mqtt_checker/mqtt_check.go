package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var connected mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Printf("Connected: %v\n", client.IsConnectionOpen())
}

var host = flag.String("host", "", "host")

// var port = flag.Int("port", 1883, "port")
var repeatNum = flag.Int("num", 10, "repeat num")
var caFile = flag.String("ca", "", "ca file")
var password = flag.String("password", "", "The password (optional)")
var user = flag.String("user", "", "The User (optional)")

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(*caFile)
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

func main() {
	flag.Parse()

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// url := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Println("conneting to:", *host)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(*host)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetOnConnectHandler(connected)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetCleanSession(true)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	if *caFile != "" {
		tlsconfig := NewTLSConfig()
		opts.SetTLSConfig(tlsconfig)
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("test_checker", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println("sub error: ", token.Error())
		os.Exit(1)
	}

	for i := 0; i < *repeatNum; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("test_checker", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}

	// time.Sleep(6 * time.Second)

	if token := c.Unsubscribe("test_checker"); token.Wait() && token.Error() != nil {
		fmt.Println("unsub error: ", token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)

	time.Sleep(1 * time.Second)
}
