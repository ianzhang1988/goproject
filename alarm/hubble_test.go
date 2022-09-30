package alarm

import (
	"encoding/json"
	"testing"
)

func TestStuff(t *testing.T) {
	a := AlarmMsg{Group: "IPES", Triggername: "test", Status: "OK", Alertlevel: "P3", Value: "something is wrong", Pushedtags: Tags{Test: "hi"}}
	data, _ := json.Marshal(a)
	t.Log("alarm: ", string(data))

	b := AlarmMsg{Group: "IPES", Triggername: "test", Status: "OK", Alertlevel: "P3", Value: "something is wrong", Pushedtags: map[string]string{"test": "hello"}}
	data, _ = json.Marshal(b)
	t.Log("alarm: ", string(data))

	// err := SendAlarm(&a)
	// t.Log("send:", err)

	config := HubbleConfig{}
	configString := `
	{
		"RetryNum": 1,
		"RetryInterval": 15
	}
	`
	err := json.Unmarshal([]byte(configString), &config)
	if err != nil {
		t.Logf("config err: %s", err)
	}

	t.Logf("config: %v", config)
}
