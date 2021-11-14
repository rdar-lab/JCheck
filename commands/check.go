package commands

import (
	"errors"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"strconv"
)

func GetCheckCommand() components.Command {
	return components.Command{
		Name:        "check",
		Description: "Run the checks on the platform.",
		Aliases:     []string{"run", "exec"},
		Arguments:   getCheckArguments(),
		Flags:       getCheckFlags(),
		EnvVars:     getCheckEnvVar(),
		Action: func(c *components.Context) error {
			return checkCmd(c)
		},
	}
}

func getCheckArguments() []components.Argument {
	return []components.Argument{}
}

func getCheckFlags() []components.Flag {
	return []components.Flag{}
}

func getCheckEnvVar() []components.EnvVar {
	return []components.EnvVar{}
}

func checkCmd(c *components.Context) error {
	if len(c.Arguments) != 0 {
		return errors.New("Wrong number of arguments. Expected: 0, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}

	return errors.New("Not implemented yet")
}
