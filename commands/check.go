package commands

import (
	context2 "context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/rdar-lab/JCheck/common"
	"github.com/rodaine/table"
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
	return []components.Argument{
		{
			Name:        "what",
			Description: "The names of the checks you want to run. It can be a group name, check name or ALL for all",
		},
	}
}

func getCheckFlags() []components.Flag {
	return []components.Flag{
		components.BoolFlag{
			Name:         "readOnlyMode",
			Description:  "Only run checks which are read only.",
			DefaultValue: false,
		},
	}
}

func getCheckEnvVar() []components.EnvVar {
	return []components.EnvVar{}
}

type checkConfiguration struct {
	what         string
	readOnlyMode bool
}

func checkCmd(c *components.Context) error {
	if len(c.Arguments) != 1 {
		return errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}

	var conf = new(checkConfiguration)
	conf.what = c.Arguments[0]
	conf.readOnlyMode = c.GetBoolFlagValue("readOnlyMode")
	return doCheck(conf)
}

func doCheck(conf *checkConfiguration) error {
	failureInd := false
	results := make(map[*common.CheckDef]*common.CheckResult)
	for _, check := range common.GetRegistry().GetAllChecks() {
		if conf.what == "" || conf.what == "ALL" || conf.what == check.Name || conf.what == check.Group {
			if check.IsReadOnly || !conf.readOnlyMode {
				result := runCheck(check)
				results[check] = result
				if !result.Success {
					failureInd = true
				}
			}
		}
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "Is Success", "Message")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for check, result := range results {
		tbl.AddRow(check.Name, result.Success, result.Message)
	}
	fmt.Println()
	fmt.Println()
	tbl.Print()
	fmt.Println()
	fmt.Println()

	if failureInd {
		return errors.New("Errors detected")
	} else {
		return nil
	}
}

func runCheck(check *common.CheckDef) (result *common.CheckResult) {
	context := context2.Background()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Check failed - Panic Detected: %v\n", r)
			result = &common.CheckResult{
				Success: false,
				Message: "Check failure due to panic",
			}
		}
		err := check.CleanupFunc(context)
		if err != nil {
			fmt.Printf("Error on cleanup - %v\n", err)
		}
	}()
	fmt.Printf("** Running check: %s...\n", check.Name)
	result = check.CheckFunc(context)
	fmt.Printf("Finished running check: %s, result=%v, message=%v\n", check.Name, result.Success, result.Message)
	return result
}
