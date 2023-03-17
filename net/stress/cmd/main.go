package main

import (
	"flag"
	stress "goproject/net/stress/internal"
	"io/ioutil"
	"log"
	"os"
)

var (
	Ip        = flag.String("ip", "127.0.0.1", "target ip")
	Port      = flag.String("port", "80", "target port")
	Num       = flag.Int("num", 1000000, "")
	WorkerNum = flag.Int("worker", 1, "")
	Type      = flag.String("t", "udp", "udp/http")
	Rate      = flag.Int("rate", 10000, "per second")
	// KeepAlive = flag.Int("keepalive", 1, "")
	KeepAlive = flag.Bool("keepalive", true, "")
	ConnPool  = flag.Bool("connpool", false, "")
)

func main() {
	flag.Parse()

	f, err := os.Open("bin.dat")
	if err != nil {
		log.Fatalf("openfile failed: %s\n", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("read file failed: %s\n", err)
	}

	log.Println("type:", *Type)

	if *Type == "udp" {
		UdpCases := []*stress.UdpStressCase{}
		for i := 0; i < *WorkerNum; i++ {
			Case, err := stress.NewUdpStressCase(data, *Ip, *Port, false)
			if err != nil {
				log.Fatalf("Case failed: %s\n", err)
			}
			UdpCases = append(UdpCases, Case)
		}

		stress.UdpStress(UdpCases, *Num, *Rate)
		stress.ReleaseAllUdpStressCase(UdpCases)
	} else if *Type == "http" {
		HttpCases := []*stress.HTTPStressCase{}
		for i := 0; i < *WorkerNum; i++ {
			Case, err := stress.NewHttpStressCase(data, *Ip, *Port, "httpdata", *KeepAlive, *ConnPool)
			if err != nil {
				log.Fatalf("Case failed: %s\n", err)
			}
			HttpCases = append(HttpCases, Case)
		}

		stress.HttpStress(HttpCases, *Num, *Rate)
	} else if *Type == "longhttp" {
		HttpCases := []stress.StressCase{}
		for i := 0; i < *WorkerNum; i++ {
			Case, err := stress.NewLongHttpStressCase(data, *Ip, *Port, "httpdata", 5)
			if err != nil {
				log.Fatalf("Case failed: %s\n", err)
			}
			HttpCases = append(HttpCases, Case)
		}

		stress.DoStress(HttpCases, *Num, *Rate)
	} else {
		log.Println("type:", *Type, "not supported")
	}
}
