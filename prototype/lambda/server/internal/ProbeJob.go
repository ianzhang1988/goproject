package lambda_service

type FunctionInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Affinity struct {
	All      bool
	Isp      string   `json:"isp"`
	Province string   `json:"province"`
	City     string   `json:"city"`
	Sn       []string `json:"sn"`
}

type Uplaod struct {
	Type    string `json:"type"`
	Cluster string `json:"cluster"`
	Topic   string `json:"topic"`
}

type Behaviour struct {
	CustomStrategy string `json:"custom_strategy"`
	Replica        uint   `json:"replica"`
}

type TaskArgs map[string]interface{}

type Task struct {
	Affinity []string `json:"affinity"`
	Args     TaskArgs `json:"args"`
}

type Job struct {
	FuncInfo        FunctionInfo        `json:"func_info"`
	Task            []Task              `json:"task"`
	Affinity        map[string]Affinity `json:"affinity"`
	Upload          Uplaod              `json:"upload"`
	LambdaBehaviour Behaviour           `json:"lambda_behaviour"`
}

type JobResult struct {
	Msg string `json:"msg"`
}
