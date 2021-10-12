package types

type DebugCommand struct {
	Arguments []string `json:"arguments"`
	Name      string   `json:"name"`
	Command   string   `json:"command"`
	Filter    string   `json:"filter"`
	Errors    []string `json:"errors"`
	ExitCode  int      `json:"exit_code"`
	StdOut    string   `json:"stdout"`
	StdErr    string   `json:"stderr"`
}

type Config struct {
	ShouldTestDockerPull bool           `json:"should_test_docker_pull"`
	ShellTests           []DebugCommand `json:"shell_tests"`
}
