package main

var moduleTemplate = `from __future__ import print_function
import {{.PackageName}} as pkg

def {{.FuncName}}(event, context):
    c = pkg.MakeContext()
    c.SetMemoryLimit(context.memory_limit_in_mb)
    c.FunctionName = context.function_name
    c.FunctionVersion = context.function_version
    c.InvokedFunctionARN = context.invoked_function_arn
    c.RequestID = context.aws_request_id
    c.LogGroupName = context.log_group_name
    c.LogStreamName = context.log_stream_name
    if context.identity != null:
        c.SetIdentity(context.identity.cognito_identity_id,
            context.identity.cognito_identity_pool_id)
    if context.client_context != null:
        c.SetClientContext(context.client_context.client.installation_id,
            context.client_context.client.app_title,
            context.client_context.client.app_version_name,
            context.client_context.client.app_version_code,
            context.client_context.client.app_package_name)
    result = pkg.{{.PackageFunc}}(passDict(event), c)
    return result

def passDict(dict):
    d = pkg.MakeDict(len(dict))
    if len(dict) == 0:
        return d
    options = {
        int: d.SetInt,
        long: d.SetInt64,
        float: d.SetFloat64,
        str: d.SetString,
        bool: d.SetBool,
    }
    for k, v in dict.iteritems():
        for opt in options:
            if isinstance(v, opt):
                options[opt](k, v)
    return d
`
