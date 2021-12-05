package common

import "context"

type CheckFuncDef func(c context.Context) (string, error)
type CleanupFuncDef func(c context.Context) error

type CheckDef struct {
	Name        string
	Group       string
	Description string
	IsReadOnly  bool
	CheckFunc   CheckFuncDef
	CleanupFunc CleanupFuncDef
}
