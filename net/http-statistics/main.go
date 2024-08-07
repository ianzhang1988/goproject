package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

func timeGet(url string) {
	req, _ := http.NewRequest("GET", url, nil)

	// body := bytes.NewBufferString("test")
	// req, _ := http.NewRequest("POST", url, body)

	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Done: %v\n", time.Since(dns))
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			fmt.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			fmt.Printf("Connect time: %v\n", time.Since(connect))
		},

		GotFirstResponseByte: func() {
			fmt.Printf("Time from start to first byte: %v\n", time.Since(start))
		},

		WroteRequest: func(wri httptrace.WroteRequestInfo) {
			fmt.Printf("Time wrote request: %v\n", time.Since(start))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	// if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
	// 	log.Fatal(err)
	// }
	// if _, err := http.DefaultClient.Do(req); err != nil {
	// 	log.Fatal(err)
	// }
	client := http.Client{}
	if _, err := client.Do(req); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}

func main() {
	timeGet("http://www.baidu.com")
}
