package main

import (
	"mcp-server-admintasks/pkg/systemctl"
	"mcp-server-admintasks/pkg/utils"
	"mcp-server-admintasks/pkg/zypper"
)

func main() {
	utils.INIT(utils.Debug)
	systemctl.INIT(utils.Test, utils.Typed)
	zypper.INIT(utils.Test, utils.Typed)
	utils.RUN()
}
