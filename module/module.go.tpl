package main

//#cgo pkg-config: python-2.7 --cflags --libs
import "C"

import (
	"unsafe"

	lambda "in"
)

//export lambda_handler
type lambda_handler unsafe.Pointer

//export get_lambda_handler
func get_lambda_handler() lambda_handler {
	return (lambda_handler)(unsafe.Pointer(&lambda.{{.PackageFunc}}))
}

//export lambda_handler_call
func lambda_handler_call(p lambda_handler, event, context *C.char) (*C.char, C.size_t) {
	return (*(*func(*C.char, *C.char) (*C.char, C.size_t))(unsafe.Pointer(p)))(event, context)
}

func main() {}
