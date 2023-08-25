package lambda_service

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type DevGroup struct {
	idx  uint32
	devs []*Device
}

type AffinityGroup map[string]*DevGroup

type TaskMatch struct {
	Task     *TaskArgs
	Affinity string
	Dev      *Device
}

func Affinity2Path(name, version string, aff Affinity) string {

	pathSlice := []string{}
	pathSlice = append(pathSlice, name)
	pathSlice = append(pathSlice, version)
	if aff.Isp != "" {
		pathSlice = append(pathSlice, aff.Isp)
	}
	if aff.Province != "" {
		pathSlice = append(pathSlice, aff.Province)
	}
	if aff.City != "" {
		pathSlice = append(pathSlice, aff.City)
	}

	return strings.Join(pathSlice, "/")
}

func ParseVersion(ver string) ([]string, error) {
	return []string{ver}, nil
}

func MatchDev(job *Job) (JobResult, error) {
	mdb := GetDeviceMemDB()
	replica := job.LambdaBehaviour.Replica
	funcName := job.FuncInfo.Name
	funcVersion, err := ParseVersion(job.FuncInfo.Version)
	if err != nil {
		return JobResult{}, fmt.Errorf("parse version failed: %s", err)
	}
	sn2dev := GetSn2DevMap()

	affinityGrp := AffinityGroup{}

	rand.Seed(time.Now().UnixNano())
	randStart := rand.Uint32()

	if job.Affinity == nil {

		dg := &DevGroup{
			idx: randStart,
		}

		for _, ver := range funcVersion {
			devs := mdb.FindPrefix(Affinity2Path(funcName, ver, Affinity{}))
			for _, d := range devs {
				dg.devs = append(dg.devs, d.(*Device))
			}
		}

		affinityGrp[""] = dg
	}

	for k, v := range job.Affinity {
		if len(v.Sn) > 0 {
			continue
		}

		dg := &DevGroup{
			idx: randStart,
		}

		for _, ver := range funcVersion {
			devs := mdb.FindPrefix(Affinity2Path(funcName, ver, v))
			for _, d := range devs {
				dg.devs = append(dg.devs, d.(*Device))
			}
		}

		affinityGrp[k] = dg
	}

	for k, v := range job.Affinity {
		if len(v.Sn) == 0 {
			continue
		}

		dg := &DevGroup{
			idx: randStart,
		}

		for _, sn := range v.Sn {
			if d, ok := sn2dev[sn]; ok {
				dg.devs = append(dg.devs, d)
			}
		}

		affinityGrp[k] = dg
	}

	tasks := []TaskMatch{}
	for i := 0; i < int(replica); i++ {
		for idx := range job.Task {
			taskargs := &(job.Task[idx].Args)
			for _, aff := range job.Task[idx].Affinity {
				tm := TaskMatch{
					Task:     taskargs,
					Affinity: aff,
				}
				tasks = append(tasks, tm)
			}
			if len(job.Task[idx].Affinity) == 0 {
				tm := TaskMatch{
					Task:     taskargs,
					Affinity: "",
				}
				tasks = append(tasks, tm)
			}
		}
	}

	for _, t := range tasks {
		grp, ok := affinityGrp[t.Affinity]
		if !ok {
			// process error
			continue
		}
		// t.Dev =
	}

	return JobResult{}, nil
}
