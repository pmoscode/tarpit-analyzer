package modules

import (
	"endlessh-analyzer/database/schemas"
)

type JSON struct {
}

func (r *JSON) Export(data *[]schemas.Data) (*[]string, error) {
	result := make([]string, len(*data))

	return &result, nil
}
