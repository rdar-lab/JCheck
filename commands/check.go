package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/rdar-lab/JCheck/common"
	"github.com/rodaine/table"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var whitespaceRe = regexp.MustCompile(`\r\n|[\r\n\v\f\x{0085}\x{2028}\x{2029}]`)

type checkResult struct {
	Time    time.Time `json:"time"`
	Success bool      `json:"is_success"`
	Message string    `json:"message"`
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
		components.StringFlag{
			Name:         "loopSleep",
			Description:  "Sleep time (in seconds) between loops.",
			DefaultValue: "0",
		},
		components.BoolFlag{
			Name:         "json",
			Description:  "Return JSON result",
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
	loop         int
	loopSleep    int
	json         bool
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
	loopSleep, err := strconv.Atoi(c.GetStringFlagValue("loopSleep"))
	if err != nil {
		return err
	}
	conf.loopSleep = loopSleep
	conf.json = c.GetBoolFlagValue("json")

	return doCheck(conf)
}

type resultPair struct {
	Check  *common.CheckDef `json:"check_def"`
	Result *checkResult     `json:"result"`
}

func doCheck(conf *checkConfiguration) error {
	failureInd := false
	results := make([]*resultPair, 0, len(common.GetRegistry().GetAllChecks())*conf.loop)
	checksCount := 0

	for i := 0; i < conf.loop; i++ {
		for _, check := range common.GetRegistry().GetAllChecks() {
			if conf.what == "" ||
				strings.ToUpper(conf.what) == "ALL" ||
				strings.Contains(strings.ToLower(check.Name), strings.ToLower(conf.what)) ||
				strings.Contains(strings.ToLower(check.Group), strings.ToLower(conf.what)) {
				if check.IsReadOnly || !conf.readOnlyMode {
					result := runCheck(check)
					results = append(results,
						&resultPair{
							Check:  check,
							Result: result,
						},
					)
					if !result.Success {
						failureInd = true
					}
					checksCount++
				}
			}
		}
		if conf.loopSleep > 0 {
			time.Sleep(time.Second * time.Duration(conf.loopSleep))
		}
	}

	if conf.json {
		err := outputResultAsJson(results)
		if err != nil {
			return err
		}
	} else {
		outputResultTable(results)
	}

	if checksCount == 0 {
		return errors.New("no checks performed")
	}

	if failureInd {
		return errors.New("failures detected")
	} else {
		return nil
	}
}

func outputResultAsJson(results []*resultPair) error {
	jsonData, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		return err
	}

	fmt.Printf(string(jsonData))
	return nil
}

func outputResultTable(results []*resultPair) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Time", "Name", "Failure Ind", "Message")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithPadding(5)

	for _, pair := range results {
		msg := pair.Result.Message
		msg = whitespaceRe.ReplaceAllString(msg, " ")

		failureStr := ""
		if !pair.Result.Success {
			failureStr = "FAIL"
		}

		tbl.AddRow(pair.Result.Time.Format(time.StampMilli), pair.Check.Name, failureStr, msg)
	}
	fmt.Println()
	fmt.Println()
	tbl.Print()
	fmt.Println()
	fmt.Println()
}

func runCheck(check *common.CheckDef) (result *checkResult) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "State", make(map[string]interface{}))
	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Sprintf("Check failed - Panic Detected: %v\n", r))
			result = &checkResult{
				Success: false,
				Message: "Check failure due to panic",
			}
		}
		if check.CleanupFunc != nil {
			err := check.CleanupFunc(ctx)
			if err != nil {
				log.Error(fmt.Sprintf("Error on cleanup - %v\n", err))
			}
		}
	}()
	log.Info(fmt.Sprintf("** Running check: %s...\n", check.Name))
	message, err := check.CheckFunc(ctx)

	if err == nil {
		result = &checkResult{
			Time:    time.Now(),
			Success: true,
			Message: message,
		}
	} else {
		result = &checkResult{
			Time:    time.Now(),
			Success: false,
			Message: err.Error(),
		}
	}

	log.Info(fmt.Sprintf("Finished running check: %s, result=%v, message=%v\n", check.Name, result.Success, result.Message))
	return result
}
