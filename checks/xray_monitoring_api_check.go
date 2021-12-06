package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/rdar-lab/JCheck/common"
	"net/http"
	"time"
)

func GetXrayMonitoringAPICheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "XrayMonitoringAPICheck",
		Group:       "Xray",
		Description: "Performs a check that calls XRAY monitoring API",
		IsReadOnly:  true,
		CheckFunc: func(c context.Context) (string, error) {

			serverConf, err := config.GetDefaultServerConf()
			if err != nil {
				return "", err
			}

			if serverConf.XrayUrl == "" {
				return "", errors.New("xray service is not configured")
			}

			xrayServiceMgr, err := commands.CreateXrayServiceManager(serverConf)
			if err != nil {
				return "", err
			}

			httpClientsDetails := xrayServiceMgr.Config().GetServiceDetails().CreateHttpClientDetails()

			url := clientutils.AddTrailingSlashIfNeeded(serverConf.XrayUrl) + "api/v1/monitor"

			resp, body, _, err := xrayServiceMgr.Client().SendGet(url, true, &httpClientsDetails)
			if err != nil {
				return "", err
			}
			if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
				return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
			}

			type MonitoringSystemProblem struct {
				Severity string `json:"severity"`
				// Same problem can be potentially experienced by more than one service
				Services       []string  `json:"services"`
				Problem        string    `json:"problem"`
				ProblemTime    time.Time `json:"problem_time"`
				ShouldShowInUI bool      `json:"should_show_in_ui"`
			}

			type MonitoringSystemStatus struct {
				Problems []MonitoringSystemProblem `json:"problems"`
			}

			monitoringStatus := MonitoringSystemStatus{}

			err = json.Unmarshal(body, &monitoringStatus)

			if err != nil {
				return "", errors.New("failed unmarshalling monitoring response")
			}

			if len(monitoringStatus.Problems) > 0 {
				for _, problem := range monitoringStatus.Problems {
					log.Warn(fmt.Sprintf("Problem detected: %s, services=%v, time=%v\n", problem.Problem, problem.Services, problem.ProblemTime))
				}
				return "", errors.New(fmt.Sprintf("detected %d problems", len(monitoringStatus.Problems)))
			} else {
				return "No problems found", nil
			}

		},
	}
}
