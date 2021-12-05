package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/rdar-lab/JCheck/common"
	"net/http"
)

func GetRTHasProjectsCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "RTHasProjectsCheck",
		Group:       "Artifactory",
		Description: "Performs a check that validates that RT has configured projects",
		IsReadOnly:  true,
		CheckFunc: func(c context.Context) (string, error) {
			serverConf, err := config.GetDefaultServerConf()
			if err != nil {
				return "", err
			}
			serviceManager, err := utils.CreateServiceManager(serverConf, -1, false)
			if err != nil {
				return "", err
			}

			httpClientsDetails := serviceManager.GetConfig().GetServiceDetails().CreateHttpClientDetails()

			url := clientutils.AddTrailingSlashIfNeeded(serverConf.AccessUrl) + "api/v1/projects"

			resp, body, _, err := serviceManager.Client().SendGet(url, true, &httpClientsDetails)
			if err != nil {
				return "", err
			}
			if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
				return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
			}

			projects := make([]struct{}, 0)

			err = json.Unmarshal(body, &projects)

			if err != nil {
				return "", errors.New("failed unmarshalling projects response")
			}

			return fmt.Sprintf("detected %d projects", len(projects)), nil
		},
	}
}
