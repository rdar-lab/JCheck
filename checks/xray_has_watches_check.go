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
	"github.com/rdar-lab/JCheck/common"
	"net/http"
)

func GetXrayHasWatchesCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "XrayHasWatchesCheck",
		Group:       "Xray",
		Description: "Performs a check that validates that XRAY has configured watches",
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

			url := clientutils.AddTrailingSlashIfNeeded(serverConf.XrayUrl) + "api/v2/watches"

			resp, body, _, err := xrayServiceMgr.Client().SendGet(url, true, &httpClientsDetails)
			if err != nil {
				return "", err
			}
			if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
				return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
			}

			watches := make([]struct{}, 0)

			err = json.Unmarshal(body, &watches)

			if err != nil {
				return "", errors.New("failed unmarshalling watches response")
			}

			if len(watches) > 0 {

				return fmt.Sprintf("detected %d watches", len(watches)), nil
			} else {
				return "", errors.New("detected 0 watches")
			}

		},
	}
}
