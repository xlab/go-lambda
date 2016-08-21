package lambda

import (
	"encoding/json"
	"unsafe"
)

import "C"

type HandlerFunc func(event json.RawMessage, context *Context) []byte

// Use the provided HandlerFunc as handler for incoming requests. This function returns
// a bridge that manages values passing between the caller (python) and the target HandlerFunc.
//
// Note: do not pass the bridge directly to C, otherwise CGo checker will complain,
// instead use a map of handles, see module/module.go for example.
//
// Warning: bytes of eventData should be copied explicitly if you want to use them outside
// the HandlerFunc scope (e.g. in goroutines), they are valid until the Bridge func returns.
func Use(fn HandlerFunc) bridge {
	return func(eventData, ctxData *C.char) (result *C.char, size C.size_t) {
		var context *Context
		json.Unmarshal(bytesFrom(ctxData), &context)
		event := json.RawMessage(bytesFrom(eventData))

		buf := fn(event, context)
		result = C.CString(string(buf)) // TODO(xlab): optimize if needed
		size = (C.size_t)(len(buf))
		return
	}
}

type bridge func(event, context *C.char) (result *C.char, size C.size_t)

func bytesFrom(p *C.char) []byte {
	var slice []byte
	if p != nil && *p != 0 {
		h := (*sliceHeader)(unsafe.Pointer(&slice))
		h.Data = uintptr(unsafe.Pointer(p))
		for *p != 0 {
			p = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 1)) // p++
		}
		h.Len = int(uintptr(unsafe.Pointer(p)) - h.Data)
		h.Cap = h.Len
	}
	return slice
}

type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
