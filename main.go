package main

import (
	"github.com/jfrog/jfrog-cli-core/v2/plugins"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/rdar-lab/JCheck/commands"
)

func main() {
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
