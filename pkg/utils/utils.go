package utils

import (
	"fmt"
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

type SingleCmd struct {
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
		"0.0.1",
		server.WithToolCapabilities(false),
	)
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
