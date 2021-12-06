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
	"github.com/rdar-lab/JCheck/common"
	"math"
	"net/http"
)

func GetXrayFreeDiskSpaceCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "XrayFreeDiskSpaceCheck",
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

			if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
				return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(respBody)))
			} else {
				//strResp := string(respBody)
				reader := bytes.NewReader(respBody)
				mf, err := common.ParseMF(reader)
				if err != nil {
					return "", err
				}
				diskFreeValue := *mf["app_disk_free_bytes"].GetMetric()[0].Gauge.Value
				diskFreeValueInGB := diskFreeValue / math.Pow(2, 30)

				shouldFail := diskFreeValueInGB < 100
				if shouldFail {
					return "", errors.New(fmt.Sprintf("Xray disk free space is lower than 100Gb (%.2f Gb)", diskFreeValueInGB))
				}

				return fmt.Sprintf("Xray free disk space is above 100Gb (%.2f Gb)", diskFreeValueInGB), nil
			}
		},
		CleanupFunc: func(c context.Context) error {
			return nil
		},
	}
}
