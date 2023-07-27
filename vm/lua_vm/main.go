package main

import (
	"fmt"
	"io"
	"os"

	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

func SayHi(L *lua.LState) int {
	name := L.ToString(1)
	line := fmt.Sprintln("hi", name)
	L.Push(lua.LString(line))
	return 1
}

func main() {
	L := lua.NewState()
	defer L.Close()
	libs.Preload(L)

	L.SetGlobal("gohi", L.NewFunction(SayHi))

	f, err := os.Open("my.lua")
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
}
