package commands

import (
	"errors"
	"github.com/fatih/color"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/rdar-lab/JCheck/common"
	"github.com/rodaine/table"
	"strconv"
)

func GetListCommand() components.Command {
	return components.Command{
		Name:        "list",
		Description: "Get the list of all implemented checks.",
		Aliases:     []string{},
		Arguments:   getListArguments(),
		Flags:       getListFlags(),
		EnvVars:     getListEnvVar(),
		Action: func(c *components.Context) error {
			return listCmd(c)
		},
	}
}

func getListArguments() []components.Argument {
	return []components.Argument{}
}

func getListFlags() []components.Flag {
	return []components.Flag{}
}

func getListEnvVar() []components.EnvVar {
	return []components.EnvVar{}
}

func listCmd(c *components.Context) error {
	if len(c.Arguments) != 0 {
		return errors.New("Wrong number of arguments. Expected: 0, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "Description", "Group", "Is Read Only")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, check := range common.GetRegistry().GetAllChecks() {
		tbl.AddRow(check.Name, check.Description, check.Group, check.IsReadOnly)
	}
	tbl.Print()
	return nil
}
