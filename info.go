package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/apcera/termtables"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/cheggaaa/pb"
)

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
			fmt.Sprintf("%.2fK", float32(*f.CodeSize)/1024.0),
			fmt.Sprintf("%dM", *f.MemorySize),
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

func findFuction(svc *lambda.Lambda, id int, name string) *lambda.FunctionConfiguration {
	listResp, err := svc.ListFunctions(nil)
	if err != nil {
		log.Fatalln(err)
	} else if len(listResp.Functions) == 0 {
		fmt.Println("error: 0 lambda functions available.")
		os.Exit(-1)
	}
	for i, f := range listResp.Functions {
		if i+1 == id || *f.FunctionName == name {
			return f
		}
	}
	switch {
	case id > 0:
		fmt.Printf("error: no such function with ID=%d found.\n", id)
	case len(name) > 0:
		fmt.Printf("error: no such function with NAME=%s found.\n", name)
	}
	os.Exit(-1)
	return nil
}

func functionSource(svc *lambda.Lambda, f *lambda.FunctionConfiguration, path string, onlyURL bool) {
	getResp, err := svc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: f.FunctionName,
		Qualifier:    f.Version,
	})
	if err != nil {
		log.Println("error:", err)
		return
	}
	if onlyURL {
		log.Println(*getResp.Code.Location)
		return
	}
	resp, err := http.Get(*getResp.Code.Location)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	bar := pb.New(int(resp.ContentLength)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	barproxy := bar.NewProxyReader(resp.Body)
	bar.ShowSpeed = true
	bar.Start()
	io.Copy(out, barproxy)
	out.Close()
	bar.Finish()
}

func functionInfo(f *lambda.FunctionConfiguration, region string) {
	table := termtables.CreateTable()
	table.UTF8Box()
	table.AddTitle(fmt.Sprintf("AWS LAMBDA FUNCTION %s (%s)", *f.FunctionName, region))
	table.AddRow("SHA256 Hash", *f.CodeSha256)
	table.AddRow("Code Size", fmt.Sprintf("%.2fK", float32(*f.CodeSize)/1024.0))
	table.AddRow("Description", *f.Description)
	table.AddRow("Amazon Resource Name", *f.FunctionArn)
	table.AddRow("Handler", *f.Handler)
	table.AddRow("Last Modified", parseDate(*f.LastModified).Format(time.RFC822))
	table.AddRow("Memory Size", fmt.Sprintf("%dM", *f.MemorySize))
	table.AddRow("Role", *f.Role)
	table.AddRow("Runtime", *f.Runtime)
	table.AddRow("Timeout", time.Second*time.Duration(*f.Timeout))
	table.AddRow("Version", *f.Version)
	fmt.Println(table.Render())
}

func listRegions() {
	table := termtables.CreateTable()
	table.UTF8Box()
	table.AddTitle("AWS LAMBDA REGIONS (2015-12-20)")
	table.AddRow("us-east-1", "US East (N. Virginia)")
	table.AddRow("us-west-2", "US West (Oregon)")
	table.AddRow("eu-west-1", "EU (Ireland)")
	table.AddRow("ap-northeast-1", "Asia Pacific (Tokyo)")
	fmt.Println(table.Render())
}
