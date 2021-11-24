package modules

import (
	"endlessh-analyzer/database/schemas"
)

type KML struct {
}

func (r *KML) Export(data *[]schemas.Data) (*[]string, error) {
	result := make([]string, len(*data))

	return &result, nil
}
