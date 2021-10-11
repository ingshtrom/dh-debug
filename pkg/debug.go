package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/ingshtrom/dh-debug/pkg/types"
)

func RunDebugTests(configFile, filePath string) {
	rawDCS, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		os.Exit(1)
	}

	var dcs *[]types.DebugCommand
	err = json.Unmarshal(rawDCS, &dcs)
	if err != nil {
		fmt.Printf("error parsing config file: %v", err)
		os.Exit(1)
	}

	//dcs := []types.DebugCommand{
	//  // DNS
	//  {
	//    Name:      "Docker Hub DNS A",
	//    Command:   "dig",
	//    Arguments: []string{"hub.docker.com", "A"},
	//  },
	//  {
	//    Name:      "Docker Hub DNS AAAA",
	//    Command:   "dig",
	//    Arguments: []string{"hub.docker.com", "AAAA"},
	//  },
	//  {
	//    Name:      "Docker Registry DNS A",
	//    Command:   "dig",
	//    Arguments: []string{"registry-1.docker.io", "A"},
	//  },
	//  {
	//    Name:      "Docker Registry DNS AAAA",
	//    Command:   "dig",
	//    Arguments: []string{"registry-1.docker.io", "AAAA"},
	//  },
	//  {
	//    Name:      "Docker Registry Auth DNS A",
	//    Command:   "dig",
	//    Arguments: []string{"auth.docker.io", "A"},
	//  },
	//  {
	//    Name:      "Docker Registry Auth DNS AAAA",
	//    Command:   "dig",
	//    Arguments: []string{"auth.docker.io", "AAAA"},
	//  },
	//  {
	//    Name:      "Cloudflare DNS A",
	//    Command:   "dig",
	//    Arguments: []string{"production.cloudflare.docker.com", "A"},
	//  },
	//  {
	//    Name:      "Cloudflare DNS AAAA",
	//    Command:   "dig",
	//    Arguments: []string{"production.cloudflare.docker.com", "AAAA"},
	//  },

	//  // CURL
	//  {
	//    Name:      "Docker Registry IPv4",
	//    Command:   "curl",
	//    Arguments: []string{"-4svo /dev/null https://registry-1.docker.io/"},
	//  },
	//  {
	//    Name:      "Docker Registry IPv6",
	//    Command:   "curl",
	//    Arguments: []string{"-6svo /dev/null https://registry-1.docker.io/"},
	//  },
	//  {
	//    Name:      "Docker Hub IPv4",
	//    Command:   "curl",
	//    Arguments: []string{"-4svo /dev/null https://hub.docker.com/"},
	//  },
	//  {
	//    Name:      "Docker Hub IPv6",
	//    Command:   "curl",
	//    Arguments: []string{"-6svo /dev/null https://hub.docker.com/"},
	//  },
	//  {
	//    Name:      "Docker Registry Auth IPv4",
	//    Command:   "curl",
	//    Arguments: []string{"-4svo /dev/null https://auth.docker.io/"},
	//  },
	//  {
	//    Name:      "Docker Registry Auth IPv6",
	//    Command:   "curl",
	//    Arguments: []string{"-6svo /dev/null https://auth.docker.io/"},
	//  },
	//  {
	//    Name:      "Cloudflare IPv4",
	//    Command:   "curl",
	//    Arguments: []string{"-4svo /dev/null https://production.cloudflare.docker.com/"},
	//  },
	//  {
	//    Name:      "Cloudflare IPv6",
	//    Command:   "curl",
	//    Arguments: []string{"-6svo", "/dev/null", "https://production.cloudflare.docker.com/"},
	//  },


	//  // Cloudflare Trace
	//  {
	//    Name: "Cloudflare Trace IPv4",
	//    Command: "curl",
	//    Arguments: []string{"-4sv", "http://production.cloudflare.docker.com/cdn-cgi/trace"},

	//  },
	//  {
	//    Name: "Cloudflare Trace IPv4",
	//    Command: "curl",
	//    Arguments: []string{"-6sv", "http://production.cloudflare.docker.com/cdn-cgi/trace", },

	//  },
	//}

	out := make([]types.DebugCommand, 0)

	for _, dc := range *dcs {
		out = append(out, runCommand(dc))
	}

	saveDebugData(out, filePath)
}

func runCommand(dc types.DebugCommand) types.DebugCommand {
	dc = addExtraArgs(dc)
	cmd := exec.Command(dc.Command, dc.Arguments...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		dc.Errors = append(dc.Errors, err.Error())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		dc.Errors = append(dc.Errors, err.Error())
	}

	err = cmd.Start()
	if err != nil {
		dc.Errors = append(dc.Errors, err.Error())
		return dc
	}

	slurpOut, err := io.ReadAll(stdout)
	if err != nil {
		dc.Errors = append(dc.Errors, err.Error())
	} else {
		dc.StdOut = string(slurpOut)
	}

	slurpErr, err := io.ReadAll(stderr)
	if err != nil {
		dc.Errors = append(dc.Errors, err.Error())
	} else {
		dc.StdErr = string(slurpErr)
	}

	err = cmd.Wait()
	if err != nil {
		dc.Errors = append(dc.Errors, err.Error())
	}

	dc.ExitCode = cmd.ProcessState.ExitCode()

	return dc
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
