package alarm

import (
	"log"
	"sync"
	"time"
)

type AlarmMetric struct {
	State      int
	ValueInt   uint64
	ValueFloat float64
}

type Alarm struct {
	AlarmMetrics map[string]*AlarmMetric
	Lock         map[string]*sync.Mutex
}

type AlarmEvent struct {
	Name        string
	AlarmStatus bool
	Value       string
	Tags        map[string]string
}

type ProcessFunc func(met *AlarmMetric, time time.Time) *AlarmEvent
type SendAlarmFunc func(event *AlarmEvent)

var GAlarm Alarm
var ProcessFuncMap map[string]ProcessFunc
var LastAlarmEvent map[string]*AlarmEvent
var AlarmWhenToggle map[string]bool
var SendFuncs []SendAlarmFunc
var WorkFlag bool

func init() {
	WorkFlag = true
	GAlarm = Alarm{AlarmMetrics: map[string]*AlarmMetric{}, Lock: map[string]*sync.Mutex{}}
	ProcessFuncMap = map[string]ProcessFunc{}
	LastAlarmEvent = map[string]*AlarmEvent{}
	AlarmWhenToggle = map[string]bool{}
}

func RegisterProcessFunc(name string, f ProcessFunc) {
	RegisterProcessFuncExt(name, f, true)
}

func RegisterProcessFuncExt(name string, f ProcessFunc, alarmWhenToggle bool) {
	ProcessFuncMap[name] = f
	// seprate init between processfunc and metric,
	// why? if one collect metric "a", and not register a process func "a"
	// alarm Set would panic due to nil pointer
	// GAlarm.AlarmMetrics[name] = &AlarmMetric{}
	// GAlarm.Lock[name] = &sync.Mutex{}
	AlarmWhenToggle[name] = alarmWhenToggle
}

func RegisterSendFunc(f SendAlarmFunc) {
	SendFuncs = append(SendFuncs, f)
}

func Strat() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for WorkFlag {

			<-ticker.C

			now := time.Now()

			for name, f := range ProcessFuncMap {
				met := GAlarm.Get(name)

				if met == nil {
					log.Println("alarm get metirc is nil, should not happen, check your code!")
					continue
				}

				event := f(met, now)
				if event == nil { // not determined
					continue
				}
				event.Name = name

				// process only when alarm status toggled
				if AlarmWhenToggle[name] {
					lastEvent := LastAlarmEvent[name]
					if lastEvent != nil && lastEvent.AlarmStatus == event.AlarmStatus {
						continue
					}
				}

				LastAlarmEvent[name] = event

				// send alarm
				for _, f := range SendFuncs {
					go f(event)
				}
			}
		}
	}()
}

func Stop() {
	WorkFlag = false
}

func (a *Alarm) Init(name string) {
	a.AlarmMetrics[name] = &AlarmMetric{}
	a.Lock[name] = &sync.Mutex{}
}

func (a *Alarm) Check(name string) bool {
	_, lockOk := a.Lock[name]
	_, metricOk := a.AlarmMetrics[name]
	return lockOk && metricOk
}

func (a *Alarm) Get(name string) *AlarmMetric {
	if !a.Check(name) {
		return nil
	}

	a.Lock[name].Lock()
	defer a.Lock[name].Unlock()

	met := AlarmMetric{}
	met = *a.AlarmMetrics[name]
	return &met
}

func (a *Alarm) SetState(name string, state int) {
	if !a.Check(name) {
		return
	}

	a.Lock[name].Lock()
	defer a.Lock[name].Unlock()
	a.AlarmMetrics[name].State = state
}

func (a *Alarm) SetInt(name string, value uint64) {
	if !a.Check(name) {
		return
	}

	a.Lock[name].Lock()
	defer a.Lock[name].Unlock()
	a.AlarmMetrics[name].ValueInt = value
}

func (a *Alarm) AddInt(name string, value uint64) {
	if !a.Check(name) {
		return
	}

	a.Lock[name].Lock()
	defer a.Lock[name].Unlock()
	a.AlarmMetrics[name].ValueInt += value
}

func (a *Alarm) DecInt(name string, value uint64) {
	if !a.Check(name) {
		return
	}

	a.Lock[name].Lock()
	defer a.Lock[name].Unlock()
	a.AlarmMetrics[name].ValueInt -= value
}

func (a *Alarm) SetFloat(name string, value float64) {
	if !a.Check(name) {
		return
	}

	a.Lock[name].Lock()
	defer a.Lock[name].Unlock()
	a.AlarmMetrics[name].ValueFloat = value
}

func (a *Alarm) AddFloat(name string, value float64) {
	if !a.Check(name) {
		return
	}

	a.Lock[name].Lock()
	defer a.Lock[name].Unlock()
	a.AlarmMetrics[name].ValueFloat += value
}

func (a *Alarm) DecFloat(name string, value float64) {
	if !a.Check(name) {
		return
	}

	a.Lock[name].Lock()
	defer a.Lock[name].Unlock()
	a.AlarmMetrics[name].ValueFloat -= value
}
