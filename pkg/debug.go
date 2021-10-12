package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ingshtrom/dh-debug/pkg/types"
)

func RunDebugTests(configFile, filePath string) {
	rawDCS, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		os.Exit(1)
	}

	var config *types.Config
	err = json.Unmarshal(rawDCS, &config)
	if err != nil {
		fmt.Printf("error parsing config file: %v", err)
		os.Exit(1)
	}

	out := make([]types.DebugCommand, 0)

	for _, dc := range *&config.ShellTests {
		out = append(out, runCommand(dc))
	}

	saveDebugData(out, filePath)
}

func runCommand(dc types.DebugCommand) types.DebugCommand {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dc = addExtraArgs(dc)
	cmd := exec.CommandContext(ctx, dc.Command, dc.Arguments...)

	fmt.Printf("$ %s %s ...", dc.Command, strings.ReplaceAll(strings.Join(dc.Arguments, " "), "\n", "\\n"))

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		dc.Errors = append(dc.Errors, err.Error())
		fmt.Printf("❌ Failed to Run\n")
		return dc
	}

	err = cmd.Wait()
	if err != nil {
		dc.Errors = append(dc.Errors, err.Error())
	}

	dc.ExitCode = cmd.ProcessState.ExitCode()
	dc.StdErr = string(stderr.Bytes())

	stdoutData := string(stdout.Bytes())
	if dc.Filter != "" {
		stdoutData = grepFilter(stdoutData, dc.Filter)
	}
	dc.StdOut = stdoutData

	if ctx.Err() == context.DeadlineExceeded {
		dc.Errors = append(dc.Errors, ctx.Err().Error())
		fmt.Printf("❌ Timeout\n")
		return dc
	}

	fmt.Printf("✅\n")

	return dc
}

func grepFilter(stdout, filter string) string {
	lines := strings.Split(stdout, "\n")
	linesToKeep := make([]string, 0)
	for _, l := range lines {
		if strings.Contains(strings.ToLower(l), filter) {
			linesToKeep = append(linesToKeep, l)
		}
	}
	return strings.Join(linesToKeep, "\n")
}

func addExtraArgs(dc types.DebugCommand) types.DebugCommand {
	if dc.Command == "curl" {
		dc.Arguments = append(dc.Arguments, "-w \"Content Type: %{content_type} \nHTTP Code: %{http_code} \nHTTP Connect:%{http_connect} \nNumber Connects: %{num_connects} \nNumber Redirects: %{num_redirects} \nRedirect URL: %{redirect_url} \nSize Download: %{size_download} \nSize Upload: %{size_upload} \nSSL Verify: %{ssl_verify_result} \nTime Handshake: %{time_appconnect} \nTime Connect: %{time_connect} \nName Lookup Time: %{time_namelookup} \nTime Pretransfer: %{time_pretransfer} \nTime Redirect: %{time_redirect} \nTime Start Transfer: %{time_starttransfer} \nTime Total: %{time_total} \nEffective URL: %{url_effective}\"")
	}

	return dc
}

func saveDebugData(data []types.DebugCommand, filePath string) error {
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	err = os.WriteFile(filePath, rawData, 0777)
	if err != nil {
		return err
	}

	return nil
}
