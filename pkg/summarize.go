package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/fatih/color"
	"github.com/ingshtrom/dh-debug/pkg/types"
)

var summarizeFuncs = template.FuncMap{
	"red": func(arg int) string {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		return red(fmt.Sprint(arg))
	},
	"bold": func(arg string) string {
		bold := color.New(color.Bold).SprintFunc()
		return bold(arg)
	},
	"stringify": func(arg interface{}) string {
		return fmt.Sprint(arg)
	},
}

var summaryTemplate = `
  {{ "Tests Run:" | bold }} {{ .TestsRun }}
  {{ "Tests with Errors:" | bold }} {{ .Errors | red }}
  {{ "Exit Codes:" | bold }} {{ range $k, $v := .ExitCodes}}{{$k}}=>{{$v}} {{end}}
`

type SummarizeData struct {
	TestsRun  int
	Errors    int
	ExitCodes map[int]int
}

func SummarizeTestResults(testResultsFile string) {
	data, err := os.ReadFile(testResultsFile)
	if err != nil {
		fmt.Printf("error reading debug file: %v", err)
		os.Exit(1)
	}

	var dcs *[]types.DebugCommand

	err = json.Unmarshal(data, &dcs)
	if err != nil {
		fmt.Printf("error parsing debug file: %v", err)
		os.Exit(1)
	}

	testSummary := SummarizeData{
		ExitCodes: map[int]int{},
		TestsRun:  len(*dcs),
	}

	for _, dc := range *dcs {
		if dc.Errors != nil && len(dc.Errors) > 0 {
			testSummary.Errors = testSummary.Errors + 1
		}

		if val, ok := testSummary.ExitCodes[dc.ExitCode]; ok {
			testSummary.ExitCodes[dc.ExitCode] = val + 1
		} else {
			testSummary.ExitCodes[dc.ExitCode] = 1
		}
	}

	fmt.Println(testSummary)

	t, err := template.New("Summarize").Funcs(summarizeFuncs).Parse(summaryTemplate)
	if err != nil {
		fmt.Printf("error creating template: %v", err)
		os.Exit(1)
	}
	err = t.Execute(os.Stdout, testSummary)
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		os.Exit(1)
	}
}
