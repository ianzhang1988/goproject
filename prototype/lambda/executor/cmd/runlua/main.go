package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
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

func main() {
	flag.Parse()

	jsonData, err := LoadInput(*fInputFile)
	if err != nil {
		fmt.Println("LoadInput:", err)
		return
	}
	intputFunc := NewLuaInputFunc(jsonData)

	L := lua.NewState()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	L.SetContext(ctx)
	defer cancel()
	defer L.Close()
	libs.Preload(L)

	metric := map[string]string{}
	MetricUD := L.NewUserData()
	MetricUD.Value = metric

	L.SetGlobal("lib_metric", MetricUD)
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

	metricLV := L.GetGlobal("lib_metric")
	metricUd, ok := metricLV.(*lua.LUserData)
	if ok && metricUd.Value != nil {
		if keyValue, ok := metricUd.Value.(map[string]string); ok {
			fmt.Println("kv:", keyValue)
		}
	}
}
