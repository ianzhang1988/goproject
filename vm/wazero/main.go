package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// greetWasm was compiled using `tinygo build -o greet.wasm -scheduler=none --no-debug -target=wasi greet.go`
//
//go:embed testdata/http.wasm
var greetWasm []byte

type HttpRequest struct {
	Url string
}

type HttpResponse struct {
	Body string
	Code int
	Err  string
}

// main shows how to interact with a WebAssembly function that was compiled
// from TinyGo.
//
// See README.md for a full description.
func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Instantiate a Go-defined module named "env" that exports a function to
	// log to the console.
	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(logString).Export("log").
		NewFunctionBuilder().WithFunc(httpGet).Export("http_get").
		Instantiate(ctx)
	if err != nil {
		log.Panicln(err)
	}

	// Note: testdata/greet.go doesn't use WASI, but TinyGo needs it to
	// implement functions such as panic.
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Instantiate a WebAssembly module that imports the "log" function defined
	// in "env" and exports "memory" and functions we'll use in this example.
	mod, err := r.Instantiate(ctx, greetWasm)
	if err != nil {
		log.Panicln(err)
	}

	// Get references to WebAssembly functions we'll use in this example.
	greet := mod.ExportedFunction("greet")
	greeting := mod.ExportedFunction("greeting")
	// These are undocumented, but exported. See tinygo-org/tinygo#2788
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	// Let's use the argument to this main function in Wasm.
	name := os.Args[1]
	nameSize := uint64(len(name))

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := malloc.Call(ctx, nameSize)
	if err != nil {
		log.Panicln(err)
	}
	namePtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, namePtr)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(uint32(namePtr), []byte(name)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			namePtr, nameSize, mod.Memory().Size())
	}

	// Now, we can call "greet", which reads the string we wrote to memory!
	_, err = greet.Call(ctx, namePtr, nameSize)
	if err != nil {
		log.Panicln(err)
	}

	// Finally, we get the greeting message "greet" printed. This shows how to
	// read-back something allocated by TinyGo.
	ptrSize, err := greeting.Call(ctx, namePtr, nameSize)
	if err != nil {
		log.Panicln(err)
	}

	greetingPtr := uint32(ptrSize[0] >> 32)
	greetingSize := uint32(ptrSize[0])

	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	if greetingPtr != 0 {
		defer func() {
			_, err := free.Call(ctx, uint64(greetingPtr))
			if err != nil {
				log.Panicln(err)
			}
		}()
	}

	// The pointer is a linear memory offset, which is where we write the name.
	if bytes, ok := mod.Memory().Read(greetingPtr, greetingSize); !ok {
		log.Panicf("Memory.Read(%d, %d) out of range of memory size %d",
			greetingPtr, greetingSize, mod.Memory().Size())
	} else {
		fmt.Println("go >>", string(bytes))
	}

	run := mod.ExportedFunction("run")
	run.Call(ctx)

}

func SetGoReturn(ctx context.Context, data string, mod api.Module) {
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")
	set := mod.ExportedFunction("set_return_str")

	dataSize := uint64(len(data))

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := malloc.Call(ctx, dataSize)
	if err != nil {
		log.Panicln(err)
	}
	dataPtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, dataPtr)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(uint32(dataPtr), []byte(data)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			dataPtr, dataSize, mod.Memory().Size())
	}

	// Now, we can call "greet", which reads the string we wrote to memory!
	_, err = set.Call(ctx, dataPtr, dataSize)
	if err != nil {
		log.Panicln(err)
	}
}

func logString(_ context.Context, m api.Module, offset, byteCount uint32) {
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	fmt.Println(string(buf))
}

func httpGet(ctx context.Context, m api.Module, offset, byteCount uint32) uint32 {
	fmt.Println("go >> httpGet")
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}

	fmt.Println("buf:", string(buf))

	hr := HttpRequest{}

	json.Unmarshal(buf, &hr)

	resp, err := http.Get(hr.Url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("code:", resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	SetGoReturn(ctx, string(data), m)

	// fmt.Println(string(data))
	return 42
}
