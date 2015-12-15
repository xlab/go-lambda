package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/apcera/termtables"
)

func listRegions() {
	table := termtables.CreateTable()
	table.UTF8Box()
	table.AddTitle("AWS LAMBDA REGIONS (2015-12-15)")
	table.AddRow("us-east-1", "US East (N. Virginia)")
	table.AddRow("us-west-2", "US West (Oregon)")
	table.AddRow("eu-west-1", "EU (Ireland)")
	table.AddRow("ap-northeast-1", "Asia Pacific (Tokyo)")
	fmt.Println(table.Render())
}

func getMain(packageName, packageFunc string) []byte {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, `from __future__ import print_function
# import json
import %s as pkg`, packageName)

	fmt.Fprintln(buf, "\n")
	fmt.Fprintf(buf, `def lol(event, context):
    c = pkg.Context()
    c.FunctionName = context.function_name
    c.FunctionVersion = context.function_version
    c.InvokedFunctionARN = context.invoked_function_arn
    # c.MemoryLimit = context.memory_limit_in_mb
    c.AWSRequestID = context.aws_request_id
    c.LogGroupName = context.log_group_name
    c.LogStreamName = context.log_stream_name
    result = pkg.%s(c)
    return result`, strings.Title(packageFunc))
	return buf.Bytes()
}
