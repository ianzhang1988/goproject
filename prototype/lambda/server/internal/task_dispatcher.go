package lambda_service

import "fmt"

//

type TaskDispatcher struct {
}

func Setup() {
	UpdateDB()
}

func SendTask(tasks []TaskMatch) error {
	for i := range tasks {
		t := &(tasks[i])
		// if t.Task == nil || t.Dev == nil {
		// 	fmt.Printf("%d: Task: %+v Affinity: %s Dev: %+v Upload %+v \n", i, t.Task, t.Affinity, t.Dev, t.Upload)
		// 	break
		// }
		fmt.Printf("%d: Task: %+v Affinity: %s Dev: %+v Upload %+v \n", i, *(t.Task), t.Affinity, *(t.Dev), t.Upload)
	}
	return nil
}

func DispatchJob(job *Job) {
	
}

func DispatchTask(job *Job, affinityDevs AffinityDevs) error {

	affinityGrp := GetAffinityDevGroup(job, affinityDevs)
	tasks, err := MatchDev(job, affinityGrp)
	if err != nil {
		return err
	}

	// filter task without device, record error

	err = SendTask(tasks)
	if err != nil {
		return err
	}

	return nil
}
