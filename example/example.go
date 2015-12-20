package example

import "C"
import "unsafe"

//export Handler
var Handler = func(event, context *C.char) (result *C.char, size C.size_t) {

	str := "hello world"
	hdr := (*stringHeader)(unsafe.Pointer(&str))
	result = (*C.char)(unsafe.Pointer(hdr.Data))
	size = (C.size_t)(hdr.Len)
	return
	// buf := new(bytes.Buffer)
	// fmt.Fprintf(buf, "Hello from %s! Mem allocated %s\n", c.FunctionName, c.MemoryLimit)
	// fmt.Fprintf(buf, "Params you've passed: %#v\n", event)
	// fmt.Fprintln(buf, "Current time:", time.Now().Format(time.Kitchen))
	// return buf.String()
}

type stringHeader struct {
	Data uintptr
	Len  int
}
