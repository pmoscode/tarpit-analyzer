package modules

import (
	"endlessh-analyzer/database/schemas"
	"strconv"
	"strings"
	"time"
)

type CSV struct {
	Separator string
}

func (r *CSV) Export(data *[]schemas.Data) (*[]string, error) {
	result := make([]string, len(*data))

	for idx, dataItem := range *data {
		result[idx] = strings.Join(r.mapToCSV(dataItem), r.Separator)
	}

	return &result, nil
}

func (r CSV) mapToCSV(data schemas.Data) []string {
	return []string{
		data.Begin.Format(time.RFC3339),
		data.End.Format(time.RFC3339),
		data.Ip,
		strconv.FormatFloat(float64(data.Duration), 'f', -1, 32),
	}
}
