package alarm

import "fmt"

func SendToStd(event *AlarmEvent) {
	fmt.Printf("name %s state %v, value %s \n", event.Name, event.AlarmStatus, event.Value)
}
