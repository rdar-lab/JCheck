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

func GetXrayHasPoliciesCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "XrayHasPoliciesCheck",
		Group:       "Xray",
		Description: "Performs a check that validates that XRAY has configured policies",
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

			url := clientutils.AddTrailingSlashIfNeeded(serverConf.XrayUrl) + "api/v2/policies"

			resp, body, _, err := xrayServiceMgr.Client().SendGet(url, true, &httpClientsDetails)
			if err != nil {
				return "", err
			}
			if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
				return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
			}

			policies := make([]struct{}, 0)

			err = json.Unmarshal(body, &policies)

			if err != nil {
				return "", errors.New("failed unmarshalling policies response")
			}

			if len(policies) > 0 {

				return fmt.Sprintf("detected %d policies", len(policies)), nil
			} else {
				return "", errors.New("detected 0 policies")
			}

		},
	}
}
