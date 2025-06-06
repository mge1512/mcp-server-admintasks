package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
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
	Typed
)

type SingleSubCmd struct {
	CmdGroup       string   `json:"cmd_group"`
	Summary        string   `json:"summary"`
	Description    string   `json:"description"`
	IsEnabled      bool     `json:"is_enabled"`
	IsRootRequired bool     `json:"is_root_required"`
	Parameters     []string `json:"parameters"`
}

type SystemCmd struct {
	Executable        string                  `json:"executable"`
	Description       string                  `json:"description"`
	NeedsRootHandling bool                    `json:"needs_root_handling"`
	DefaultParameters []string                `json:"default_parameters"`
	SubCommands       map[string]SingleSubCmd `json:"subcommands"`
}

func readSystemCmdJSONIntoStruct(directoryPath string) (map[string]SystemCmd, error) {
	// Get list of files
	allSystemCmds := make(map[string]SystemCmd)
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}
	// Loop through each file
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue // Skip directories and non-JSON files
		}
		filePath := filepath.Join(directoryPath, file.Name())

		// Open JSON file
		file, err := os.Open(filePath)
		if err != nil {
			// fmt.Println("Error opening file:", filePath, err)
			continue
		}
		defer file.Close()
		// Decode JSON into struct
		var newSystemCmd SystemCmd
		err = json.NewDecoder(file).Decode(&newSystemCmd)
		if err != nil {
			// fmt.Println("Error decoding JSON:", err)
			continue
		}
		// Store the parsed struct in the map using Executable as key
		allSystemCmds[newSystemCmd.Executable] = newSystemCmd
	}
	return allSystemCmds, nil
}

func startMCPServer() {
	AdminTasksMCPServer = server.NewMCPServer(
		"mcp_server_admintasks",
		"0.0.2",
		server.WithToolCapabilities(false),
	)
}

func ExecuteSystemCall(systemCmd SystemCmd, fullHelpText string, isRootRequired bool, subcmd string, subcmd_params ...string) string {
	sysLog, syslogerr := syslog.New(syslog.LOG_INFO, "ExecuteSystemCall")
	defer sysLog.Close()
	if syslogerr != nil {
		return ("Failed to connect to syslog: %v")
	}
	if subcmd == "help" {
		if utilsDebug {
			sysLog.Info(string(fullHelpText))
		}
		return (string(fullHelpText))
	} else {
		if utilsDebug {
			sysLog.Info(systemCmd.Executable)
			sysLog.Info(subcmd)
		}
		// This will be the cmdline parameters after the main executable
		var strArgs []string
		// First add the default options
		for _, param := range systemCmd.DefaultParameters {
			strArgs = append(strArgs, param)
		}
		// Second add the subcommand
		strArgs = append(strArgs, subcmd)
		// Third, add the subcmd parameters
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
		if utilsDebug {
			sysLog.Info(cmd.String())
		}
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

func AddToolToMCPServer(systemCmd SystemCmd, fullHelpText string, cmdName string, newCmd SingleSubCmd) {

	if newCmd.IsEnabled {

		newCmdName := systemCmd.Executable + "_" + cmdName

		sysLog, syslogerr := syslog.New(syslog.LOG_INFO, newCmdName)
		defer sysLog.Close()

		if syslogerr != nil {
			log.Fatalf("Failed to connect to syslog: %v", syslogerr)
		}
		if utilsDebug {
			sysLog.Info(newCmdName)
		}

		var summaryWithParameters string = newCmd.Summary + ". Parameters: "
		for _, arg := range newCmd.Parameters {
			summaryWithParameters = summaryWithParameters + fmt.Sprint(arg)
		}

		mcpToolZypper := mcp.NewTool(newCmdName,
			mcp.WithDescription(newCmd.Summary),
			mcp.WithArray("Parameters", mcp.Items(map[string]any{"type": "string"})),
		)

		AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			systemctlargs := req.GetArguments()["Parameters"]
			var strList []string
			argsSlice, ok := systemctlargs.([]interface{})
			if ok {
				for _, arg := range argsSlice {
					strList = append(strList, fmt.Sprint(arg.(string)))
				}
			}
			return mcp.NewToolResultText(fmt.Sprintf("%s", ExecuteSystemCall(systemCmd, fullHelpText, newCmd.IsRootRequired, cmdName, strList...))), nil
		})

	}

}

func RUN() {
	// readSystemCmdJSONIntoStruct("/usr/share/mcp-server-admintasks/")
	readSystemCmdJSONIntoStruct("/tmp/json/")
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
