package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/apcera/termtables"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jawher/mow.cli"
)

var golambda = cli.App("go-lambda", "A tiny AWS Lambda manager for golang deployments.")

func init() {
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
	golambda.Command("list-regions", "Lists regions available for AWS Lambda service.", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			listRegions()
		}
	})
	golambda.Command("info", "Gets info about specific AWS Lambda function (specified by ID or NAME).",
		func(cmd *cli.Cmd) {
			cmd.Spec = "(ID | NAME)"
			idStr := cmd.StringArg("ID", "", "Specify function ID (as in `list` result).")
			cmd.StringArg("NAME", "", "Specify function NAME.")
			cmd.Action = func() {
				svc := lambda.New(session.New(&aws.Config{
					Region: aws.String(*region),
				}))
				id, _ := strconv.Atoi(*idStr)
				functionInfo(svc, id, *idStr, *region)
			}
		})
	if err := golambda.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

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

func listFunctions(svc *lambda.Lambda, region string, rx *regexp.Regexp) {
	resp, err := svc.ListFunctions(nil)
	if err != nil {
		log.Fatalln(err)
	} else if len(resp.Functions) == 0 {
		fmt.Println("0 lambda functions available.")
		return
	}
	table := termtables.CreateTable()
	table.UTF8Box()
	table.AddHeaders("#", "NAME", "UPDATED", "SIZE", "MEM LIMIT", "TIMEOUT", "DESCRIPTION")
	var filtered int
	for i, f := range resp.Functions {
		if rx != nil && !rx.MatchString(*f.FunctionName) {
			continue
		}
		filtered++
		table.AddRow(strconv.Itoa(i+1), *f.FunctionName,
			parseDate(*f.LastModified).Format(time.RFC822),
			fmt.Sprintf("%dB", *f.CodeSize),
			fmt.Sprintf("%dMB", *f.MemorySize),
			time.Second*time.Duration(*f.Timeout), *f.Description)
	}
	if filtered > 0 {
		table.AddTitle(fmt.Sprintf("AWS LAMBDA FUNCTIONS: %d (%s)", filtered, region))
		fmt.Println(table.Render())
		return
	}
	fmt.Printf("0 lambda functions matched `%s` pattern.\n", rx.String())
}

func parseDate(v string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02T15:04:05.999-0700", v, time.UTC)
	return t
}

func functionInfo(svc *lambda.Lambda, id int, name, region string) {
	resp, err := svc.ListFunctions(nil)
	if err != nil {
		log.Fatalln(err)
	} else if len(resp.Functions) == 0 {
		fmt.Println("error: 0 lambda functions available.")
		os.Exit(-1)
	}
	var f *lambda.FunctionConfiguration
	for i, function := range resp.Functions {
		if i+1 == id || *function.FunctionName == name {
			f = function
			name = *function.FunctionName
		}
	}
	if f == nil {
		switch {
		case id > 0:
			fmt.Printf("error: no such function with ID=%d found.\n", id)
		case len(name) > 0:
			fmt.Printf("error: no such function with NAME=%s found.\n", name)
		}
		os.Exit(-1)
	}

	table := termtables.CreateTable()
	table.UTF8Box()
	table.AddTitle(fmt.Sprintf("AWS LAMBDA FUNCTION %s (%s)", *f.FunctionName, region))
	table.AddRow("SHA256 Hash", *f.CodeSha256)
	table.AddRow("Code Size", fmt.Sprintf("%dB", *f.CodeSize))
	table.AddRow("Description", *f.Description)
	table.AddRow("Amazon Resource Name", *f.FunctionArn)
	table.AddRow("Handler", *f.Handler)
	table.AddRow("Last Modified", parseDate(*f.LastModified).Format(time.RFC822))
	table.AddRow("Memory Size", fmt.Sprintf("%dMB", *f.MemorySize))
	table.AddRow("Role", *f.Role)
	table.AddRow("Runtime", *f.Runtime)
	table.AddRow("Timeout", time.Second*time.Duration(*f.Timeout))
	table.AddRow("Version", *f.Version)
	fmt.Println(table.Render())
}
