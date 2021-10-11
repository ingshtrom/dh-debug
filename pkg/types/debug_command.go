package types

import (
	"fmt"
	"text/template"
	"os"


	"github.com/fatih/color"
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
  {{ "Command:" | red }} "{{ .Command }} {{ range $index, $element := .Arguments}}{{if $index}} {{end}}{{$element}}{{end}}"
  {{ "ExitCode:" | red }} {{ .ExitCode }}
  {{ "StdOut:" | red }}{{ .StdOut }}
  {{ "StdErr:" | red }}{{ .StdErr }}
  {{ "Processing Errors:" | red }} {{ range $index, $element := .Errors}}{{"\n"}}  - {{$element}}{{end}}
{{ "------------------------------------------" | green }}
`

type DebugCommand struct {
	Name string `json:"name"`
	Command string `json:"command"`
	Arguments []string `json:"arguments"`
	Errors []string  `json:"errors"`
	ExitCode int `json:"exit_code"`
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}


func (dc DebugCommand) Print() {
	t, err := template.New("DebugCommand").Funcs(funcs).Parse(debugCommandTemplate)
	if err != nil {
		fmt.Printf("error creating template: %v", err)
		os.Exit(1)
	}
	err = t.Execute(os.Stdout, dc)
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		os.Exit(1)
	}
}
