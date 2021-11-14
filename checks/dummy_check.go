package checks

import (
	"context"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/rdar-lab/JCheck/common"
)

func GetDummyCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "DummyCheck",
		Group:       "DummyGroup",
		Description: "Nothing here",
		IsReadOnly:  true,
		CheckFunc: func(c context.Context) *common.CheckResult {
			shouldFail, _ := utils.GetBoolEnvValue("DummyFailure", false)
			if shouldFail {
				panic("test")
			}
			return &common.CheckResult{
				Success: true,
				Message: "Everything OK",
			}
		},
		CleanupFunc: func(c context.Context) error {
			return nil
		},
	}
}
