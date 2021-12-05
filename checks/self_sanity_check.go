package checks

import (
	"context"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/rdar-lab/JCheck/common"
)

func GetSelfCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "SelfCheck",
		Group:       "Self",
		Description: "A sanity check that should pass",
		IsReadOnly:  true,
		CheckFunc: func(c context.Context) *common.CheckResult {
			shouldPanic, _ := utils.GetBoolEnvValue("PanicTest", false)
			if shouldPanic {
				panic("Panic indication detected")
			}
			shouldFail, _ := utils.GetBoolEnvValue("FailureTest", false)
			if shouldFail {
				return &common.CheckResult{
					Success: false,
					Message: "Failure indication detected",
				}
			}
			return &common.CheckResult{
				Success: true,
				Message: "Self check passed",
			}
		},
		CleanupFunc: func(c context.Context) error {
			return nil
		},
	}
}
