package example

import (
	"bytes"
	"fmt"
	"time"

	"github.com/xlab/go-lambda/lambda"
)

type Context struct {
	*lambda.Context
}

func MakeContext() Context {
	return Context{&lambda.Context{}}
}

type Dict struct {
	*lambda.Dict
}

func MakeDict(size int) Dict {
	return Dict{lambda.MakeDict(size)}
}

func Handler(event Dict, c Context) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "Hello from %s! Mem allocated %s\n", c.FunctionName, c.MemoryLimit)
	fmt.Fprintf(buf, "Params you've passed: %#v\n", event)
	fmt.Fprintln(buf, "Current time:", time.Now().Format(time.Kitchen))
	return buf.String()
}
