package types

type DebugCommand struct {
	Arguments      []string `json:"arguments"`
	Command        string   `json:"command"`
	Errors         []string `json:"errors"`
	ExitCode       int      `json:"exitCode"`
	ExtraArguments []string `json:"extraArguments"`
	Filter         string   `json:"filter"`
	Name           string   `json:"name"`
	StdOut         string   `json:"stdout"`
	StdErr         string   `json:"stderr"`
}
