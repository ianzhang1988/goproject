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

func (dg *DevGroup) Next() *Device {
	dev := dg.devs[dg.idx]
	dg.idx += 1
	dg.idx = dg.idx % uint32(len(dg.devs))
	return dev
}

func (dg *DevGroup) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(dg.devs), func(i, j int) {
		dg.devs[i], dg.devs[j] = dg.devs[j], dg.devs[i]
	})
}

type AffinityDevGroup map[string]*DevGroup

type TaskMatch struct {
	Task     *TaskArgs
	Affinity string
	Dev      *Device
	Upload   Uplaod
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

func GetAffinityDevGroup(job *Job, affinityDevs AffinityDevs) AffinityDevGroup {

	funcName := job.FuncInfo.Name
	// funcVersion, err := ParseVersion(job.FuncInfo.Version)
	// if err != nil {
	// 	return JobResult{}, fmt.Errorf("parse version failed: %s", err)
	// }
	funcVersion := job.FuncInfo.Version
	// sn2dev := GetSn2DevMap()

	affinityGrp := AffinityDevGroup{}

	if job.Affinity == nil {

		dg := &DevGroup{}

		devs, err := affinityDevs.GetAffinityDevs(funcName, funcVersion, Affinity{All: true})
		if err != nil {
			fmt.Printf("GetAffinityDevGroup get devs 1 err: %s", err)
		}
		dg.devs = append(dg.devs, devs...)

		dg.Shuffle()

		affinityGrp[""] = dg
	}

	for k, v := range job.Affinity {
		dg := &DevGroup{}

		devs, err := affinityDevs.GetAffinityDevs(funcName, funcVersion, v)
		if err != nil {
			fmt.Printf("GetAffinityDevGroup get devs 2 err: %s", err)
		}
		dg.devs = append(dg.devs, devs...)

		dg.Shuffle()

		if len(dg.devs) > 0 {
			affinityGrp[k] = dg
		}
	}

	return affinityGrp
}

// func GetAffinityDevGroupOld(job *Job) AffinityDevGroup {
// 	mdb := GetDeviceMemDB()
// 	funcName := job.FuncInfo.Name
// 	// funcVersion, err := ParseVersion(job.FuncInfo.Version)
// 	// if err != nil {
// 	// 	return JobResult{}, fmt.Errorf("parse version failed: %s", err)
// 	// }
// 	funcVersion := job.FuncInfo.Version
// 	sn2dev := GetSn2DevMap()

// 	affinityGrp := AffinityDevGroup{}

// 	if job.Affinity == nil {

// 		dg := &DevGroup{}

// 		devs := mdb.FindPrefix(Affinity2Path(funcName, funcVersion, Affinity{}))

// 		for _, d := range devs {
// 			dg.devs = append(dg.devs, d.(*Device))
// 		}

// 		dg.Shuffle()

// 		affinityGrp[""] = dg
// 	}

// 	for k, v := range job.Affinity {
// 		if len(v.Sn) > 0 {
// 			continue
// 		}

// 		dg := &DevGroup{}

// 		devs := mdb.FindPrefix(Affinity2Path(funcName, funcVersion, v))
// 		for _, d := range devs {
// 			dg.devs = append(dg.devs, d.(*Device))
// 		}

// 		dg.Shuffle()

// 		if len(dg.devs) > 0 {
// 			affinityGrp[k] = dg
// 		}
// 	}

// 	for k, v := range job.Affinity {
// 		if len(v.Sn) == 0 {
// 			continue
// 		}

// 		dg := &DevGroup{}

// 		for _, sn := range v.Sn {
// 			if d, ok := sn2dev[sn]; ok {
// 				dg.devs = append(dg.devs, d)
// 			}
// 		}

// 		dg.Shuffle()

// 		affinityGrp[k] = dg
// 	}

// 	return affinityGrp
// }

func MatchDev(job *Job, affinityGrp AffinityDevGroup) ([]TaskMatch, error) {
	replica := job.LambdaBehaviour.Replica

	for k, v := range affinityGrp {
		fmt.Print(k, ":")
		num := 10
		if len(v.devs) < num {
			num = len(v.devs)
		}
		for _, d := range v.devs[:num] {
			fmt.Print(" ", *d)
		}
		fmt.Println("")
	}

	tasks := []TaskMatch{}
	for i := 0; i < int(replica); i++ {
		for idx := range job.Task {
			taskargs := &(job.Task[idx].Args)
			for _, aff := range job.Task[idx].Affinity {
				tm := TaskMatch{
					Task:     taskargs,
					Affinity: aff,
					Upload:   job.Upload,
				}
				tasks = append(tasks, tm)
			}
			if len(job.Task[idx].Affinity) == 0 {
				tm := TaskMatch{
					Task:     taskargs,
					Affinity: "",
					Upload:   job.Upload,
				}
				tasks = append(tasks, tm)
			}
		}
	}

	for i, t := range tasks {
		grp, ok := affinityGrp[t.Affinity]
		if !ok {
			// process error
			continue
		}

		tasks[i].Dev = grp.Next()

		// fmt.Println(t)
	}

	return tasks, nil
}
