package checks

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/rdar-lab/JCheck/common"
)

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

			if resp.StatusCode != 200 {
				return "", errors.New("got http error for metrics")
			} else {
				//strResp := string(respBody)
				reader := bytes.NewReader(respBody)
				mf, err := common.ParseMF(reader)
				if err != nil {
					return "", err
				}

				monitorLabels := make(map[string]bool)
				monitorLabels["alertImpactAnalysis"] = true
				monitorLabels["ticketing"] = true
				monitorLabels["report"] = true
				monitorLabels["persistExistingContent"] = true

				queueMetric := *mf["queue_messages_total"]
				shouldFail := false
				for _, qm := range queueMetric.Metric {
					label := qm.Label[0].Value
					value := qm.Gauge.Value
					if monitorLabels[*label] {
						if *value > 100 {
							shouldFail = true
							break
						}
					}
				}

				if shouldFail {
					return "", errors.New(fmt.Sprintf("Rabbit MQ check has failed. Please check for overflowed queues"))
				}

				return fmt.Sprintf("Rabbit MQ check is valid"), nil
			}
		},
	}
}
