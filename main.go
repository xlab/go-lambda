package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jawher/mow.cli"
)

var golambda = cli.App("go-lambda", "A tiny AWS Lambda manager for Go deployments.")
var debug = golambda.BoolOpt("verbose", false, "Run in verbose mode.")

func init() {
	log.SetFlags(log.Lshortfile)
	golambda.Version("v version", "0.1")
}

func main() {
	region := golambda.StringOpt("r region", "eu-west-1", "Specify the region.")
	golambda.Command("list", "Lists all defined AWS Lambda functions in the region.", func(cmd *cli.Cmd) {
		filter := cmd.StringOpt("f filter", "", "Filter names by regexp.")
		cmd.Action = func() {
			rx, err := regexp.Compile(*filter)
			if len(*filter) > 0 && err != nil {
				log.Fatalln(err)
			}
			svc := lambda.New(session.New(&aws.Config{
				Region: aws.String(*region),
			}))
			listFunctions(svc, *region, rx)
		}
	})
	golambda.Command("regions", "Lists regions available for AWS Lambda service.", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			listRegions()
		}
	})
	golambda.Command("gopy", "Generates Python bindings for a Go package using a dockerized tool from gopy project.",
		func(cmd *cli.Cmd) {
			pkg := cmd.StringArg("PACKAGE", "", "Fully qualified Go package name.")
			cmd.Action = func() {
				runGopy(*pkg)
			}
		})
	golambda.Command("info", "Gets info about specific AWS Lambda function (specified by ID or NAME).",
		func(cmd *cli.Cmd) {
			cmd.Spec = "(ID | NAME)"
			idStr := cmd.StringArg("ID", "", "Function ID (as in `list` result).")
			cmd.StringArg("NAME", "", "Function NAME.")
			cmd.Action = func() {
				svc := lambda.New(session.New(&aws.Config{
					Region: aws.String(*region),
				}))
				id, _ := strconv.Atoi(*idStr)
				functionInfo(findFuction(svc, id, *idStr), *region)
			}
		})
	golambda.Command("source", "Gets .zip of an AWS Lambda function source (specified by ID or NAME).",
		func(cmd *cli.Cmd) {
			cmd.Spec = "(ID | NAME)"
			idStr := cmd.StringArg("ID", "", "Function ID (as in `list` result).")
			cmd.StringArg("NAME", "", "Function NAME.")
			source := cmd.StringOpt("o out", "source.zip", "Specify path of output file.")
			onlyURL := cmd.BoolOpt("u url", false, "Show file URL only (valid for 10m), do not download.")
			cmd.Action = func() {
				svc := lambda.New(session.New(&aws.Config{
					Region: aws.String(*region),
				}))
				id, _ := strconv.Atoi(*idStr)
				f := findFuction(svc, id, *idStr)
				functionSource(svc, f, *source, *onlyURL)
			}
		})
	golambda.Command("delete", "Deletes an AWS Lambda function (specified by ID or NAME).",
		func(cmd *cli.Cmd) {
			cmd.Spec = "(ID | NAME)"
			idStr := cmd.StringArg("ID", "", "Function ID (as in `list` result).")
			cmd.StringArg("NAME", "", "Function NAME.")
			version := cmd.StringOpt("v version", "", "Function version (qualifier) to delete.")
			cmd.Action = func() {
				svc := lambda.New(session.New(&aws.Config{
					Region: aws.String(*region),
				}))
				id, _ := strconv.Atoi(*idStr)
				f := findFuction(svc, id, *idStr)
				deleteFunction(svc, *f.FunctionName, *version)
			}
		})
	golambda.Command("create", "Creates an AWS Lambda function and uploads specified Go package as its source.",
		func(cmd *cli.Cmd) {
			cmd.Spec = "[OPTIONS] FUNC PACKAGE [FILES...]"
			packageFunc := cmd.StringArg("FUNC", "", "Function name in Go package to be called from within Python wrapper.")
			packagePath := cmd.StringArg("PACKAGE", ".", "Fully qualified import path of a Go package.")
			additionalFiles := cmd.StringsArg("FILES", nil, "A list of additional static files to be included in archive.")
			funcName := cmd.StringOpt("n name", "package-func", "The name you want to assign to the function you are uploading, should be Amazon Resource Name (ARN) or unqualified.")
			funcDescription := cmd.StringOpt("d description", "", "A short, user-defined function description. Lambda does not use this value.")
			memorySize := cmd.IntOpt("m memsize", 128, "The amount of memory, in MB, your Lambda function is given. The value must be a multiple of 64MB.")
			role := cmd.StringOpt("r role", "arn:aws:iam::account-id:role/lambda_basic_execution", "The Amazon Resource Name (ARN) of the IAM role that Lambda assumes when it executes your function.")
			timeout := cmd.IntOpt("t timeout", 3, "The function execution time at which Lambda should terminate the function.")
			writeZip := cmd.StringOpt("w write-zip", "", "Path to write the produced .zip of an AWS Lambda function source.")
			dryRun := cmd.BoolOpt("dry", false, "Run in dry mode, do not actually upload anything, but all the processing will be done.")

			cmd.Action = func() {
				tmp := getTempDir()
				buildModuleBridge(tmp, *packagePath, *packageFunc)
				packageName := packageName(*packagePath)
				if *funcName == "package-func" {
					*funcName = fmt.Sprintf("%s-%s", packageName, *packageFunc)
				}
				module := getModuleSource(packageName, *packageFunc)
				zip := makeZip(module, "module.py", "module.so", *additionalFiles...)
				if len(*writeZip) > 0 {
					if err := ioutil.WriteFile(*writeZip, zip, 0644); err != nil {
						log.Fatalln(err)
					}
				}

				if *dryRun {
					return
				}
				svc := lambda.New(session.New(&aws.Config{
					Region: aws.String(*region),
				}))
				handler := fmt.Sprintf("%s.%s", *funcName, *packageFunc)
				cfg := &deploymentConfig{
					HandlerName:       handler,
					LambdaName:        *funcName,
					LambdaDescription: *funcDescription,
					MemorySize:        *memorySize,
					Role:              *role,
					Timeout:           *timeout,
					SourceZip:         zip,
				}
				createFunction(svc, cfg, *region)
			}
		})
	golambda.Command("update", "Updates source of an AWS Lambda function (specified by ID or NAME).",
		func(cmd *cli.Cmd) {
			cmd.Spec = "[OPTIONS] (ID | NAME) FUNC PACKAGE [FILES...]"
			idStr := cmd.StringArg("ID", "", "Function ID (as in `list` result).")
			cmd.StringArg("NAME", "", "Function NAME.")
			packageFunc := cmd.StringArg("FUNC", "", "Function name in Go package to be called from within Python wrapper.")
			packagePath := cmd.StringArg("PACKAGE", ".", "Fully qualified import path of a Go package.")
			additionalFiles := cmd.StringsArg("FILES", nil, "A list of additional static files to be included in archive.")
			funcName := cmd.StringOpt("n name", "package-func", "The name you want to assign to the function you are uploading, should be Amazon Resource Name (ARN) or unqualified.")
			writeZip := cmd.StringOpt("w write-zip", "", "Path to write the produced .zip of an AWS Lambda function source.")
			dryRun := cmd.BoolOpt("dry", false, "Run in dry mode, do not actually upload anything, but all the processing will be done.")

			cmd.Action = func() {
				tmp := getTempDir()
				buildModuleBridge(tmp, *packagePath, *packageFunc)
				packageName := packageName(*packagePath)
				if *funcName == "package-func" {
					*funcName = fmt.Sprintf("%s-%s", packageName, *packageFunc)
				}
				module := getModuleSource(packageName, *packageFunc)
				zip := makeZip(module, "module.py", "module.so", *additionalFiles...)
				if len(*writeZip) > 0 {
					if err := ioutil.WriteFile(*writeZip, zip, 0644); err != nil {
						log.Fatalln(err)
					}
				}
				if *dryRun {
					return
				}
				svc := lambda.New(session.New(&aws.Config{
					Region: aws.String(*region),
				}))
				id, _ := strconv.Atoi(*idStr)
				f := findFuction(svc, id, *idStr)
				updateFunction(svc, *f.FunctionName, *region, zip)
			}
		})
	if err := golambda.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
