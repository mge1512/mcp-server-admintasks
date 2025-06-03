package utils

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/server"
)

var utilsDebug = false
var AdminTasksMCPServer *server.MCPServer

type RunningMode int
const (
	Test RunningMode = iota
	Debug
	Production
)

type ToolsInitMode int
const (
	Single ToolsInitMode = iota
	All
)

type CmdParameter struct {
	Description string
	Type        string
	IsMandatory bool
}

type SystemCmd struct {
	Executable		string
	Description		string
	NeedsRootHandling	bool
	Parameters		[]CmdParameter
}

type SingleSubCmd struct {
	CmdGroup       string
	Summary        string
	Description    string
	IsEnabled      bool
	IsRootRequired bool
	Parameters     []CmdParameter
}

func startMCPServer() {
	AdminTasksMCPServer = server.NewMCPServer(
		"mcp_server_admintasks",
		"0.0.2",
		server.WithToolCapabilities(false),
	)
}

func ExecuteSystemCall(systemCmd SystemCmd, fullHelpText string, isRootRequired bool, subcmd string, subcmd_params ...interface{}) string {
	if subcmd == "help" {
		return(string(fullHelpText))
	} else {
		var strArgs []string
		for _, param := range systemCmd.Parameters {
			strArgs = append(strArgs, param.Description)
		}
		strArgs = append(strArgs, subcmd)
		for _, arg := range subcmd_params {
			strArgs = append(strArgs, fmt.Sprint(arg)) 
		}
		var cmd *exec.Cmd 
		if isRootRequired {
			sudoArgsA := append([]string{systemCmd.Executable}, strArgs...) 
			sudoArgsB := append([]string{"-b"}, sudoArgsA...) 
			cmd = exec.Command("sudo", sudoArgsB...)
		} else {
			cmd = exec.Command(systemCmd.Executable, strArgs...)
		}
		// Buffer to capture the output
		var out bytes.Buffer
		var resultstring string
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return ("Error running systemctl command: %v")
		} else {
			if len(out.String()) == 0 {
				resultstring = "{\"message\": \"success\"}"
			} else {
				resultstring = out.String()
			}
		}
		// return the (hopefully) JSON string ...
		return (resultstring)
	}
}

func RUN() {
	if err := server.ServeStdio(AdminTasksMCPServer); err != nil {
		fmt.Printf("Server error: %v\n")
	}
}

func INIT(mode RunningMode) {
	switch mode {
	case Production:
		utilsDebug = false
		startMCPServer()
	case Debug:
		utilsDebug = true
		startMCPServer()
	case Test:
		utilsDebug = true
	}
}
