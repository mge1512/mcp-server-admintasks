package zypper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/syslog"
	"strconv"

	"mcp-server-admintasks/pkg/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

var zypperDebug bool

var jsonZypperSubCmds string

var zypperCmd utils.SystemCmd = utils.SystemCmd{
	Executable:        "zypper",
	Description:       "Command-line interface to ZYpp system management library (libzypp)",
	NeedsRootHandling: true,
	DefaultParameters: []string{"--xmlout", "--terse", "--non-interactive"},
	SubCommands: map[string]utils.SingleSubCmd{
		"search": {
			CmdGroup:       "Querying Commands",
			Summary:        "DEFAULT action of zypper. Search for packages matching a PATTERN.",
			Description:    "PATTERN can be a regular expression. Recommended when searching for unknown packages or patterns or when the name of a package might be vague/unclear. For a more extensive search try to add the SEARCHOPTION '--search-description'.",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"PATTERN or PACKAGE name", "SEARCHOPTION"},
		},
		"help": {
			CmdGroup:       "General Commands",
			Summary:        "Print zypper help",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"repos": {
			CmdGroup:       "RepositoryManagement Commands",
			Summary:        "List all defined repositories.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"addrepo": {
			CmdGroup:       "RepositoryManagement Commands",
			Summary:        "Add a new repository.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"removerepo": {
			CmdGroup:       "RepositoryManagement Commands",
			Summary:        "Remove specified repository.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"renamerepo": {
			CmdGroup:       "RepositoryManagement Commands",
			Summary:        "Rename specified repository.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"modifyrepo": {
			CmdGroup:       "RepositoryManagement Commands",
			Summary:        "Modify specified repository.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"refresh": {
			CmdGroup:       "RepositoryManagement Commands",
			Summary:        "Refresh all repositories.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: true,
			Parameters:     []string{},
		},
		"clean": {
			CmdGroup:       "RepositoryManagement Commands",
			Summary:        "Clean local caches.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"services": {
			CmdGroup:       "ServiceManagement Commands",
			Summary:        "List all defined services.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"addservice": {
			CmdGroup:       "ServiceManagement Commands",
			Summary:        "Add a new service.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"modifyservice": {
			CmdGroup:       "ServiceManagement Commands",
			Summary:        "Modify specified service.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"removeservice": {
			CmdGroup:       "ServiceManagement Commands",
			Summary:        "Remove specified service.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"refresh-services": {
			CmdGroup:       "ServiceManagement Commands",
			Summary:        "Refresh all services.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"install": {
			CmdGroup:       "SoftwareManagement Commands",
			Summary:        "Install packages.",
			Description:    "If installation fails adding the INSTALLOPTION '--no-confirm' might help",
			IsEnabled:      true,
			IsRootRequired: true,
			Parameters:     []string{"PATTERN or PACKAGE name", "INSTALLOPTION"},
		},
		"remove": {
			CmdGroup:       "SoftwareManagement Commands",
			Summary:        "Remove packages.",
			Description:    "If installation fails adding the INSTALLOPTION '--no-confirm' might help",
			IsEnabled:      true,
			IsRootRequired: true,
			Parameters:     []string{"PATTERN or PACKAGE name", "REMOVEOPTION"},
		},
		"removeptf": {
			CmdGroup:       "SoftwareManagement Commands",
			Summary:        "Remove (not only) PTFs.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"verify": {
			CmdGroup:       "SoftwareManagement Commands",
			Summary:        "Verify integrity of package dependencies.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"source-install": {
			CmdGroup:       "SoftwareManagement Commands",
			Summary:        "Install source packages and their build dependencies.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"install-new-recommends": {
			CmdGroup:       "SoftwareManagement Commands",
			Summary:        "Install newly added packages recommended by installed packages.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"update": {
			CmdGroup:       "UpdateManagement Commands",
			Summary:        "Update installed packages with newer versions.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: true,
			Parameters:     []string{},
		},
		"list-updates": {
			CmdGroup:       "UpdateManagement Commands",
			Summary:        "List available updates.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"patch": {
			CmdGroup:       "UpdateManagement Commands",
			Summary:        "Install needed patches.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"list-patches": {
			CmdGroup:       "UpdateManagement Commands",
			Summary:        "List available patches.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"dist-upgrade": {
			CmdGroup:       "UpdateManagement Commands",
			Summary:        "Perform a distribution upgrade.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"patch-check": {
			CmdGroup:       "UpdateManagement Commands",
			Summary:        "Check for patches.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"info": {
			CmdGroup:       "Querying Commands",
			Summary:        "Show full information for specified packages.",
			Description:    "Ask for full/detailed information about a single package with a known name (version, size, status, description, installation status).--  Do not use, if the name is not fully clear (use zypper search for that). ",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"PACKAGE name"},
		},
		"patch-info": {
			CmdGroup:       "Querying Commands",
			Summary:        "Show full information for specified patches.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"pattern-info": {
			CmdGroup:       "Querying Commands",
			Summary:        "Show full information for specified patterns.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"product-info": {
			CmdGroup:       "Querying Commands",
			Summary:        "Show full information for specified products.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"patches": {
			CmdGroup:       "Querying Commands",
			Summary:        "List all available patches.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"packages": {
			CmdGroup:       "Querying Commands",
			Summary:        "List all available packages.",
			Description:    "If the package name is known, better use zypper search ",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"patterns": {
			CmdGroup:       "Querying Commands",
			Summary:        "List all available patterns.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"products": {
			CmdGroup:       "Querying Commands",
			Summary:        "List all available products.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"what-provides": {
			CmdGroup:       "Querying Commands",
			Summary:        "List packages providing specified capability.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"addlock": {
			CmdGroup:       "PackageLocks Commands",
			Summary:        "Add a package lock.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"removelock": {
			CmdGroup:       "PackageLocks Commands",
			Summary:        "Remove a package lock.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"locks": {
			CmdGroup:       "PackageLocks Commands",
			Summary:        "List current package locks.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"cleanlocks": {
			CmdGroup:       "PackageLocks Commands",
			Summary:        "Remove useless locks.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"locales": {
			CmdGroup:       "LocaleManagement Commands",
			Summary:        "List requested locales (languages codes).",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"addlocale": {
			CmdGroup:       "LocaleManagement Commands",
			Summary:        "Add locale(s) to requested locales.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"removelocale": {
			CmdGroup:       "LocaleManagement Commands",
			Summary:        "Remove locale(s) from requested locales.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"versioncmp": {
			CmdGroup:       "Other Commands",
			Summary:        "Compare two version strings.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"targetos": {
			CmdGroup:       "Other Commands",
			Summary:        "Print the target operating system ID string.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"licenses": {
			CmdGroup:       "Other Commands",
			Summary:        "Print report about licenses and EULAs of installed packages.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"download": {
			CmdGroup:       "Other Commands",
			Summary:        "Download rpms specified on the commandline to a local directory.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"source-download": {
			CmdGroup:       "Other Commands",
			Summary:        "Download source rpms for all installed packages to a local directory.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"needs-rebooting": {
			CmdGroup:       "Other Commands",
			Summary:        "Check if the reboot-needed flag was set.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"ps": {
			CmdGroup:       "Other Commands",
			Summary:        "List running processes which might still use files and libraries deleted",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"purge-kernels": {
			CmdGroup:       "Other Commands",
			Summary:        "Remove old kernels.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"system-architecture": {
			CmdGroup:       "Other Commands",
			Summary:        "Print the detected system architecture.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
		"subcommand": {
			CmdGroup:       "Subcommands Commands",
			Summary:        "Lists available subcommands.",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{},
		},
	},
}

func addSingleToolToMCPServer(cmdName string, newCmd utils.SingleSubCmd) {

	if newCmd.IsEnabled {

		newCmdName := "zypper_" + cmdName

		var numOfParameters = 0
		if newCmd.Parameters != nil {
			numOfParameters = len(newCmd.Parameters)
		}

		sysLog, syslogerr := syslog.New(syslog.LOG_INFO, "mcp-server-zypper")
		defer sysLog.Close()
		if syslogerr != nil {
			log.Fatalf("Failed to connect to syslog: %v", syslogerr)
		}
		if zypperDebug {
			sysLog.Info(newCmdName)
			sysLog.Info(strconv.Itoa(numOfParameters))
		}

		switch numOfParameters {
		case 1:
			mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary),
				mcp.WithString("zypperp00", mcp.Required(), mcp.Description(newCmd.Parameters[0])))
			utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				zypperp00, ok := req.GetArguments()["zypperp00"].(string)
				if !ok {
					return nil, errors.New("Error in addSingleToolToMCPServer -> utils.AdminTasksMCPServer.AddTool - 1 parameter")
				}
				return mcp.NewToolResultText(fmt.Sprintf("%s", utils.ExecuteSystemCall(zypperCmd, jsonZypperSubCmds, newCmd.IsRootRequired, cmdName, zypperp00))), nil
			})
		case 2:
			mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary),
				mcp.WithString("zypperp00", mcp.Required(), mcp.Description(newCmd.Parameters[0])),
				mcp.WithString("zypperp01", mcp.Required(), mcp.Description(newCmd.Parameters[1])),
			)
			utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				zypperp00, ok := req.GetArguments()["zypperp00"].(string)
				zypperp01, ok := req.GetArguments()["zypperp01"].(string)
				if !ok {
					return nil, errors.New("Error in addSingleToolToMCPServer -> utils.AdminTasksMCPServer.AddTool - 2 parameters")
				}
				return mcp.NewToolResultText(fmt.Sprintf("%s", utils.ExecuteSystemCall(zypperCmd, jsonZypperSubCmds, newCmd.IsRootRequired, cmdName, zypperp00, zypperp01))), nil
			})
		case 3:
			mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary),
				mcp.WithString("zypperp00", mcp.Required(), mcp.Description(newCmd.Parameters[0])),
				mcp.WithString("zypperp01", mcp.Required(), mcp.Description(newCmd.Parameters[1])),
				mcp.WithString("zypperp02", mcp.Required(), mcp.Description(newCmd.Parameters[2])),
			)
			utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				zypperp00, ok := req.GetArguments()["zypperp00"].(string)
				zypperp01, ok := req.GetArguments()["zypperp01"].(string)
				zypperp02, ok := req.GetArguments()["zypperp02"].(string)
				if !ok {
					return nil, errors.New("Error in addSingleToolToMCPServer -> utils.AdminTasksMCPServer.AddTool - 3 parameters")
				}
				return mcp.NewToolResultText(fmt.Sprintf("%s", utils.ExecuteSystemCall(zypperCmd, jsonZypperSubCmds, newCmd.IsRootRequired, cmdName, zypperp00, zypperp01, zypperp02))), nil
			})
		default:
			mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary))
			utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcp.NewToolResultText(fmt.Sprintf("%s", utils.ExecuteSystemCall(zypperCmd, jsonZypperSubCmds, newCmd.IsRootRequired, cmdName))), nil
			})
		}
	}
}

