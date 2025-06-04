package systemctl

import (
	"encoding/json"
	"fmt"
	"mcp-server-admintasks/pkg/utils"
)

var systemctlDebug bool

var jsonSystemCtlSubCmds string

var systemCtlCmd utils.SystemCmd = utils.SystemCmd{
	Executable:        "systemctl",
	Description:       "Query or send control commands to the system manager",
	NeedsRootHandling: false,
	DefaultParameters: []string{"--output=json-pretty", "--full", "--no-pager"},
	SubCommands: map[string]utils.SingleSubCmd{
		"list-units": {
			CmdGroup:       "Unit Commands",
			Summary:        "List units currently in memory. DEFAULT action of systemctl, recommended to use with OPTION='--all'.",
			Description:    "List units currently in memory. DEFAULT action of systemctl. Use OPTION='--all' to see also those units which are installed, but not enabled. This is the default of systemctl and should be called first to get an overview.",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"list-automounts": {
			CmdGroup:       "Unit Commands",
			Summary:        "List automount units currently in memory, ordered by path. Use PATTERN='--all' to see also those which are installed, but not enabled.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"list-paths": {
			CmdGroup:       "Unit Commands",
			Summary:        "List path units currently in memory, ordered by path. Use PATTERN='--all' to see also those which are installed, but not enabled.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"list-sockets": {
			CmdGroup:       "Unit Commands",
			Summary:        "List socket units currently in memory, ordered by address. Use PATTERN='--all' to see also those which are installed, but not enabled.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"list-timers": {
			CmdGroup:       "Unit Commands",
			Summary:        "List timer units currently in memory, ordered by next elapse. Use PATTERN='--all' to see also those which are installed, but not enabled.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"is-readonlycmd": {
			CmdGroup:       "Unit Commands",
			Summary:        "Check whether units are readonlycmd",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"is-failed": {
			CmdGroup:       "Unit Commands",
			Summary:        "Check whether units are failed or system is in degraded state",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"status": {
			CmdGroup:       "Unit Commands",
			Summary:        "Show runtime status of one or more units",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression or PID / ProcessID"},
		},
		"show": {
			CmdGroup:       "Unit Commands",
			Summary:        "Show properties of one or more units/jobs or the manager",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression or jobID"},
		},
		"cat": {
			CmdGroup:       "Unit Commands",
			Summary:        "Show files and drop-ins of specified units",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"help": {
			CmdGroup:       "Unit Commands",
			Summary:        "Show manual for one or more units. This includes extended information about the respective unit/service, which most often cannot be directly accessed by systemctl, but can be useful for either a human administrator or another MCP server to deal with.",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression or PID / ProcessID"},
		},
		"list-dependencies": {
			CmdGroup:       "Unit Commands",
			Summary:        "Recursively show units which are required or wanted by the units or by which those units are required or wanted",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"start": {
			CmdGroup:       "Unit Commands",
			Summary:        "Start (activate) one or more units",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"stop": {
			CmdGroup:       "Unit Commands",
			Summary:        "Stop (deactivate) one or more units",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"reload": {
			CmdGroup:       "Unit Commands",
			Summary:        "Reload one or more units",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"restart": {
			CmdGroup:       "Unit Commands",
			Summary:        "Start or restart one or more units",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"try-restart": {
			CmdGroup:       "Unit Commands",
			Summary:        "Restart one or more units if readonlycmd",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"reload-or-restart": {
			CmdGroup:       "Unit Commands",
			Summary:        "Reload one or more units if possible, otherwise start or restart",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"try-reload-or-restart": {
			CmdGroup:       "Unit Commands",
			Summary:        "If readonlycmd, reload one or more units, if supported, otherwise restart",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"isolate": {
			CmdGroup:       "Unit Commands",
			Summary:        "Start one unit and stop all others",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"kill": {
			CmdGroup:       "Unit Commands",
			Summary:        "Send signal to processes of a unit",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"clean": {
			CmdGroup:       "Unit Commands",
			Summary:        "Clean runtime, cache, state, logs or configuration of unit",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"freeze": {
			CmdGroup:       "Unit Commands",
			Summary:        "Freeze execution of unit processes",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"thaw": {
			CmdGroup:       "Unit Commands",
			Summary:        "Resume execution of a frozen unit",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"set-property": {
			CmdGroup:       "Unit Commands",
			Summary:        "set-property UNIT PROPERTY=VALUE... Sets one or more properties of a unit",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name", "PROPERTY name", "VALUE"},
		},
		"bind": {
			CmdGroup:       "Unit Commands",
			Summary:        "Bind-mount a path from the host into a unit's namespace",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name", "PATH"},
		},
		"mount-image": {
			CmdGroup:       "Unit Commands",
			Summary:        "Mount an image from the host into a unit's namespace",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name", "PATH", "OPTIONS"},
		},
		"service-log-level": {
			CmdGroup:       "Unit Commands",
			Summary:        "Get/set logging threshold for service.",
			Description:    "Get/set logging threshold for service. If the optional argument LEVEL is provided, then change the current log level of the service to LEVEL. The log level should be a typical syslog log level, i.e. a value in the range 0...7 or one of the strings emerg, alert, crit, err, warning, notice, info, debug; see syslog(3) for details.",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"SERVICE name", "LEVEL"},
		},
		"service-log-target": {
			CmdGroup:       "Unit Commands",
			Summary:        "Get/set logging target for service",
			Description:    "If the optional argument TARGET is provided, then change the current log target of the service to TARGET. The log target should be one of the strings console (for log output to the service's standard error stream), kmsg (for log output to the kernel log buffer), journal (for log output to systemd-journald.service(8) using the native journal protocol), syslog (for log output to the classic syslog socket /dev/log), null (for no log output whatsoever) or auto (for an automatically determined choice, typically equivalent to console if the service is invoked interactively, and journal or syslog otherwise).",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"SERVICE name", "TARGET"},
		},
		"reset-failed": {
			CmdGroup:       "Unit Commands",
			Summary:        "Reset failed state for all, one, or more units",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name or PATTERN / regular expression"},
		},
		"whoami": {
			CmdGroup:       "Unit Commands",
			Summary:        "Return unit caller or specified PIDs are part of",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"PID"},
		},
		"list-unit-files": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "list-unit-files [PATTERN...]        List installed unit files",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"PATTERN / regular expression"},
		},
		"enable": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Enable one or more unit files",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT or PATH: what to enable"},
		},
		"disable": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Disable one or more unit files",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"UNIT or PATH: what to disable"},
		},
		"reenable": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Reenable one or more unit files",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"preset": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Enable/disable one or more unit files based on preset configuration",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"preset-all": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Enable/disable all unit files based on preset configuration",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"is-enabled": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Check whether unit files are enabled",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT to check whether it is enabled"},
		},
		"mask": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Mask one or more units",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"unmask": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Unmask one or more units",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"link": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Link one or more units files into the search path",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"PATH"},
		},
		"revert": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Revert one or more unit files to vendor version",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"add-wants": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Add 'Wants' dependency for the target on specified one or more units",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"TARGET", "UNIT name"},
		},
		"add-requires": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Add 'Requires' dependency for the target on specified one or more units",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"TARGET", "UNIT name"},
		},
		"edit": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Edit one or more unit files",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"UNIT name"},
		},
		"get-default": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Get the name of the default target",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"set-default": {
			CmdGroup:       "UnitFile Commands",
			Summary:        "Set the default target",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
			Parameters:     []string{"TARGET name"},
		},
		"list-machines": {
			CmdGroup:       "Machine Commands",
			Summary:        "list-machines [PATTERN...]          List local containers and host",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"PATTERN / regular expression"},
		},
		"list-jobs": {
			CmdGroup:       "Job Commands",
			Summary:        "list-jobs [PATTERN...]              List jobs",
			Description:    "",
			IsEnabled:      true,
			IsRootRequired: false,
			Parameters:     []string{"PATTERN / regular expression"},
		},
		"cancel": {
			CmdGroup:       "Job Commands",
			Summary:        "cancel [JOB...]                     Cancel all, one, or more jobs",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"show-environment": {
			CmdGroup:       "Environment Commands",
			Summary:        "show-environment                    Dump environment",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"set-environment": {
			CmdGroup:       "Environment Commands",
			Summary:        "set-environment VARIABLE=VALUE...   Set one or more environment variables",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"unset-environment": {
			CmdGroup:       "Environment Commands",
			Summary:        "unset-environment VARIABLE...       Unset one or more environment variables",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"import-environment": {
			CmdGroup:       "Environment Commands",
			Summary:        "import-environment VARIABLE...      Import all or some environment variables",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"daemon-reload": {
			CmdGroup:       "ManagerState Commands",
			Summary:        "daemon-reload                       Reload systemd manager configuration",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"daemon-reexec": {
			CmdGroup:       "ManagerState Commands",
			Summary:        "daemon-reexec                       Reexecute systemd manager",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"log-level": {
			CmdGroup:       "ManagerState Commands",
			Summary:        "log-level [LEVEL]                   Get/set logging threshold for manager",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"log-target": {
			CmdGroup:       "ManagerState Commands",
			Summary:        "log-target [TARGET]                 Get/set logging target for manager",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"service-watchdogs": {
			CmdGroup:       "ManagerState Commands",
			Summary:        "service-watchdogs [BOOL]            Get/set service watchdog state",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"is-system-running": {
			CmdGroup:       "System Commands",
			Summary:        "is-system-running                   Check whether system is fully running",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"default": {
			CmdGroup:       "System Commands",
			Summary:        "default                             Enter system default mode",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"rescue": {
			CmdGroup:       "System Commands",
			Summary:        "rescue                              Enter system rescue mode",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"emergency": {
			CmdGroup:       "System Commands",
			Summary:        "emergency                           Enter system emergency mode",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"halt": {
			CmdGroup:       "System Commands",
			Summary:        "halt                                Shut down and halt the system",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"poweroff": {
			CmdGroup:       "System Commands",
			Summary:        "poweroff                            Shut down and power-off the system",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"reboot": {
			CmdGroup:       "System Commands",
			Summary:        "reboot                              Shut down and reboot the system",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"kexec": {
			CmdGroup:       "System Commands",
			Summary:        "kexec                               Shut down and reboot the system with kexec",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"soft-reboot": {
			CmdGroup:       "System Commands",
			Summary:        "soft-reboot                         Shut down and reboot userspace",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"exit": {
			CmdGroup:       "System Commands",
			Summary:        "exit [EXIT_CODE]                    Request user instance or container exit",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"switch-root": {
			CmdGroup:       "System Commands",
			Summary:        "switch-root [ROOT [INIT]]           Change to a different root file system",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"sleep": {
			CmdGroup:       "System Commands",
			Summary:        "Put the system to sleep (through one of the operations below)",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"suspend": {
			CmdGroup:       "System Commands",
			Summary:        "Suspend the system",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"hibernate": {
			CmdGroup:       "System Commands",
			Summary:        "Hibernate the system",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"hybrid-sleep": {
			CmdGroup:       "System Commands",
			Summary:        "Hibernate and suspend the system",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
		"suspend-then-hibernate": {
			CmdGroup:       "System Commands",
			Summary:        "Suspend the system, wake after a period of time, and hibernate",
			Description:    "",
			IsEnabled:      false,
			IsRootRequired: false,
		},
	},
}

func runTests() {

}

func INIT(debugMode utils.RunningMode, initMode utils.ToolsInitMode) {
	// Convert to JSON
	tmpSystemCtlSubCmds, err := json.MarshalIndent(systemCtlCmd.SubCommands, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}
	jsonSystemCtlSubCmds = string(tmpSystemCtlSubCmds)
	// For reference: how to convert back
	//     // Convert JSON to map
	//    var systemCtlCmd.SubCommands map[string]SingleSubCmd
	//    err := json.Unmarshal([]byte(jsonData), &systemCtlCmd.SubCommands)
	//    if err != nil {
	//        fmt.Println("Error unmarshalling JSON:", err)
	//        return
	//    }
	//    fmt.Println(systemCtlCmd.SubCommands)
	switch debugMode {
	case utils.Production:
		systemctlDebug = false
		if initMode == utils.Typed {
			for key := range systemCtlCmd.SubCommands {
				utils.AddToolToMCPServer(systemCtlCmd, jsonSystemCtlSubCmds, key, systemCtlCmd.SubCommands[key])
			}
		}
	case utils.Debug:
		systemctlDebug = true
		if initMode == utils.Typed {
			for key := range systemCtlCmd.SubCommands {
				utils.AddToolToMCPServer(systemCtlCmd, jsonSystemCtlSubCmds, key, systemCtlCmd.SubCommands[key])
			}
		}
	case utils.Test:
		systemctlDebug = true
		runTests()
	}
}
