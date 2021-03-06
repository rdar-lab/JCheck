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
	app.Name = "jcheck"
	app.Description = " A Micro-UTP, plug-able sanity checker for any on-prem JFrog platform instance."
	app.Version = "v1.0.0"
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
	common.GetRegistry().Register(checks.GetRTConnectionCheck())
	common.GetRegistry().Register(checks.GetRTHasRepositoriesCheck())
	//common.GetRegistry().Register(checks.GetRTHasProjectsCheck()) // Disabled due to permission issues
	common.GetRegistry().Register(checks.GetXrayConnectionCheck())
	common.GetRegistry().Register(checks.GetXrayHasPoliciesCheck())
	common.GetRegistry().Register(checks.GetXrayHasWatchesCheck())
	common.GetRegistry().Register(checks.GetXrayHasIndexedResourcesCheck())
	common.GetRegistry().Register(checks.GetXrayViolationsCountCheck())
	common.GetRegistry().Register(checks.GetXrayFreeDiskSpaceCheck())
	common.GetRegistry().Register(checks.GetXrayMonitoringAPICheck())
	common.GetRegistry().Register(checks.GetXrayDbConnectionPoolCheck())
	common.GetRegistry().Register(checks.GetRTDeployCheck())
	common.GetRegistry().Register(checks.GetXrayRabbitMQCheck())
}