func addToolsToMCPServer() {

	mcpToolZypper := mcp.NewTool("tool_zypper",
		mcp.WithDescription("Send a single cmd to zypper and get output back in XML (or JSON)"),
		mcp.WithString("zyppercmd",
			mcp.Required(),
			mcp.Description(string(jsonZypperSubCmds)),
		),
		mcp.WithString("zypperp01",
			mcp.Required(),
			mcp.Description("PACKAGES, PATTERNS, ... or the like for zypper. Do not use zyppercmd here!"),
		),
		mcp.WithString("zypperp02",
			mcp.Required(),
			mcp.Description("PACKAGES, PATTERNS, ... or the like for zypper. Do not use zyppercmd here!"),
		),
	)

	utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		zyppercmd, ok := req.GetArguments()["zyppercmd"].(string)
		zypperp01, ok := req.GetArguments()["zypperp01"].(string)
		zypperp02, ok := req.GetArguments()["zypperp02"].(string)

		if !ok {
			return nil, errors.New("Error in addToolsToMCPServer -> utils.AdminTasksMCPServer.AddTool")
		}

		return mcp.NewToolResultText(fmt.Sprintf("%s", utils.ExecuteSystemCall(zypperCmd, jsonZypperSubCmds, false, zyppercmd, zypperp01, zypperp02))), nil
	})

}

func runTests() {
	addToolsToMCPServer()
}

func INIT(debugMode utils.RunningMode, initMode utils.ToolsInitMode) {
	tmpZypperSubCmds, err := json.MarshalIndent(zypperCmd.SubCommands, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}
	jsonZypperSubCmds = string(tmpZypperSubCmds)
	switch debugMode {
	case utils.Production:
		zypperDebug = false
		if initMode == utils.Single {
			addToolsToMCPServer()
		} else if initMode == utils.Typed {
			for key := range zypperCmd.SubCommands {
				utils.AddToolToMCPServer(zypperCmd, jsonZypperSubCmds, key, zypperCmd.SubCommands[key])
			}
		}
	case utils.Debug:
		zypperDebug = true
		if initMode == utils.Single {
			addToolsToMCPServer()
		} else if initMode == utils.Typed {
			for key := range zypperCmd.SubCommands {
				utils.AddToolToMCPServer(zypperCmd, jsonZypperSubCmds, key, zypperCmd.SubCommands[key])
			}
		}
	case utils.Test:
		zypperDebug = true
		runTests()
	}
}
