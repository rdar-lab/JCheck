package checks

import (
	"context"
	"errors"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/rdar-lab/JCheck/common"
)

func GetSelfCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "SelfCheck",
		Group:       "Self",
		Description: "A sanity check that should pass",
		IsReadOnly:  true,
		CheckFunc: func(c context.Context) (string, error) {
			shouldPanic, _ := utils.GetBoolEnvValue("PanicTest", false)
			if shouldPanic {
				panic("Panic indication detected")
			}
			shouldFail, _ := utils.GetBoolEnvValue("FailureTest", false)
			if shouldFail {
				return "", errors.New("failure indication detected")
			}
			return "Self check passed", nil
		},
	}
}
