package common

import (
	"context"
	"errors"
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

func GetStateMapFromContext(ctx context.Context) (map[string]interface{}, error) {
	state := ctx.Value("State")
	stateMap, ok := state.(map[string]interface{})
	if !ok {
		return nil, errors.New("error with testing platform")
	}
	return stateMap, nil
}
