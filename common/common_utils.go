package common

import (
	"context"
	"errors"
)

func GetStateMapFromContext(ctx context.Context) (map[string]interface{}, error) {
	state := ctx.Value("State")
	stateMap, ok := state.(map[string]interface{})
	if !ok {
		return nil, errors.New("error with testing platform")
	}
	return stateMap, nil
}
