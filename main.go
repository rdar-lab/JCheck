package main

import (
	"github.com/jfrog/jfrog-cli-core/v2/plugins"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/rdar-lab/JCheck/checks"
	"github.com/rdar-lab/JCheck/commands"
	"github.com/rdar-lab/JCheck/common"
)

func main() {
	registerChecks()
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "JCheck"
	app.Description = " A Micro-UTP, plug-able sanity checker for any on-prem JFrog platform instance."
	app.Version = "v0.1.0"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		commands.GetListCommand(),
		commands.GetCheckCommand(),
	}
}

// TODO: Add ability to inject external checks via configuration
func registerChecks() {
	common.GetRegistry().Register(checks.GetSelfCheck())
	common.GetRegistry().Register(checks.GetRTPingCheck())
	common.GetRegistry().Register(checks.GetXrayPingCheck())
}
