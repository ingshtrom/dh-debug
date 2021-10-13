package pkg

import (
	"testing"

	"github.com/ingshtrom/dh-debug/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDebugGrepFilter(t *testing.T) {
	tests := []struct {
		Name           string
		Filter         string
		ExpectedOutput string
	}{
		{
			Name:   "WorksCustom",
			Filter: "custom",
			ExpectedOutput: `custom_black=#2A3236
custom_blue=#A7DAF8
custom_comment_grey=#546E7A
custom_cyan=#64FCDA
custom_green=#5CF19E
custom_orange=#F2CD86
custom_red=#FF5252
custom_visual_grey=#4B5962
custom_white=#EDEFF1
custom_yellow=#FFD740`,
		},
		{
			Name:   "WorksProxyIgnoreCase",
			Filter: "proxy",
			ExpectedOutput: `http_proxy=http:127.0.0.1
hTtP_ProXy2=http:127.0.0.1`,
		},
	}

	input := `WINDOWID=21525
XPC_FLAGS=0x0
XPC_SERVICE_NAME=0
ZSH=/Users/alexhokanson/.oh-my-zsh
ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE=fg=blue,bold
ZSH_AUTOSUGGEST_USE_ASYNC=true
ZSH_THEME=agnoster
__CFBundleIdentifier=net.kovidgoyal.kitty
__CF_USER_TEXT_ENCODING=0x1F5:0x0:0x0
custom_black=#2A3236
custom_blue=#A7DAF8
custom_comment_grey=#546E7A
custom_cyan=#64FCDA
custom_green=#5CF19E
custom_orange=#F2CD86
custom_red=#FF5252
custom_visual_grey=#4B5962
custom_white=#EDEFF1
custom_yellow=#FFD740
http_proxy=http:127.0.0.1
hTtP_ProXy2=http:127.0.0.1
_=/usr/local/bin/go`

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.ExpectedOutput, grepFilter(input, test.Filter))
		})
	}
}

func TestDebugAddExtraArgs(t *testing.T) {
	tests := []struct {
		Name           string
		DebugCommand   types.DebugCommand
		ExpectedArguments []string
		ExpectedExtraArguments []string
	}{
		{
			Name:   "Curl",
			DebugCommand: types.DebugCommand{
				Name: "Curl",
				Command: "curl",
				Arguments: []string{"-v"},
			},
			ExpectedArguments: []string{"-v"},
			ExpectedExtraArguments: []string{"-L", "-w", "\"HTTP Connect:%{http_connect} \nNumber Connects: %{num_connects} \nNumber Redirects: %{num_redirects} \nRedirect URL: %{redirect_url} \nSize Download: %{size_download} \nSize Upload: %{size_upload} \nSSL Verify: %{ssl_verify_result} \nTime Handshake: %{time_appconnect} \nTime Connect: %{time_connect} \nName Lookup Time: %{time_namelookup} \nTime Pretransfer: %{time_pretransfer} \nTime Redirect: %{time_redirect} \nTime Start Transfer: %{time_starttransfer} \nTime Total: %{time_total} \nEffective URL: %{url_effective}\""},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			dc := addExtraArgs(test.DebugCommand)
			assert.Equal(t, test.ExpectedArguments, dc.Arguments)
			assert.Equal(t, test.ExpectedExtraArguments, dc.ExtraArguments)
		})
	}
}
