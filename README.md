# go-lambda [![aws](https://d0.awsstatic.com/logos/aws/AWS_Logo_PoweredBy_127px.png)](https://aws.amazon.com) [![Build Status](https://travis-ci.org/xlab/go-lambda.svg)](https://travis-ci.org/xlab/go-lambda)

`go-lambda` is a multi-purpose tool for creating and managing AWS Lambda instances backed by arbitrary Go code. Since there is no official support of Go, this tool automatically generates a wrappig module in Python 2.7 which is able to pass data back and forth to the Go land. Check out an [example](/example/example.go) function.

### Features at glance

* `<10ms` startup time, feels like native experience;
* Resulting `source.zip` size is only `1.0M` and in most cases has 2 files;
* Easy to use: start writing your own lambdas in Go just in few minutes;
* Relies on the official [AWS SDK for Go](https://github.com/aws/aws-sdk-go) while making all the requests, security is guaranteed;
* No any boilerplate or "all-in-one" aims: the tool is made to do its job and nothing else. Yes, this is also a feature.

### Installation

```
$ go get github.com/xlab/go-lambda
```

Keep in mind that if you are using something other than `Linux x86_64`, this tool will invoke the [xlab/go-lambda](https://hub.docker.com/r/xlab/go-lambda/) Docker image in order to create cross-platform modules for AWS. Prepare your Docker setup beforehand.

### Reference

```
Usage: go-lambda [OPTIONS] COMMAND [arg...]

A tiny AWS Lambda manager for Go deployments.

Options:
  --verbose=false            Run in verbose mode.
  -v, --version=false        Show the version and exit
  -r, --region="eu-west-1"   Specify the region.

Commands:
  list         Lists all defined AWS Lambda functions in the region.
  regions      Lists regions available for AWS Lambda service.
  info         Gets info about specific AWS Lambda function (specified by ID or NAME).
  source       Gets .zip of an AWS Lambda function source (specified by ID or NAME).
  delete       Deletes an AWS Lambda function (specified by ID or NAME).
  create       Creates an AWS Lambda function and uploads specified Go package as its source.
  update       Updates source of an AWS Lambda function (specified by ID or NAME).

Run 'go-lambda COMMAND --help' for more information on a command.
```

### Examples

#### List available lambda functions in default region:

```
$ ./go-lambda list
╭──────────────────────────────────────────────────────────────────────────────────────────╮
│                           AWS LAMBDA FUNCTIONS: 2 (eu-west-1)                            │
├───┬─────────────────┬─────────────────────┬──────────┬───────────┬─────────┬─────────────┤
│ # │ NAME            │ UPDATED             │ SIZE     │ MEM LIMIT │ TIMEOUT │ DESCRIPTION │
├───┼─────────────────┼─────────────────────┼──────────┼───────────┼─────────┼─────────────┤
│ 1 │ xxx-handler     │ 20 Dec 15 11:24 UTC │ 459.22K  │ 128M      │ 3s      │             │
│ 2 │ yyy-handler     │ 19 Dec 15 09:20 UTC │ 384.10K  │ 128M      │ 3s      │             │
╰───┴─────────────────┴─────────────────────┴──────────┴───────────┴─────────┴─────────────╯
```

#### Create a new lambda backed by Go code:

```
$ go-lambda create --role arn:aws:iam::account-id:role/lambda_basic_execution handler github.com/xlab/go-lambda/example

╭───────────────────────────────────────────────────────────────────────────────────────╮
│                    AWS LAMBDA FUNCTION example-handler (eu-west-1)                    │
├──────────────────────┬────────────────────────────────────────────────────────────────┤
│ SHA256 Hash          │ rcUmH7y39GVzmnwEbeV4aCj9e4vkPdoZIF14oj+nyWM=                   │
│ Code Size            │ 1067.95K                                                       │
│ Description          │                                                                │
│ Amazon Resource Name │ arn:aws:lambda:eu-west-1:account-id:function:example-handler   │
│ Handler              │ example-handler.handler                                        │
│ Last Modified        │ 20 Dec 15 14:00 UTC                                            │
│ Memory Size          │ 128M                                                           │
│ Role                 │ arn:aws:iam::account-id:role/lambda_basic_execution            │
│ Runtime              │ python2.7                                                      │
│ Timeout              │ 3s                                                             │
│ Version              │ 25                                                             │
╰──────────────────────┴────────────────────────────────────────────────────────────────╯
```

Now you can manage your function in the AWS dashboard, create API endpoints and event handlers.

#### Update existing function after fixing its Go code:

```
$ go-lambda update example-handler handler github.com/xlab/go-lambda/example
```

### Roadmap

- [ ] Replace JSON marshalling with something efficient (maybe);
- [ ] Try to convince AWS staff that we need the native Go support in their cloud.

Please, spread the word about `go-lambda` — let people use their favourite language for AWS Lambda. 🍻

![go-lambda](http://cl.ly/1w1U1n3w3W2n/go-lamda-alt.png)

### Useful links

* [AWS Serverless Multi-Tier Architectures Whitepaper](https://d0.awsstatic.com/whitepapers/AWS_Serverless_Multi-Tier_Architectures.pdf)
* [Lambda limits](http://docs.aws.amazon.com/lambda/latest/dg/limits.html)
* [The Twelve Days of Lambda](https://aws.amazon.com/blogs/compute/the-twelve-days-of-lambda/)
* [GoSparta Project Limitations](http://gosparta.io/docs/limitations/)

### License

[MIT](/LICENSE)
