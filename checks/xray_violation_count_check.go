package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands"
	artUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/rdar-lab/JCheck/common"
	"net/http"
	"time"
)

const LIMIT = 10000

func GetXrayViolationsCountCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "XrayViolationCountCheck",
		Group:       "Xray",
		Description: "Performs a check that checks that Xray is not generating too many violations",
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

			url := clientutils.AddTrailingSlashIfNeeded(serverConf.XrayUrl) + "api/v1/violations"

			type filters struct {
				CreatedFrom time.Time `json:"created_from,omitempty"`
			}

			type pagination struct {
				Limit  int `json:"limit,omitempty"`
				Offset int `json:"offset,omitempty"`
			}

			requestBody := struct {
				Filters    filters    `json:"filters,omitempty"`
				Pagination pagination `json:"pagination,omitempty"`
			}{
				Filters: filters{
					CreatedFrom: time.Now().AddDate(0, 0, -1),
				},
				Pagination: pagination{
					Limit:  1,
					Offset: 1,
				},
			}
			content, err := json.Marshal(requestBody)
			if err != nil {
				return "", err
			}

			artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
			resp, body, err := xrayServiceMgr.Client().SendPost(url, content, &httpClientsDetails)
			if err != nil {
				return "", err
			}
			if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
				return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
			}

			response := struct {
				TotalViolations int `json:"total_violations,omitempty"`
			}{}

			err = json.Unmarshal(body, &response)

			if err != nil {
				return "", errors.New("failed unmarshalling violations response")
			}

			if response.TotalViolations <= LIMIT {

				return fmt.Sprintf("detected %d violations in last 24 hours", response.TotalViolations), nil
			} else {
				return "", errors.New(fmt.Sprintf("detected %d violations in last 24 hours", response.TotalViolations))
			}

		},
	}
}
