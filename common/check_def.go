package common

import (
	"context"
)

type CheckFuncDef func(c context.Context) (string, error)
type CleanupFuncDef func(c context.Context) error

type CheckDef struct {
	Name        string         `json:"name"`
	Group       string         `json:"group"`
	Description string         `json:"description"`
	IsReadOnly  bool           `json:"is_read_only"`
	CheckFunc   CheckFuncDef   `json:"-"`
	CleanupFunc CleanupFuncDef `json:"-"`
}
