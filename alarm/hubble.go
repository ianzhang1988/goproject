package alarm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type HubbleConfig struct {
	RetryNum      int
	RetryInterval int
	ServerUrl     string
	UserToken     string
}

func SendRetry(n int, interval time.Duration, f func() error) {
	for i := 0; i < n; i++ {
		if f() == nil {
			break
		}

		time.Sleep(interval)
	}
}

func NewSendToHubbleFunc(config HubbleConfig) SendAlarmFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	return func(event *AlarmEvent) {
		SendRetry(config.RetryNum, time.Duration(config.RetryInterval)*time.Second, func() error {
			status := "OK"
			if event.AlarmStatus {
				status = "PROBLEM"
			}

			msg := AlarmMsg{Group: "Test", Triggername: event.Name, Status: status, Alertlevel: "P2", Value: event.Value, Endpoint: hostname}
			return doHubbleSend(config.ServerUrl, config.UserToken, &msg)
		})
	}
}

type AlarmMsg struct {
	Group       string      `json:"grp"`
	Triggername string      `json:"triggername"`
	Status      string      `json:"status"`
	Alertlevel  string      `json:"alertlevel"`
	Value       string      `json:"value"`
	Endpoint    string      `json:"endpoint"`
	Pushedtags  interface{} `json:"pushedtags"`
}

type Tags struct {
	Test string
}

func doHubbleSend(url, usrtoken string, msg *AlarmMsg) error {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshal error")
		return err
	}
	body := strings.NewReader(string(data))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println("NewRequest error")
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("usertoken", usrtoken)

	clt := http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		log.Println("http error")
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("not 200")
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ReadAll error")
		return err
	}
	log.Println("resp", string(content))
	return nil
}
