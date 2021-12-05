package checks

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/rdar-lab/JCheck/common"
	"io"
	"math"
)

func parseMF(reader io.Reader) (map[string]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, err
	}
	return mf, nil
}

func GetXrayMertricsFreeDiskSpaceCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "GetXrayMertricsFreeDiskSpaceCheck",
		Group:       "Xray",
		Description: "Performs a check that free disk space is above 1Gb",
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
				mf, err := parseMF(reader)
				if err != nil {
					return "", err
				}
				diskFreeValue := *mf["app_disk_free_bytes"].GetMetric()[0].Gauge.Value
				shouldFail := diskFreeValue < math.Pow(2, 30) // 1G
				if shouldFail {
					return "", errors.New(fmt.Sprintf("Xray disk free space is lower than 1G (%.f bytes)", diskFreeValue))
				}

				return fmt.Sprintf("Xray free disk space is above 1G (%.f bytes)", diskFreeValue), nil
			}
		},
		CleanupFunc: func(c context.Context) error {
			return nil
		},
	}
}
