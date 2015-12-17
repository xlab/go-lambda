package example

import (
	"fmt"
	"time"

	"github.com/xlab/go-lambda/lambda"
)

type Context lambda.Context

func Handler(c Context) string {
	return fmt.Sprintf("Hello from %s! Time is %v", c.FunctionName, time.Now().Format(time.Kitchen))
}
