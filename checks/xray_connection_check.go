package checks

import (
	"context"
	"errors"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands"
	"github.com/rdar-lab/JCheck/common"
)

func GetXrayConnectionCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "XrayConnectionCheck",
		Group:       "Xray",
		Description: "Performs a check that validates that a connection to XRAY works",
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

			version, err := xrayServiceMgr.GetVersion()

			if err != nil {
				return "", err
			}

			if version == "" {
				return "", errors.New("empty version returned by Xray")
			} else {
				return "Xray version " + version + " was detected", nil
			}
		},
	}
}
