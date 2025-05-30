package systemctl

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

var systemctlDebug bool

var allSystemCtlCmds_json string
var allSystemCtlCmds = map[string]utils.SingleCmd{
	"list-units": {
		CmdGroup:      "Unit Commands",
		Summary:       "List units currently in memory. DEFAULT action of systemctl, recommended to use with OPTION='--all'.",
		Description:   "List units currently in memory. Use OPTION='--all' to see also those units which are installed, but not enabled. This is the default of systemctl and should be called first to get an overview.",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN /regular expression", "string", false}},
	},
	"list-automounts": {
		CmdGroup:      "Unit Commands",
		Summary:       "List automount units currently in memory, ordered by path. Use PATTERN='--all' to see also those which are installed, but not enabled.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"list-paths": {
		CmdGroup:      "Unit Commands",
		Summary:       "List path units currently in memory, ordered by path. Use PATTERN='--all' to see also those which are installed, but not enabled.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"list-sockets": {
		CmdGroup:      "Unit Commands",
		Summary:       "List socket units currently in memory, ordered by address. Use PATTERN='--all' to see also those which are installed, but not enabled.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"list-timers": {
		CmdGroup:      "Unit Commands",
		Summary:       "List timer units currently in memory, ordered by next elapse. Use PATTERN='--all' to see also those which are installed, but not enabled.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"is-readonlycmd": {
		CmdGroup:      "Unit Commands",
		Summary:       "Check whether units are readonlycmd",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"is-failed": {
		CmdGroup:      "Unit Commands",
		Summary:       "Check whether units are failed or system is in degraded state",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"status": {
		CmdGroup:      "Unit Commands",
		Summary:       "Show runtime status of one or more units",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression or PID / ProcessID", "string or number", false}},
	},
	"show": {
		CmdGroup:      "Unit Commands",
		Summary:       "Show properties of one or more units/jobs or the manager",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression or jobID", "string", false}},
	},
	"cat": {
		CmdGroup:      "Unit Commands",
		Summary:       "Show files and drop-ins of specified units",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"help": {
		CmdGroup:      "Unit Commands",
		Summary:       "Show manual for one or more units. This includes extended information about the respective unit/service, which most often cannot be directly accessed by systemctl, but can be useful for either a human administrator or another MCP server to deal with.",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression or PID / ProcessID", "string", false}},
	},
	"list-dependencies": {
		CmdGroup:      "Unit Commands",
		Summary:       "Recursively show units which are required or wanted by the units or by which those units are required or wanted",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"start": {
		CmdGroup:      "Unit Commands",
		Summary:       "Start (activate) one or more units",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"stop": {
		CmdGroup:      "Unit Commands",
		Summary:       "Stop (deactivate) one or more units",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"reload": {
		CmdGroup:      "Unit Commands",
		Summary:       "Reload one or more units",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"restart": {
		CmdGroup:      "Unit Commands",
		Summary:       "Start or restart one or more units",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"try-restart": {
		CmdGroup:      "Unit Commands",
		Summary:       "Restart one or more units if readonlycmd",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"reload-or-restart": {
		CmdGroup:      "Unit Commands",
		Summary:       "Reload one or more units if possible, otherwise start or restart",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"try-reload-or-restart": {
		CmdGroup:      "Unit Commands",
		Summary:       "If readonlycmd, reload one or more units, if supported, otherwise restart",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"isolate": {
		CmdGroup:      "Unit Commands",
		Summary:       "Start one unit and stop all others",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"kill": {
		CmdGroup:      "Unit Commands",
		Summary:       "Send signal to processes of a unit",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"clean": {
		CmdGroup:      "Unit Commands",
		Summary:       "Clean runtime, cache, state, logs or configuration of unit",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}},
	},
	"freeze": {
		CmdGroup:      "Unit Commands",
		Summary:       "Freeze execution of unit processes",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"thaw": {
		CmdGroup:      "Unit Commands",
		Summary:       "Resume execution of a frozen unit",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"set-property": {
		CmdGroup:      "Unit Commands",
		Summary:       "set-property UNIT PROPERTY=VALUE... Sets one or more properties of a unit",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}, {"PROPERTY name", "string", false}, {"VALUE", "string or number", false}},
	},
	"bind": {
		CmdGroup:      "Unit Commands",
		Summary:       "Bind-mount a path from the host into a unit's namespace",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}, {"PATH", "string", false}},
	},
	"mount-image": {
		CmdGroup:      "Unit Commands",
		Summary:       "Mount an image from the host into a unit's namespace",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", false}, {"PATH", "string", false}, {"OPTIONS", "string", false}},
	},
	"service-log-level": {
		CmdGroup:      "Unit Commands",
		Summary:       "Get/set logging threshold for service.",
		Description:   "Get/set logging threshold for service. If the optional argument LEVEL is provided, then change the current log level of the service to LEVEL. The log level should be a typical syslog log level, i.e. a value in the range 0...7 or one of the strings emerg, alert, crit, err, warning, notice, info, debug; see syslog(3) for details.",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"SERVICE name", "string", false}, {"LEVEL", "string", false}},
	},
	"service-log-target": {
		CmdGroup:      "Unit Commands",
		Summary:       "Get/set logging target for service",
		Description:   "If the optional argument TARGET is provided, then change the current log target of the service to TARGET. The log target should be one of the strings console (for log output to the service's standard error stream), kmsg (for log output to the kernel log buffer), journal (for log output to systemd-journald.service(8) using the native journal protocol), syslog (for log output to the classic syslog socket /dev/log), null (for no log output whatsoever) or auto (for an automatically determined choice, typically equivalent to console if the service is invoked interactively, and journal or syslog otherwise).",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"SERVICE name", "string", false}, {"TARGET", "string", false}},
	},
	"reset-failed": {
		CmdGroup:      "Unit Commands",
		Summary:       "Reset failed state for all, one, or more units",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name or PATTERN / regular expression", "string", false}},
	},
	"whoami": {
		CmdGroup:      "Unit Commands",
		Summary:       "Return unit caller or specified PIDs are part of",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"PID", "string or number", false}},
	},
	"list-unit-files": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "list-unit-files [PATTERN...]        List installed unit files",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"PATTERN / regular expression", "string", false}},
	},
	"enable": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Enable one or more unit files",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT or PATH: what to enable", "string", true}},
	},
	"disable": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Disable one or more unit files",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT or PATH: what to disable", "string", true}},
	},
	"reenable": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Reenable one or more unit files",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"preset": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Enable/disable one or more unit files based on preset configuration",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"preset-all": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Enable/disable all unit files based on preset configuration",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"is-enabled": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Check whether unit files are enabled",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT to check whether it is enabled", "string", true}},
	},
	"mask": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Mask one or more units",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", true}},
	},
	"unmask": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Unmask one or more units",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", true}},
	},
	"link": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Link one or more units files into the search path",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"PATH", "string", true}},
	},
	"revert": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Revert one or more unit files to vendor version",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", true}},
	},
	"add-wants": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Add 'Wants' dependency for the target on specified one or more units",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"TARGET", "string", true}, {"UNIT name", "string", true}},
	},
	"add-requires": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Add 'Requires' dependency for the target on specified one or more units",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"TARGET", "string", true}, {"UNIT name", "string", true}},
	},
	"edit": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Edit one or more unit files",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"UNIT name", "string", true}},
	},
	"get-default": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Get the name of the default target",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"set-default": {
		CmdGroup:      "UnitFile Commands",
		Summary:       "Set the default target",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"TARGET name", "string", true}},
	},
	"list-machines": {
		CmdGroup:      "Machine Commands",
		Summary:       "list-machines [PATTERN...]          List local containers and host",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"PATTERN / regular expression", "string", false}},
	},
	"list-jobs": {
		CmdGroup:      "Job Commands",
		Summary:       "list-jobs [PATTERN...]              List jobs",
		Description:   "",
		IsEnabled:	true,
		IsRootRequired: false,
		Parameters:    []utils.CmdParameter{{"PATTERN / regular expression", "string", false}},
	},
	"cancel": {
		CmdGroup:      "Job Commands",
		Summary:       "cancel [JOB...]                     Cancel all, one, or more jobs",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"show-environment": {
		CmdGroup:      "Environment Commands",
		Summary:       "show-environment                    Dump environment",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"set-environment": {
		CmdGroup:      "Environment Commands",
		Summary:       "set-environment VARIABLE=VALUE...   Set one or more environment variables",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"unset-environment": {
		CmdGroup:      "Environment Commands",
		Summary:       "unset-environment VARIABLE...       Unset one or more environment variables",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"import-environment": {
		CmdGroup:      "Environment Commands",
		Summary:       "import-environment VARIABLE...      Import all or some environment variables",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"daemon-reload": {
		CmdGroup:      "ManagerState Commands",
		Summary:       "daemon-reload                       Reload systemd manager configuration",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"daemon-reexec": {
		CmdGroup:      "ManagerState Commands",
		Summary:       "daemon-reexec                       Reexecute systemd manager",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"log-level": {
		CmdGroup:      "ManagerState Commands",
		Summary:       "log-level [LEVEL]                   Get/set logging threshold for manager",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"log-target": {
		CmdGroup:      "ManagerState Commands",
		Summary:       "log-target [TARGET]                 Get/set logging target for manager",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"service-watchdogs": {
		CmdGroup:      "ManagerState Commands",
		Summary:       "service-watchdogs [BOOL]            Get/set service watchdog state",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"is-system-running": {
		CmdGroup:      "System Commands",
		Summary:       "is-system-running                   Check whether system is fully running",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"default": {
		CmdGroup:      "System Commands",
		Summary:       "default                             Enter system default mode",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"rescue": {
		CmdGroup:      "System Commands",
		Summary:       "rescue                              Enter system rescue mode",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"emergency": {
		CmdGroup:      "System Commands",
		Summary:       "emergency                           Enter system emergency mode",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"halt": {
		CmdGroup:      "System Commands",
		Summary:       "halt                                Shut down and halt the system",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"poweroff": {
		CmdGroup:      "System Commands",
		Summary:       "poweroff                            Shut down and power-off the system",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"reboot": {
		CmdGroup:      "System Commands",
		Summary:       "reboot                              Shut down and reboot the system",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"kexec": {
		CmdGroup:      "System Commands",
		Summary:       "kexec                               Shut down and reboot the system with kexec",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"soft-reboot": {
		CmdGroup:      "System Commands",
		Summary:       "soft-reboot                         Shut down and reboot userspace",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"exit": {
		CmdGroup:      "System Commands",
		Summary:       "exit [EXIT_CODE]                    Request user instance or container exit",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"switch-root": {
		CmdGroup:      "System Commands",
		Summary:       "switch-root [ROOT [INIT]]           Change to a different root file system",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"sleep": {
		CmdGroup:      "System Commands",
		Summary:       "Put the system to sleep (through one of the operations below)",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"suspend": {
		CmdGroup:      "System Commands",
		Summary:       "Suspend the system",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"hibernate": {
		CmdGroup:      "System Commands",
		Summary:       "Hibernate the system",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"hybrid-sleep": {
		CmdGroup:      "System Commands",
		Summary:       "Hibernate and suspend the system",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
	"suspend-then-hibernate": {
		CmdGroup:      "System Commands",
		Summary:       "Suspend the system, wake after a period of time, and hibernate",
		Description:   "",
		IsEnabled:	false,
		IsRootRequired: false,
	},
}

func callSystemCtl(isRootRequired bool, systemctl_cmd string, systemctl_params ...interface{}) string {

	if systemctl_cmd == "help" {

		return(string(allSystemCtlCmds_json))

	} else {
		// Construct the cmdline for systemctl: Always JSON, never shorten names, no pager
		var strArgs = []string{"-o", "json-pretty", "--full", "--no-pager", systemctl_cmd}
		for _, arg := range systemctl_params {
			strArgs = append(strArgs, fmt.Sprint(arg)) // Convert each argument to a string
		}
		var cmd *exec.Cmd 
		if isRootRequired {
		} else {
			cmd = exec.Command("systemctl", strArgs...)
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

func addSingleToolToMCPServer(cmdName string, newCmd utils.SingleCmd) {

	if newCmd.IsEnabled {
		
		newCmdName := "systemctl_" + cmdName

		var numOfParameters = 0
		if newCmd.Parameters != nil {
			numOfParameters = len(newCmd.Parameters)
		}

		sysLog, syslogerr := syslog.New(syslog.LOG_INFO, "mcp-server-systemctl")
		defer sysLog.Close()
		if syslogerr != nil {
			log.Fatalf("Failed to connect to syslog: %v", syslogerr)
		}
		if systemctlDebug {
			sysLog.Info(newCmdName)
			sysLog.Info(strconv.Itoa(numOfParameters))
		}

		switch numOfParameters {
			case 1: 
				mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary),
					mcp.WithString("systemctlp00", mcp.Required(), mcp.Description(newCmd.Parameters[0].Description),),)
				utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					systemctlp00, ok := req.GetArguments()["systemctlp00"].(string)
					if !ok { return nil, errors.New("Error in addSingleToolToMCPServer -> utils.AdminTasksMCPServer.AddTool - 1 parameter") }
					return mcp.NewToolResultText(fmt.Sprintf("%s", callSystemCtl(newCmd.IsRootRequired, cmdName, systemctlp00))), nil
				})
			case 2: 
				mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary),
					mcp.WithString("systemctlp00", mcp.Required(), mcp.Description(newCmd.Parameters[0].Description),),
					mcp.WithString("systemctlp01", mcp.Required(), mcp.Description(newCmd.Parameters[1].Description),),
				)
				utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					systemctlp00, ok := req.GetArguments()["systemctlp00"].(string)
					systemctlp01, ok := req.GetArguments()["systemctlp01"].(string)
					if !ok { return nil, errors.New("Error in addSingleToolToMCPServer -> utils.AdminTasksMCPServer.AddTool - 2 parameters") }
					return mcp.NewToolResultText(fmt.Sprintf("%s", callSystemCtl(newCmd.IsRootRequired, cmdName, systemctlp00, systemctlp01))), nil
				})
			case 3: 
				mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary),
					mcp.WithString("systemctlp00", mcp.Required(), mcp.Description(newCmd.Parameters[0].Description),),
					mcp.WithString("systemctlp01", mcp.Required(), mcp.Description(newCmd.Parameters[1].Description),),
					mcp.WithString("systemctlp02", mcp.Required(), mcp.Description(newCmd.Parameters[2].Description),),
				)
				utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					systemctlp00, ok := req.GetArguments()["systemctlp00"].(string)
					systemctlp01, ok := req.GetArguments()["systemctlp01"].(string)
					systemctlp02, ok := req.GetArguments()["systemctlp02"].(string)
					if !ok { return nil, errors.New("Error in addSingleToolToMCPServer -> utils.AdminTasksMCPServer.AddTool - 3 parameters") }
					return mcp.NewToolResultText(fmt.Sprintf("%s", callSystemCtl(newCmd.IsRootRequired, cmdName, systemctlp00, systemctlp01, systemctlp02))), nil
				})
			default: 
				mcpToolZypper := mcp.NewTool(newCmdName, mcp.WithDescription(newCmd.Summary))
				utils.AdminTasksMCPServer.AddTool(mcpToolZypper, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					return mcp.NewToolResultText(fmt.Sprintf("%s", callSystemCtl(newCmd.IsRootRequired, cmdName))), nil
				})
		}

	}

}

func addAllInOneToolsToMCPServer() {

	mcpToolSystemCtl := mcp.NewTool("tool_systemctl",
		mcp.WithDescription("Send a single cmd to systemctl and get output back in JSON. First try COMMAND 'help' to get an overview about all commands or COMMAND list-commands with OPTION --all for what is there. List units currently in memory. Use OPTION='--all' to see also those units which are installed, but not enabled. This is the default of systemctl and should be called first to get an overview."),
		mcp.WithString("systemctlcmd",
			mcp.Required(),
			mcp.Description(string(allSystemCtlCmds_json)),
		),
		mcp.WithString("systemctlparams",
			mcp.Required(),
			mcp.Description("OPTION, PATTERNS, UNITS, JOBS, ID, ... or the like for systemctl. Do not use systemctl cmds here!"),
		),
	)

	utils.AdminTasksMCPServer.AddTool(mcpToolSystemCtl, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systemctlcmd, ok := req.GetArguments()["systemctlcmd"].(string)
		systemctlparams, ok := req.GetArguments()["systemctlparams"].(string)
		// systemctlparamB, ok := req.GetArguments()["systemctlparamB"].(string)
		if !ok {
			return nil, errors.New("Error initializing tool_systemctl")
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", callSystemCtl(false, systemctlcmd, systemctlparams))), nil
	})

}

func runTests() {

}

func INIT(debugMode utils.RunningMode, initMode utils.ToolsInitMode) {
	allSystemCtlCmds_tmp, err := json.MarshalIndent(allSystemCtlCmds, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}
	allSystemCtlCmds_json = string(allSystemCtlCmds_tmp)
	switch debugMode {
	case utils.Production:
		systemctlDebug = false
		if initMode == utils.Single {
			addAllInOneToolsToMCPServer()
		} else {
			for key := range allSystemCtlCmds {
				addSingleToolToMCPServer(key, allSystemCtlCmds[key] )
			}
		}
	case utils.Debug:
		systemctlDebug = true
		if initMode == utils.Single {
			addAllInOneToolsToMCPServer()
		} else {
			for key := range allSystemCtlCmds {
				addSingleToolToMCPServer(key, allSystemCtlCmds[key] )
			}
		}
	case utils.Test:
		systemctlDebug = true
		runTests()
	}
}

