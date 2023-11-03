package main

import (
	"fmt"
	. "goproject/prototype/lambda/server/internal"
	"time"
)

func main() {
	fmt.Println("lambda server test")

	j := &Job{
		Affinity: map[string]Affinity{},
	}

	j.Affinity["a"] = Affinity{
		Isp:      "ct",
		Province: "shangdong",
	}
	j.Affinity["b"] = Affinity{
		Isp:      "cnc",
		Province: "hebei",
	}
	j.Affinity["c"] = Affinity{
		Sn: []string{
			"sn_1",
			"sn_2",
			"sn_3",
		},
	}

	affinityName := []string{"a", "b", "c"}

	args := map[string]interface{}{
		"a": "b",
	}

	num := 10
	for i := 0; i < num; i++ {

		j.Task = append(j.Task, Task{
			Args:     args,
			Affinity: []string{affinityName[i%len(affinityName)]},
		})
	}

	j.FuncInfo = FunctionInfo{
		Name: "func1",
		// Version: "v1",
		Version: "all",
	}

	j.LambdaBehaviour = Behaviour{
		Replica: 3,
	}

	j.Upload.Cluster = "test_cluster"
	j.Upload.Type = "test"

	// Setup()
	testDevs := NewTestDevs()
	affinityDevs := NewCacahedAffinityDevs(testDevs, 10*time.Minute)

	err := DispatchJob(j, affinityDevs)
	if err != nil {
		fmt.Println("DispatchJob err:", err)
	}
}
