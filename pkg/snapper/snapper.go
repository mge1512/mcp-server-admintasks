package snapper

import (
        "mcp-server-admintasks/pkg/utils"
)

var snapperDebug bool

func addToolsToMCPServer() {

}

func INIT(mode utils.RunningMode) {
	switch mode {
		case utils.Production:
			snapperDebug = false
			addToolsToMCPServer()
		case utils.Debug:
			snapperDebug = true
			addToolsToMCPServer()
		case utils.Test:
			snapperDebug = true
	}
}

