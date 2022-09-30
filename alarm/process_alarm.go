package alarm

import (
	"fmt"
	"time"
)

type TimeWindow struct {
	LastTime time.Time
	LastNum  uint64
	Interval time.Duration
}

func AlarmInTimeWindow() ProcessFunc {
	window := TimeWindow{LastTime: time.Now(), Interval: 3 * time.Second}

	return func(met *AlarmMetric, now time.Time) *AlarmEvent {
		if now.Sub(window.LastTime) > window.Interval {
			window.LastTime = now
			diff := met.ValueInt - window.LastNum
			window.LastNum = met.ValueInt
			if diff > 2 {
				//fmt.Println("alarm", met.ValueInt)
				return &AlarmEvent{AlarmStatus: true, Value: fmt.Sprintf("%d", diff)}
			} else {
				//fmt.Println("alarm ok", met.ValueInt)
				return &AlarmEvent{AlarmStatus: false}
			}

		}
		return nil
	}
}

func testTimeWindow(met *AlarmMetric, now time.Time, lastTime time.Time, interval time.Duration) (*AlarmEvent, time.Time) {
	if now.Sub(lastTime) > interval {
		if met.ValueInt > 2 {
			// fmt.Println("alarm", met.ValueInt)
			return &AlarmEvent{AlarmStatus: true, Value: fmt.Sprintf("%d", met.ValueInt)}, now
		} else {
			// fmt.Println("alarm ok", met.ValueInt)
			return &AlarmEvent{AlarmStatus: false}, now
		}
	}
	return nil, lastTime
}

func AlarmInTimeWindow2() ProcessFunc {
	window := TimeWindow{LastTime: time.Now(), Interval: 3 * time.Second}

	return func(met *AlarmMetric, now time.Time) *AlarmEvent {
		ret, lastTime := testTimeWindow(met, now, window.LastTime, window.Interval)
		if ret != nil {
			window.LastTime = lastTime
			return ret
		}
		return nil
	}
}

func AlarmDirect(met *AlarmMetric, now time.Time) *AlarmEvent {
	state := false
	if met.State != 0 {
		state = true
	}
	return &AlarmEvent{AlarmStatus: state, Value: fmt.Sprintf("%d", met.ValueInt)}
}
