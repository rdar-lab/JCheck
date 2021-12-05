package checks

import (
	"context"
	"errors"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/rdar-lab/JCheck/common"
)

func GetRTPingCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "RTPingCheck",
		Group:       "Artifactory",
		Description: "Performs a check that validates that a ping to RT works",
		IsReadOnly:  true,
		CheckFunc: func(c context.Context) (string, error) {
			rtDetails, err := config.GetDefaultServerConf()
			if err != nil {
				return "", err
			}
			serviceManager, err := utils.CreateServiceManager(rtDetails, -1, false)
			if err != nil {
				return "", err
			}
			resp, err := serviceManager.Ping()
			if err != nil {
				return "", err
			}
			respStr := string(resp)
			if respStr != "OK" {
				return "", errors.New("got unexpected response: " + respStr)
			} else {
				return "RT Ping was successful", nil
			}
		},
		CleanupFunc: func(c context.Context) error {
			return nil
		},
	}
}
