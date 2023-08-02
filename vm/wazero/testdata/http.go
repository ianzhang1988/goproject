package main

// #include <stdlib.h>
import "C"

import (
	"encoding/json"
	"fmt"
	"runtime"
	"unsafe"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

// greet prints a greeting to the console.
func greet(name string) {
	log(fmt.Sprint("wasm >> ", greeting(name)))
}

// log a message to the console using _log.
func log(message string) {
	ptr, size := stringToPtr(message)
	_log(ptr, size)
	runtime.KeepAlive(message) // keep message alive until ptr is no longer needed.
}

// _log is a WebAssembly import which prints a string (linear memory offset,
// byteCount) to the console.
//
//go:wasmimport env log
func _log(ptr, size uint32)

// greeting gets a greeting for the name.
func greeting(name string) string {
	return fmt.Sprint("Hello, ", name, "!")
}

// _greet is a WebAssembly export that accepts a string pointer (linear memory
// offset) and calls greet.
//
//export greet
func _greet(ptr, size uint32) {
	name := ptrToString(ptr, size)
	greet(name)
}

// _greeting is a WebAssembly export that accepts a string pointer (linear memory
// offset) and returns a pointer/size pair packed into a uint64.
//
// Note: This uses a uint64 instead of two result values for compatibility with
// WebAssembly 1.0.
//
//export greeting
func _greeting(ptr, size uint32) (ptrSize uint64) {
	name := ptrToString(ptr, size)
	g := greeting(name)
	ptr, size = stringToLeakedPtr(g)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

// ptrToString returns a string from WebAssembly compatible numeric types
// representing its pointer and length.
func ptrToString(ptr uint32, size uint32) string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

// stringToPtr returns a pointer and size pair for the given string in a way
// compatible with WebAssembly numeric types.
// The returned pointer aliases the string hence the string must be kept alive
// until ptr is no longer needed.
func stringToPtr(s string) (uint32, uint32) {
	ptr := unsafe.Pointer(unsafe.StringData(s))
	return uint32(uintptr(ptr)), uint32(len(s))
}

// stringToLeakedPtr returns a pointer and size pair for the given string in a way
// compatible with WebAssembly numeric types.
// The pointer is not automatically managed by TinyGo hence it must be freed by the host.
func stringToLeakedPtr(s string) (uint32, uint32) {
	size := C.ulong(len(s))
	ptr := unsafe.Pointer(C.malloc(size))
	copy(unsafe.Slice((*byte)(ptr), size), s)
	return uint32(uintptr(ptr)), uint32(size)
}

type HttpRequest struct {
	Url string
}

type HttpResponse struct {
	Body string
	Code int
	Err  string
}

var last_return string

//export set_return_str
func SetReturnStr(ptr, size uint32) {
	last_return = ptrToString(ptr, size)
}

func httpGet(req HttpRequest) {
	data, err := json.Marshal(&req)
	if err != nil {
		log(err.Error())
	}
	datastr := string(data)
	ptr, size := stringToPtr(datastr)
	log("call _httpGet")
	ptr = _httpGet(ptr, size)
	log(fmt.Sprintf("%d", ptr))
	runtime.KeepAlive(datastr) // keep message alive until ptr is no longer needed.
}

//go:wasmimport env http_get
func _httpGet(ptr, size uint32) uint32

//export run
func run() {
	log("run ...")
	req := HttpRequest{
		Url: "http://www.baidu.com",
	}
	httpGet(req)

	log("baidu data:" + last_return[:100])

	req = HttpRequest{
		Url: "http://www.iqiyi.com",
	}
	httpGet(req)

	log("iqiyi data:" + last_return[:100])
}
