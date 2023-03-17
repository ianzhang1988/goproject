package stress

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var udp_counter uint64 = 0

type UdpStressCase struct {
	Data       []byte
	TargetIP   string
	TargetPort string
	Conn       *net.UDPConn
	ReadResp   bool
}

func NewUdpStressCase(data []byte, target_ip string, target_port string, read_resp bool) (*UdpStressCase, error) {

	addr, err := net.ResolveUDPAddr("udp", target_ip+":"+target_port)
	if err != nil {
		return nil, err
		// log.Fatalln(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
		// log.Fatalln(err)
	}

	Case := &UdpStressCase{
		Data:       data,
		TargetIP:   target_ip,
		TargetPort: target_port,
		Conn:       conn,
	}

	return Case, nil
}

func ReleaseUdpStressCase(Case *UdpStressCase) error {
	return Case.Conn.Close()
}

func ReleaseAllUdpStressCase(Cases []*UdpStressCase) {
	for _, Case := range Cases {
		ReleaseUdpStressCase(Case)
	}
}

func (c *UdpStressCase) Do() error {
	_, err := c.Conn.Write(c.Data)
	if err != nil {
		return err
	}

	if c.ReadResp {
		data := make([]byte, 4096)
		_, _, err := c.Conn.ReadFromUDP(data) // 接收数据
		if err != nil {
			return err
		}
	}

	return nil
}

func doUdpstress(Case *UdpStressCase, num int, interval time.Duration) {
	// time.Nanosecond
	for i := 0; i < num; i++ {
		err := Case.Do()
		if err != nil {
			log.Printf("udp faild: %s\n", err)
		} else {
			atomic.AddUint64(&udp_counter, 1)
		}
		time.Sleep(interval)
	}

}

func ShowUdpCounter() {
	last_value := udp_counter
	t := time.NewTicker(60 * time.Second)

	for {
		<-t.C
		c := atomic.LoadUint64(&udp_counter)
		log.Println("sent:", c-last_value)
		last_value = c
	}
}

func UdpStress(Cases []*UdpStressCase, num int, rate int) error {
	go ShowUdpCounter()

	wg := sync.WaitGroup{}
	num_each := num / len(Cases)
	rate_each := rate / len(Cases)
	sleep_time := time.Duration(1) * time.Second / time.Duration(rate_each)
	log.Println("sleep_time:", sleep_time)
	for _, Case := range Cases {
		wg.Add(1)
		go func(Case *UdpStressCase) {
			doUdpstress(Case, num_each, sleep_time)
			wg.Done()
		}(Case)
	}
	wg.Wait()
	return nil
}
