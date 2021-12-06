package checks

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/rdar-lab/JCheck/common"
	"net/http"
)

const THRESHOLD = 1000

func GetXrayRabbitMQCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "XrayRabbitMQCheck",
		Group:       "Xray",
		Description: "Performs a check that critical queues are not overflowed",
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

			xrayUrl := serverConf.GetXrayUrl() + "/api/v1/metrics"
			details := httputils.HttpClientDetails{
				User:     serverConf.GetUser(),
				Password: serverConf.GetPassword(),
			}
			resp, respBody, _, err := xrayServiceMgr.Client().SendGet(xrayUrl, true, &details)

			if err != nil {
				return "", err
			}

			if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
				return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(respBody)))
			} else {
				//strResp := string(respBody)
				reader := bytes.NewReader(respBody)
				mf, err := common.ParseMF(reader)
				if err != nil {
					return "", err
				}

				totalCount := 0.0
				queueMetric := *mf["queue_messages_total"]
				for _, qm := range queueMetric.Metric {
					label := *qm.Label[0].Value
					value := *qm.Gauge.Value

					log.Info(fmt.Sprintf("Queue %s - size %.f", label, value))

					if value > THRESHOLD {
						return "", errors.New(fmt.Sprintf("Queue %s reached size of %.f", label, value))
					}
					totalCount += value
				}

				return fmt.Sprintf("Total number of messages = %.f", totalCount), nil
			}
		},
	}
}
