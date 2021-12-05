package checks

import (
	"context"
	"errors"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/rdar-lab/JCheck/common"
)

func GetRTConnectionCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "RTConnectionCheck",
		Group:       "Artifactory",
		Description: "Performs a check that validates that a connection to RT works",
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
			version, err := serviceManager.GetVersion()
			if err != nil {
				return "", err
			}
			if version == "" {
				return "", errors.New("empty version returned by RT")
			} else {
				return "RT version " + version + " was detected", nil
			}
		},
	}
}
