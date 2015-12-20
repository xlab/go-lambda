package lambda

import "strconv"

// Context defines LambdaContext object as described on
// http://docs.aws.amazon.com/lambda/latest/dg/python-context-object.html
type Context struct {
	FunctionName       string           `json:"function_name"`
	FunctionVersion    string           `json:"function_version"`
	InvokedFunctionARN string           `json:"invoked_function_arn"`
	MemoryLimitString  string           `json:"memory_limit_in_mb"`
	RequestID          string           `json:"aws_request_id"`
	LogGroupName       string           `json:"log_group_name"`
	LogStreamName      string           `json:"log_stream_name"`
	Identity           *CognitoIdentity `json:"identity"`
	ClientContext      *ClientContext   `json:"client_context"`
}

func (c *Context) MemoryLimit() int {
	i, _ := strconv.Atoi(c.MemoryLimitString)
	return i
}

// CognitoIdentity identity provider when invoked through the AWS Mobile SDK.
type CognitoIdentity struct {
	ID     string `json:"cognito_identity_id"`
	PoolID string `json:"cognito_identity_pool_id"`
}

// ClientContext holds information about the client application and device when invoked through the AWS Mobile SDK.
type ClientContext struct {
	InstallationID string `json:"installation_id"`
	AppTitle       string `json:"app_title"`
	AppVersionName string `json:"app_version_name"`
	AppVersionCode string `json:"app_version_code"`
	AppPackageName string `json:"app_package_name"`

	Custom Dict `json:"custom"`
	Env    Dict `json:"env"`
}
