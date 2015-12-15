package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type deploymentConfig struct {
	HandlerName       string
	LambdaName        string
	LambdaDescription string
	MemorySize        int
	Role              string
	Timeout           int
	SourceZip         []byte
}

func createFunction(svc *lambda.Lambda, cfg *deploymentConfig, region string) {
	input := &lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ZipFile: cfg.SourceZip,
		},
		Description:  aws.String(cfg.LambdaDescription),
		FunctionName: aws.String(cfg.LambdaName),
		Handler:      aws.String(cfg.HandlerName),
		MemorySize:   aws.Int64(int64(cfg.MemorySize)),
		Publish:      aws.Bool(true),
		Role:         aws.String(cfg.Role),
		Runtime:      aws.String(lambda.RuntimePython27),
		Timeout:      aws.Int64(int64(cfg.Timeout)),
	}
	f, err := svc.CreateFunction(input)
	if err != nil {
		log.Fatalln(err)
	}
	functionInfo(f, region)
}

func updateFunction(svc *lambda.Lambda, name, region string, sourceZip []byte) {
	input := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(name),
		Publish:      aws.Bool(true),
		ZipFile:      sourceZip,
	}
	f, err := svc.UpdateFunctionCode(input)
	if err != nil {
		log.Fatalln(err)
	}
	functionInfo(f, region)
}

func deleteFunction(svc *lambda.Lambda, name, version string) {
	input := &lambda.DeleteFunctionInput{
		FunctionName: aws.String(name),
	}
	if len(version) > 0 {
		input.Qualifier = aws.String(version)
	}
	if _, err := svc.DeleteFunction(input); err != nil {
		log.Fatalln(err)
	}
}
