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
	"strings"
)

type checkResult struct {
	Success bool
	Message string
}

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
		components.StringFlag{
			Name:         "loop",
			Description:  "Loop over times.",
			DefaultValue: "1",
		},
	}
}

func getCheckEnvVar() []components.EnvVar {
	return []components.EnvVar{}
}

type checkConfiguration struct {
	what         string
	readOnlyMode bool
	loop         int
}

func checkCmd(c *components.Context) error {
	if len(c.Arguments) != 1 {
		return errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}

	var conf = new(checkConfiguration)
	conf.what = c.Arguments[0]
	conf.readOnlyMode = c.GetBoolFlagValue("readOnlyMode")
	loop, err := strconv.Atoi(c.GetStringFlagValue("loop"))
	if err != nil {
		return err
	}
	conf.loop = loop

	return doCheck(conf)
}

type resultPair struct {
	check  *common.CheckDef
	result *checkResult
}

func doCheck(conf *checkConfiguration) error {
	failureInd := false
	results := make([]*resultPair, 0, len(common.GetRegistry().GetAllChecks())*conf.loop)
	for i := 0; i < conf.loop; i++ {
		for _, check := range common.GetRegistry().GetAllChecks() {
			if conf.what == "" || conf.what == "ALL" || conf.what == check.Name || conf.what == check.Group {
				if check.IsReadOnly || !conf.readOnlyMode {
					result := runCheck(check)
					results = append(results,
						&resultPair{
							check:  check,
							result: result,
						},
					)
					if !result.Success {
						failureInd = true
					}
				}
			}
		}
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "Is Success", "Message")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, pair := range results {
		msg := pair.result.Message
		msg = strings.ReplaceAll(msg, "\n", " - ")

		tbl.AddRow(pair.check.Name, pair.result.Success, msg)
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

func runCheck(check *common.CheckDef) (result *checkResult) {
	context := context2.Background()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Check failed - Panic Detected: %v\n", r)
			result = &checkResult{
				Success: false,
				Message: "Check failure due to panic",
			}
		}
		if check.CleanupFunc != nil {
			err := check.CleanupFunc(context)
			if err != nil {
				fmt.Printf("Error on cleanup - %v\n", err)
			}
		}
	}()
	fmt.Printf("** Running check: %s...\n", check.Name)
	message, err := check.CheckFunc(context)

	if err == nil {
		result = &checkResult{
			Success: true,
			Message: message,
		}
	} else {
		result = &checkResult{
			Success: false,
			Message: err.Error(),
		}
	}

	fmt.Printf("Finished running check: %s, result=%v, message=%v\n", check.Name, result.Success, result.Message)
	return result
}
