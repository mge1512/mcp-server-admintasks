package main

import (
	"mcp-server-admintasks/pkg/snapper"
	"mcp-server-admintasks/pkg/systemctl"
	"mcp-server-admintasks/pkg/utils"
	"mcp-server-admintasks/pkg/zypper"
)

func main() {
	utils.INIT(utils.Debug)
	snapper.INIT(utils.Production)
	systemctl.INIT(utils.Test, utils.Typed)
	zypper.INIT(utils.Test, utils.Typed)
	utils.RUN()
}
