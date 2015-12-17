package lambda

type Context struct {
	FunctionName       string
	FunctionVersion    string
	InvokedFunctionARN string
	MemoryLimit        int
	AWSRequestID       string
	LogGroupName       string
	LogStreamName      string
}
