package lambda

import (
	"encoding/json"
	"unsafe"
)

import "C"

type Dict map[string]interface{}

type HandlerFunc func(event Dict, context *Context) []byte

type Bridge func(event, context *C.char) (result *C.char, size C.size_t)

func Use(fn HandlerFunc) Bridge {
	return func(eventData, ctxData *C.char) (result *C.char, size C.size_t) {
		var event Dict
		json.Unmarshal(bytesFrom(eventData), &event)
		var context *Context
		json.Unmarshal(bytesFrom(ctxData), &context)

		r := fn(event, context)

		hdr := (*sliceHeader)(unsafe.Pointer(&r))
		result = (*C.char)(unsafe.Pointer(hdr.Data))
		size = (C.size_t)(hdr.Len)
		return
	}
}

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
