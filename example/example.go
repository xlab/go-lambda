// Package example shows how an AWS Lambda function can be implemented using Go.
//
// Build steps: `go-lambda update example-handler handler github.com/xlab/go-lambda/example`
package example

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xlab/go-lambda/lambda"
)

// Handler is a bridge that will be called by a python module wrapping this package.
var Handler = lambda.Use(lambda.HandlerFunc(exampleHandler))

func exampleHandler(event json.RawMessage, context *lambda.Context) []byte {
	buf := new(bytes.Buffer)

	// decode event data
	req := new(request)
	json.Unmarshal(event, &req)

	// read a variable from LambdaContext
	fmt.Fprintf(buf, "Hello from %s. Current time: %v\n",
		context.FunctionName, time.Now().Format(time.Kitchen))

	// use some of the request data
	if len(req.Name) > 0 {
		fmt.Fprintf(buf, "Have a nice day, %s!\n", req.Name)
	}

	return buf.Bytes()
}

type request struct {
	Name string `json:"name"`
}
