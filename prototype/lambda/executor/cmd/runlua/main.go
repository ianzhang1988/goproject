package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

var (
	fLuaFile    = flag.String("f", "main.lua", "lua to run")
	fInputFile  = flag.String("i", "input.json", "input args")
	fOutputFile = flag.String("o", "output", "output file")
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

func Report(L *lua.LState) int {
	data := L.ToString(1) /* get argument */

	err := save(data)
	if err != nil {
		L.Push(lua.LString(err.Error())) /* push result */
	} else {
		L.Push(lua.LNil)
	}

	return 1 /* number of results */
}

func doJob() {
	jsonData, err := LoadInput(*fInputFile)
	if err != nil {
		fmt.Println("LoadInput:", err)
		return
	}
	intputFunc := NewLuaInputFunc(jsonData)

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

	L.SetGlobal("ipes_report", L.NewFunction(Report))
	L.SetGlobal("ipes_input", L.NewFunction(intputFunc))

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
		panic(err)
	}

	// if keyValue, ok := libs_mat.GetMat(L); ok {
	// 	for k, v := range keyValue {
	// 		fmt.Println(k, ":", v)
	// 	}
	// }
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
	for i := 0; i < 10; i++ {
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
	for i := 0; i < 20000; i++ {
		ch <- i
		fmt.Printf("\r%d", i)
	}
	close(ch)

	// 等待所有 goroutine 结束
	wg.Wait()

	runtime.GC()

	runtime.ReadMemStats(&m)
	fmt.Println("----- after ------")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)
}

func main() {
	flag.Parse()
	// doJob()
	someTest()
}
