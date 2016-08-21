package main

//#cgo pkg-config: python-2.7 --cflags --libs
//#include <stdlib.h>
import "C"

import (
	"sync"
	"sync/atomic"
	"unsafe"

	lambda "in"
)

//export lambda_handle
type lambda_handle int64

//export get_lambda_handle
func get_lambda_handle() lambda_handle {
	handle := lambda_handle(atomic.AddInt64(&bridgeHandlePool, 1))
	bridgeMux.Lock()
	bridgeMap[handle] = unsafe.Pointer(&lambda.Handler)
	bridgeMux.Unlock()
	return handle
}

//export lambda_handle_call
func lambda_handle_call(handle lambda_handle, event, context *C.char) (*C.char, C.size_t) {
	bridgeMux.RLock()
	fn := *(*bridge)(bridgeMap[handle])
	bridgeMux.RUnlock()
	str, size := fn(event, context)
	return str, size
}

type bridge func(event, context *C.char) (result *C.char, size C.size_t)

var (
	bridgeHandlePool int64
	bridgeMap        = make(map[lambda_handle]unsafe.Pointer, 10)
	bridgeMux        = new(sync.RWMutex)
)

func main() {}
