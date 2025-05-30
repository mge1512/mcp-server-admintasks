package main

import (
	"mcp-server-admintasks/pkg/utils"
	"mcp-server-admintasks/pkg/snapper"
	"mcp-server-admintasks/pkg/systemctl"
	"mcp-server-admintasks/pkg/zypper"
)

func main() {
	utils.INIT(utils.Production)
	snapper.INIT(utils.Production)
	systemctl.INIT(utils.Production, utils.All)
	zypper.INIT(utils.Production, utils.All)
	utils.RUN()
}

