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
	"net/http"
)

func GetXrayDbConnectionPoolCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "XrayDbConnectionPoolCheck",
		Group:       "Xray",
		Description: "Performs a check that DB connection pool is not maxed",
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
				dbUsed := *mf["db_connection_pool_in_use_total"].GetMetric()[0].Gauge.Value
				dbMax := *mf["db_connection_pool_max_open_total"].GetMetric()[0].Gauge.Value

				shouldFail := dbUsed == dbMax
				if shouldFail {
					return "", errors.New(fmt.Sprintf("Xray DB connection pool is full (%.f/%.f connections)", dbUsed, dbMax))
				}

				return fmt.Sprintf("Xray DB connection pool has available connections (%.f/%.f connections)", dbUsed, dbMax), nil
			}
		},
	}
}
