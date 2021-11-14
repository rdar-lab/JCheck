package common

import "context"

type CheckResult struct {
	Success bool
	Message string
}

type CheckFuncDef func(c context.Context) *CheckResult
type CleanupFuncDef func(c context.Context) error

type CheckDef struct {
	Name        string
	Group       string
	Description string
	IsReadOnly  bool
	CheckFunc   CheckFuncDef
	CleanupFunc CleanupFuncDef
}
