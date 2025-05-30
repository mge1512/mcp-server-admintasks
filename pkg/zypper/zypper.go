package zypper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/syslog"
	"os/exec"
	"strconv"

	"mcp-server-admintasks/pkg/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

var zypperDebug bool

var allZypperCmds_json string
var allZypperCmds = map[string]utils.SingleCmd{
	"search": {
		CmdGroup:      "Querying Commands",
		Summary:       "DEFAULT action of zypper. Search for packages matching a PATTERN.",
		Description:   "PATTERN can be a regular expression. Recommended when searching for unknown packages or patterns or when the name of a package might be vague/unclear. For a more extensive search try to add the SEARCHOPTION '--search-description'.",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"PATTERN or PACKAGE name", "string", false}, {"SEARCHOPTION", "string", false}},
		// Parameters: []utils.CmdParameter{{"PATTERN or PACKAGE name", "string", false}, {"SEARCHOPTION", "string", false}, {"VALUE", "string or number", false}},
	},
	"help": {
		CmdGroup:      "General Commands",
		Summary:       "Print zypper help",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"repos": {
		CmdGroup:      "RepositoryManagement Commands",
		Summary:       "List all defined repositories.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"addrepo": {
		CmdGroup:      "RepositoryManagement Commands",
		Summary:       "Add a new repository.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"removerepo": {
		CmdGroup:      "RepositoryManagement Commands",
		Summary:       "Remove specified repository.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"renamerepo": {
		CmdGroup:      "RepositoryManagement Commands",
		Summary:       "Rename specified repository.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"modifyrepo": {
		CmdGroup:      "RepositoryManagement Commands",
		Summary:       "Modify specified repository.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"refresh": {
		CmdGroup:      "RepositoryManagement Commands",
		Summary:       "Refresh all repositories.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: true,
		Parameters:    []utils.CmdParameter{},
	},
	"clean": {
		CmdGroup:      "RepositoryManagement Commands",
		Summary:       "Clean local caches.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"services": {
		CmdGroup:      "ServiceManagement Commands",
		Summary:       "List all defined services.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"addservice": {
		CmdGroup:      "ServiceManagement Commands",
		Summary:       "Add a new service.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"modifyservice": {
		CmdGroup:      "ServiceManagement Commands",
		Summary:       "Modify specified service.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"removeservice": {
		CmdGroup:      "ServiceManagement Commands",
		Summary:       "Remove specified service.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"refresh-services": {
		CmdGroup:      "ServiceManagement Commands",
		Summary:       "Refresh all services.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"install": {
		CmdGroup:      "SoftwareManagement Commands",
		Summary:       "Install packages.",
		Description:   "If installation fails adding the INSTALLOPTION '--no-confirm' might help",
		IsEnabled:	true,
		IsRootRequired: true,
		Parameters:    []utils.CmdParameter{{"PATTERN or PACKAGE name", "string", false}, {"INSTALLOPTION", "string", false}},
	},
	"remove": {
		CmdGroup:      "SoftwareManagement Commands",
		Summary:       "Remove packages.",
		Description:   "If installation fails adding the INSTALLOPTION '--no-confirm' might help",
		IsEnabled:	true,
		IsRootRequired: true,
		Parameters:    []utils.CmdParameter{{"PATTERN or PACKAGE name", "string", false}, {"REMOVEOPTION", "string", false}},
	},
	"removeptf": {
		CmdGroup:      "SoftwareManagement Commands",
		Summary:       "Remove (not only) PTFs.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"verify": {
		CmdGroup:      "SoftwareManagement Commands",
		Summary:       "Verify integrity of package dependencies.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"source-install": {
		CmdGroup:      "SoftwareManagement Commands",
		Summary:       "Install source packages and their build dependencies.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"install-new-recommends": {
		CmdGroup:      "SoftwareManagement Commands",
		Summary:       "Install newly added packages recommended by installed packages.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"update": {
		CmdGroup:      "UpdateManagement Commands",
		Summary:       "Update installed packages with newer versions.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: true,
		Parameters:    []utils.CmdParameter{},
	},
	"list-updates": {
		CmdGroup:      "UpdateManagement Commands",
		Summary:       "List available updates.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"patch": {
		CmdGroup:      "UpdateManagement Commands",
		Summary:       "Install needed patches.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"list-patches": {
		CmdGroup:      "UpdateManagement Commands",
		Summary:       "List available patches.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"dist-upgrade": {
		CmdGroup:      "UpdateManagement Commands",
		Summary:       "Perform a distribution upgrade.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"patch-check": {
		CmdGroup:      "UpdateManagement Commands",
		Summary:       "Check for patches.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"info": {
		CmdGroup:      "Querying Commands",
		Summary:       "Show full information for specified packages.",
		Description:   "Ask for full/detailed information about a single package with a known name (version, size, status, description, installation status).--  Do not use, if the name is not fully clear (use zypper search for that). ",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"PACKAGE name", "string", false}},
	},
	"patch-info": {
		CmdGroup:      "Querying Commands",
		Summary:       "Show full information for specified patches.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"pattern-info": {
		CmdGroup:      "Querying Commands",
		Summary:       "Show full information for specified patterns.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"product-info": {
		CmdGroup:      "Querying Commands",
		Summary:       "Show full information for specified products.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"patches": {
		CmdGroup:      "Querying Commands",
		Summary:       "List all available patches.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"packages": {
		CmdGroup:      "Querying Commands",
		Summary:       "List all available packages.",
		Description:   "If the package name is known, better use zypper search ",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"patterns": {
		CmdGroup:      "Querying Commands",
		Summary:       "List all available patterns.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"products": {
		CmdGroup:      "Querying Commands",
		Summary:       "List all available products.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"what-provides": {
		CmdGroup:      "Querying Commands",
		Summary:       "List packages providing specified capability.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"addlock": {
		CmdGroup:      "PackageLocks Commands",
		Summary:       "Add a package lock.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"removelock": {
		CmdGroup:      "PackageLocks Commands",
		Summary:       "Remove a package lock.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"locks": {
		CmdGroup:      "PackageLocks Commands",
		Summary:       "List current package locks.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"cleanlocks": {
		CmdGroup:      "PackageLocks Commands",
		Summary:       "Remove useless locks.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"locales": {
		CmdGroup:      "LocaleManagement Commands",
		Summary:       "List requested locales (languages codes).",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"addlocale": {
		CmdGroup:      "LocaleManagement Commands",
		Summary:       "Add locale(s) to requested locales.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"removelocale": {
		CmdGroup:      "LocaleManagement Commands",
		Summary:       "Remove locale(s) from requested locales.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"versioncmp": {
		CmdGroup:      "Other Commands",
		Summary:       "Compare two version strings.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"targetos": {
		CmdGroup:      "Other Commands",
		Summary:       "Print the target operating system ID string.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"licenses": {
		CmdGroup:      "Other Commands",
		Summary:       "Print report about licenses and EULAs of installed packages.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"download": {
		CmdGroup:      "Other Commands",
		Summary:       "Download rpms specified on the commandline to a local directory.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"source-download": {
		CmdGroup:      "Other Commands",
		Summary:       "Download source rpms for all installed packages to a local directory.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"needs-rebooting": {
		CmdGroup:      "Other Commands",
		Summary:       "Check if the reboot-needed flag was set.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"ps": {
		CmdGroup:      "Other Commands",
		Summary:       "List running processes which might still use files and libraries deleted",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"purge-kernels": {
		CmdGroup:      "Other Commands",
		Summary:       "Remove old kernels.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"system-architecture": {
		CmdGroup:      "Other Commands",
		Summary:       "Print the detected system architecture.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
	"subcommand": {
		CmdGroup:      "Subcommands Commands",
		Summary:       "Lists available subcommands.",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{},
	},
}

func callZypper(isRootRequired bool, zypper_cmd string, zypper_params ...interface{}) string {
	
	if zypper_cmd == "help" {

		return(string(allZypperCmds_json))

	} else {

		// Construct the cmdline for zypper: Always use XML, prep for machine operations
		var strArgs = []string{"--xmlout", "--terse", "--non-interactive", zypper_cmd}
		for _, arg := range zypper_params {
			strArgs = append(strArgs, fmt.Sprint(arg)) // Convert each argument to a string
		}
		var cmd *exec.Cmd 
		if isRootRequired {
			sudoArgsA := append([]string{"/usr/bin/zypper"}, strArgs...) 
			sudoArgsB := append([]string{"-b"}, sudoArgsA...) 
			cmd = exec.Command("sudo", sudoArgsB...)
		} else {
			cmd = exec.Command("/usr/bin/zypper", strArgs...)
		}
		// Using syslog is useful for debugging
		sysLog, syslogerr := syslog.New(syslog.LOG_INFO, "mcp-server-zypper")
		if zypperDebug {
			if syslogerr != nil {
				log.Fatalf("Failed to connect to syslog: %v", syslogerr)
			}
		}
		defer sysLog.Close()
		// Buffer to capture the output
		var out bytes.Buffer
		var resultstring string = out.String()
		cmd.Stdout = &out
		err := cmd.Run()
		if zypperDebug {
			sysLog.Info(cmd.String())
		}
		if err != nil {
			return ("Error running zypper command: %v")
		} else {
			if len(out.String()) == 0 {
				resultstring = "<message>success</message>"
			} else {
				resultstring = out.String()
			}
		}
		if zypperDebug {
			sysLog.Info(resultstring)
		}
		// return XML
		return resultstring
	}
}

func addSingleToolToMCPServer(cmdName string, newCmd utils.SingleCmd) {

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
					mcp.WithString("zypperp00", mcp.Required(), mcp.Description(newCmd.Parameters[0].Description),),)
				utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					zypperp00, ok := req.GetArguments()["zypperp00"].(string)
					if !ok { return nil, errors.New("Error in addSingleToolToMCPServer -> utils.AdminTasksMCPServer.AddTool - 1 parameter") }
					return mcp.NewToolResultText(fmt.Sprintf("%s", callZypper(newCmd.IsRootRequired, cmdName, zypperp00))), nil
				})
			case 2: 
				mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary),
					mcp.WithString("zypperp00", mcp.Required(), mcp.Description(newCmd.Parameters[0].Description),),
					mcp.WithString("zypperp01", mcp.Required(), mcp.Description(newCmd.Parameters[1].Description),),
				)
				utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					zypperp00, ok := req.GetArguments()["zypperp00"].(string)
					zypperp01, ok := req.GetArguments()["zypperp01"].(string)
					if !ok { return nil, errors.New("Error in addSingleToolToMCPServer -> utils.AdminTasksMCPServer.AddTool - 2 parameters") }
					return mcp.NewToolResultText(fmt.Sprintf("%s", callZypper(newCmd.IsRootRequired, cmdName, zypperp00, zypperp01))), nil
				})
			case 3: 
				mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary),
					mcp.WithString("zypperp00", mcp.Required(), mcp.Description(newCmd.Parameters[0].Description),),
					mcp.WithString("zypperp01", mcp.Required(), mcp.Description(newCmd.Parameters[1].Description),),
					mcp.WithString("zypperp02", mcp.Required(), mcp.Description(newCmd.Parameters[2].Description),),
				)
				utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					zypperp00, ok := req.GetArguments()["zypperp00"].(string)
					zypperp01, ok := req.GetArguments()["zypperp01"].(string)
					zypperp02, ok := req.GetArguments()["zypperp02"].(string)
					if !ok { return nil, errors.New("Error in addSingleToolToMCPServer -> utils.AdminTasksMCPServer.AddTool - 3 parameters") }
					return mcp.NewToolResultText(fmt.Sprintf("%s", callZypper(newCmd.IsRootRequired, cmdName, zypperp00, zypperp01, zypperp02))), nil
				})
			default: 
				mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary))
				utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					return mcp.NewToolResultText(fmt.Sprintf("%s", callZypper(newCmd.IsRootRequired, cmdName))), nil
				})
		}
	}
}

