package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	libs "github.com/vadv/gopher-lua-libs"
	libs_mat "github.com/vadv/gopher-lua-libs/internal_matrics"
	lua "github.com/yuin/gopher-lua"
)

var (
	fLuaFile    = flag.String("f", "main.lua", "lua to run")
	fInputFile  = flag.String("i", "input.json", "input args")
	fOutputFile = flag.String("o", "output", "output file")
	fUploadUrl  = flag.String("up", "", "upload url")
)

type InputArgs struct {
}

func LoadInput(path string) (interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData) // attention: &
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func decode(T *lua.LTable, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := &lua.LTable{}
		for _, item := range converted {
			arr.Append(decode(T, item))
		}
		return arr
	case map[string]interface{}:
		tbl := &lua.LTable{}
		// L.SetMetatable(tbl, L.GetTypeMetatable(jsonTableIsObject))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), decode(tbl, item))
		}
		return tbl
	case nil:
		return lua.LNil
	}
	panic("unreachable")
}

func NewLuaInputFunc(input interface{}) lua.LGFunction {
	inputTable := lua.LTable{}
	Lv := decode(&inputTable, input)

	return func(L *lua.LState) int {
		L.Push(Lv)
		return 1
	}
}

func save(data string) error {
	f, err := os.OpenFile(*fOutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	// 写入内容到文件
	_, err = f.WriteString(data + "\n")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return nil
}

type KafkaBody struct {
	Topic string `json:"topic"`
	Key   string `json:"key"`
	Body  string `json:"body"`
}

type KafkaBatchBody []KafkaBody

type HttpResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func HttpReport(req *http.Request) error {

	hc := http.Client{}

	var err error
	retry := 3
	for i := 0; i < retry; i++ {
		var resp *http.Response
		resp, err = hc.Do(req)

		ok := func() bool {
			if err != nil {
				fmt.Printf("http report request error: %s\n", err)
				return false
			}

			var data []byte
			data, err = io.ReadAll(resp.Body)

			fmt.Printf("resp body %s\n", data)

			if err != nil {
				fmt.Printf("http report read body error: %s\n", err)
				return false
			}

			hr := HttpResponse{}
			err = json.Unmarshal(data, &hr)
			if err != nil {
				fmt.Printf("http report unmarshal error: %s\n", err)
				return false
			}

			if hr.Code != 200 {
				err = fmt.Errorf("http report code error: %s", hr.Msg)
				fmt.Printf("http report code error: %s\n", hr.Msg)
				return false
			} else {
				err = nil
				return true
			}
		}()

		resp.Body.Close()

		if ok {
			break
		}

		if i < (retry - 1) { // skip last one
			time.Sleep(time.Duration(1+i*2) * time.Second)
		}
	}

	return err
}

func KafkaReport(id, topic, key, data string) error {
	fmt.Println("KafkaReport")

	msg := KafkaBody{
		Topic: topic,
		Key:   key,
		Body:  data,
	}

	jsonStr, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("gid", id)

	fullUrl := *fUploadUrl + "?" + params.Encode()
	fmt.Println("url", fullUrl)

	req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	// resq_data, _ := httputil.DumpRequest(req, false)
	// fmt.Println(string(resq_data))

	err = HttpReport(req)
	if err != nil {
		return err
	}

	return nil
}

func KafkaReportBatch(id, topic, key, data string) error {
	fmt.Println("KafkaReportBatch")

	msgs := []json.RawMessage{}
	err := json.Unmarshal([]byte(data), &msgs)
	if err != nil {
		return fmt.Errorf("data is not json array: %s", err)
	}

	kakfaBatch := KafkaBatchBody{}
	for _, d := range msgs {
		kMsg := KafkaBody{
			Topic: topic,
			Key:   key,
			Body:  string(d),
		}
		kakfaBatch = append(kakfaBatch, kMsg)
	}

	jsonStr, err := json.Marshal(&kakfaBatch)
	if err != nil {
		return fmt.Errorf("marshal kafka batch err: %s", err)
	}

	params := url.Values{}
	params.Add("gid", id)
	params.Add("batch", "")

	fullUrl := *fUploadUrl + "?" + params.Encode()
	fmt.Println("url", fullUrl)

	req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	// resq_data, _ := httputil.DumpRequest(req, true)
	// fmt.Println(string(resq_data))

	err = HttpReport(req)
	if err != nil {
		return err
	}

	return nil
}

type ReportFunc func(string) error

func NewLuaReportFunc(input interface{}) lua.LGFunction {
	// use parameters in output
	report := func() ReportFunc {
		inputMap := input.(map[string]interface{})
		if inputMap == nil {
			return nil
		}

		if _, ok := inputMap["upload"]; !ok {
			return nil
		}

		upload := inputMap["upload"]
		uploadMap := upload.(map[string]interface{})
		if uploadMap == nil {
			return nil
		}

		uplaodType, _ := uploadMap["type"].(string) // var , _ := // if not ok, then empty string
		topic, _ := uploadMap["topic"].(string)
		id, _ := uploadMap["id"].(string)

		switch uplaodType {
		case "kafka":
			return func(data string) error {
				return KafkaReport(id, topic, "", data)
			}
		case "kafka_batch":
			return func(data string) error {
				return KafkaReportBatch(id, topic, "", data)
			}
		default:
			fmt.Println("type not support:", uplaodType)
		}

		return nil
	}()

	return func(L *lua.LState) int {
		var err error
		var data string

		if L.GetTop() == 1 {
			data = L.ToString(1) /* get argument */
			if report != nil {
				err = report(data)
			} else {
				err = save(data)
			}
		} else if L.GetTop() > 1 {
			reportType := L.ToString(1)
			if reportType == "kafka" {
				id := L.ToString(2)
				topic := L.ToString(3)
				key := L.ToString(4)
				data := L.ToString(5)
				// fmt.Printf("type:%s, cluster:%s, topic:%s, key:%s, data:%s\n", reportType, cluster, topic, key, data)
				err = KafkaReport(id, topic, key, data)
			} else if reportType == "http" {
				url := L.ToString(2)
				body := L.ToString(3)
				fmt.Printf("type:%s, url:%s, body:%s\n", reportType, url, body)
			}
		}

		if err != nil {
			L.Push(lua.LString(err.Error())) /* push result */
		} else {
			L.Push(lua.LNil)
		}

		return 1 /* number of results */
	}
}

func someTest() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("----- begin ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	var wg sync.WaitGroup
	ch := make(chan int, 10)

	// 启动 10 个 goroutine，并为每个 goroutine 增加一个等待计数
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			// 在函数退出时，调用 Done 方法，表示一个 goroutine 结束
			defer wg.Done()
			// 从 channel 中读取数据，并处理它
			for _ = range ch {
				doJob()
			}
		}()
	}

	// 向 channel 中写入所有数据
	for i := 0; i < 3000; i++ {
		ch <- i
		fmt.Printf("\r%d", i)
	}
	close(ch)

	// 等待所有 goroutine 结束
	wg.Wait()

	runtime.ReadMemStats(&m)
	fmt.Println("----- before GC ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	runtime.GC()

	runtime.ReadMemStats(&m)
	fmt.Println("----- after ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)
}

func doJob() {
	jsonData, err := LoadInput(*fInputFile)
	if err != nil {
		fmt.Println("LoadInput:", err)
		return
	}
	intputFunc := NewLuaInputFunc(jsonData)
	reportFunc := NewLuaReportFunc(jsonData)

	L := lua.NewState(lua.Options{
		RegistrySize:        1024,
		RegistryMaxSize:     1024 * 20,
		RegistryGrowStep:    32,
		MinimizeStackMemory: true,
		CallStackSize:       64,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	L.SetContext(ctx)
	defer cancel()
	defer L.Close()
	libs.Preload(L)

	L.SetGlobal("ipes_report", L.NewFunction(reportFunc))
	L.SetGlobal("ipes_input", L.NewFunction(intputFunc))

	// add lua mod in script dir
	baseDir := filepath.Dir(*fLuaFile)
	L.SetGlobal("ipes_script_dir", L.NewFunction(
		func(L *lua.LState) int {
			L.Push(lua.LString(baseDir))
			return 1
		},
	))
	err = L.DoString(fmt.Sprintf(`
	package.path = package.path .. ";%s" .. [[/?.lua]]
	`, baseDir))
	if err != nil {
		fmt.Println("lua runtime err:", err)
	}

	f, err := os.Open(*fLuaFile)
	if err != nil {
		fmt.Println("open err:", err)
		return
	}
	defer f.Close()

	lua_script, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("open err:", err)
		return
	}

	if err := L.DoString(string(lua_script)); err != nil {
		fmt.Println("lua runtime err:", err)
	}

	if keyValue, ok := libs_mat.GetMat(L); ok {
		for k, v := range keyValue {
			fmt.Println(k, ":", v)
		}
	}
}

func main() {
	flag.Parse()
	doJob()
	//someTest()
}
