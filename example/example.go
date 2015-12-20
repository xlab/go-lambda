// Package example shows how an AWS Lambda function can be implemented using Go.
//
// Building: `go-lambda update 1 handler github.com/xlab/go-lambda/example`
package example

import (
	"bytes"
	"fmt"
	"time"

	"github.com/xlab/go-lambda/lambda"
)

// Handler will be called by a Python module wrapping this package.
var Handler = lambda.Use(lambda.HandlerFunc(exampleHandler))

func exampleHandler(event lambda.Dict, context *lambda.Context) []byte {
	buf := new(bytes.Buffer)

	// Use some variable from lambda context
	fmt.Fprintf(buf, "Hello from %s. Current time: %v\n",
		context.FunctionName, time.Now().Format(time.Kitchen))

	// Read some field from event data
	if name, ok := event["name"].(string); ok {
		fmt.Fprintf(buf, "Have a nice day, %s!\n", name)
	}

	return buf.Bytes()
}
