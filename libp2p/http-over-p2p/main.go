package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
)

func main() {

	// Parse options from the command line
	port := flag.Int("p", 0, "wait for incoming connections")
	remoteId := flag.String("peer", "", "target peer to dial")
	msg := flag.String("msg", "", "msg to send")
	flag.Parse()

	r := rand.Reader

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		fmt.Printf("key err: %s\n", err)
		return
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *port)),
		// libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", *port)),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
		libp2p.NoSecurity,
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		fmt.Printf("new host err: %s\n", err)
	}

	fmt.Printf("I'm: %s\n", host.ID().Pretty())

	if *remoteId == "" {
		listener, _ := gostream.Listen(host, p2phttp.DefaultP2PProtocol)
		defer listener.Close()
		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				data = []byte(fmt.Sprintf("err: %s\n", err))
				fmt.Print(string(data))
				w.Write(data)
				return
			}
			w.Write(append([]byte("Hi!"), data...))
		})

		fmt.Printf("start listen on %d\n", *port)
		server := &http.Server{}
		server.Serve(listener)
	} else {
		tr := &http.Transport{}
		tr.RegisterProtocol("libp2p", p2phttp.NewTransport(host))
		client := &http.Client{Transport: tr}
		// res, err := client.Get(fmt.Sprintf("libp2p://%s/hello", *remoteId))
		res, err := client.Post(fmt.Sprintf("libp2p://%s/hello", *remoteId), "text/plain", strings.NewReader(*msg))
		if err != nil {
			fmt.Printf("get err: %s\n", err)
			return
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("get err: %s\n", err)
			return
		}
		fmt.Printf("resp: %s\n", string(data))
	}
}
