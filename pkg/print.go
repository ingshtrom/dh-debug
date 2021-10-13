package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/fatih/color"
	"github.com/ingshtrom/dh-debug/pkg/types"
)

var funcs = template.FuncMap{
	"red": func (arg string) string {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		return red(arg)
	},
	"blue": func (arg string) string {
		blue := color.New(color.FgBlue, color.Bold).SprintFunc()
		return blue(arg)
	},
	"green": func (arg string) string {
		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		return green(arg)
	},
}
var debugCommandTemplate = `
{{ "------------------------------------------" | green }}
### {{ .Name | red }} ###
  {{ "Command:" | red }} "{{ .Command }} {{ range $index, $element := .Arguments}}{{if $index}} {{end}}{{$element}}{{end}} {{ range $index, $element := .ExtraArguments}}{{if $index}} {{end}}{{$element}}{{end}}"
  {{ "Filter:" | red }} "{{ .Filter }}"
  {{ "ExitCode:" | red }} {{ .ExitCode }}
  {{ "StdOut:" | red }}{{ .StdOut }}
  {{ "StdErr:" | red }}{{ .StdErr }}
  {{ "Processing Errors:" | red }} {{ range $index, $element := .Errors}}{{"\n"}}  - {{$element}}{{end}}
{{ "------------------------------------------" | green }}
`

func PrintDebugTests(testResultsFile string, errorsOnly bool) {
	fmt.Println(testResultsFile, errorsOnly)
	data, err := os.ReadFile(testResultsFile)
	if err != nil {
		fmt.Printf("could not fine results file at %q\n", testResultsFile)
		os.Exit(1)
	}

	var dcs *[]types.DebugCommand

	err = json.Unmarshal(data, &dcs)
	if err != nil {
		fmt.Println("error parsing debug file: ", err)
		os.Exit(1)
	}

	for _, dc := range *dcs {
		if !errorsOnly {
			printDebugCommand(dc)
		} else if dc.ExitCode != 0 || len(dc.Errors) > 0 {
			printDebugCommand(dc)
		}
	}
}

func printDebugCommand(dc types.DebugCommand) {
	t, err := template.New("DebugCommand").Funcs(funcs).Parse(debugCommandTemplate)
	if err != nil {
		fmt.Println("error creating template: ", err)
		os.Exit(1)
	}
	err = t.Execute(os.Stdout, dc)
	if err != nil {
		fmt.Println("error executing template: ", err)
		os.Exit(1)
	}
}

