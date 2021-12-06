package checks

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/rdar-lab/JCheck/common"
)

func GetRTHasRepositoriesCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "RTHasRepositoriesCheck",
		Group:       "Artifactory",
		Description: "Performs a check that validates that RT has configured repositories",
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
			repos, err := serviceManager.GetAllRepositories()
			if err != nil {
				return "", err
			}
			if len(*repos) == 0 {
				return "", errors.New("detected 0 repositories")
			} else {
				return fmt.Sprintf("detected %d repositories", len(*repos)), nil
			}
		},
	}
}