func addToolsToMCPServer() {

	mcpToolZypper := mcp.NewTool("tool_zypper",
		mcp.WithDescription("Send a single cmd to zypper and get output back in XML (or JSON)"),
		mcp.WithString("zyppercmd",
			mcp.Required(),
			mcp.Description(string(allZypperCmds_json)),
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

		return mcp.NewToolResultText(fmt.Sprintf("%s", callZypper(false, zyppercmd, zypperp01, zypperp02))), nil
	})

}

func runTests() {
	addToolsToMCPServer()
}

func INIT(debugMode utils.RunningMode, initMode utils.ToolsInitMode) {
	allZypperCmds_tmp, err := json.MarshalIndent(allZypperCmds, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}
	allZypperCmds_json = string(allZypperCmds_tmp)
	switch debugMode {
	case utils.Production:
		zypperDebug = false
		if initMode == utils.Single {
			addToolsToMCPServer()
		} else {
			for key := range allZypperCmds {
				addSingleToolToMCPServer(key, allZypperCmds[key] )
			}
		}
	case utils.Debug:
		zypperDebug = true
		if initMode == utils.Single {
			addToolsToMCPServer()
		} else {
			for key := range allZypperCmds {
				addSingleToolToMCPServer(key, allZypperCmds[key] )
			}
		}
	case utils.Test:
		zypperDebug = true
		runTests()
	}
}

