package lambda

import "strconv"

// Context defines LambdaContext object as described on
// http://docs.aws.amazon.com/lambda/latest/dg/python-context-object.html
type Context struct {
	FunctionName       string           `lambda:"function_name"`
	FunctionVersion    string           `lambda:"function_version"`
	InvokedFunctionARN string           `lambda:"invoked_function_arn"`
	MemoryLimit        int              `lambda:"memory_limit_in_mb"`
	RequestID          string           `lambda:"aws_request_id"`
	LogGroupName       string           `lambda:"log_group_name"`
	LogStreamName      string           `lambda:"log_stream_name"`
	Identity           *CognitoIdentity `lambda:"identity"`
	ClientContext      *ClientContext   `lambda:"client_context"`
}

func (c *Context) SetMemoryLimit(str string) error {
	c.MemoryLimit, _ = strconv.Atoi(str)
	return nil
}

func (c *Context) SetCognitoIdentity(id, poolID string) error {
	c.Identity = &CognitoIdentity{
		ID: id, PoolID: poolID,
	}
	return nil
}

func (c *Context) SetClientContext(id, title, versionName, versionCode, packageName string) error {
	c.ClientContext = &ClientContext{
		InstallationID: id,
		AppTitle:       title,
		AppVersionName: versionName,
		AppVersionCode: versionCode,
		AppPackageName: packageName,
	}
	return nil
}

// CognitoIdentity identity provider when invoked through the AWS Mobile SDK.
type CognitoIdentity struct {
	ID     string `lambda:"cognito_identity_id"`
	PoolID string `lambda:"cognito_identity_pool_id"`
}

// ClientContext holds information about the client application and device when invoked through the AWS Mobile SDK.
type ClientContext struct {
	InstallationID string `lambda:"installation_id"`
	AppTitle       string `lambda:"app_title"`
	AppVersionName string `lambda:"app_version_name"`
	AppVersionCode string `lambda:"app_version_code"`
	AppPackageName string `lambda:"app_package_name"`

	Custom Dict `lambda:"custom"`
	Env    Dict `lambda:"env"`
}
