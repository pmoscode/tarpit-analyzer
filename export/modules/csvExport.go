package modules

import (
	"strconv"
	"strings"
	"tarpit-analyzer/database"
	"tarpit-analyzer/database/schemas"
	"tarpit-analyzer/export/helper"
	time2 "time"
)

type CSV struct {
	Separator string
}

func (r *CSV) Export(database *database.Database, start *time2.Time, end *time2.Time) (*[]string, error) {
	data := helper.QueryDataDB(database, start, end)

	result := make([]string, len(*data))

	for idx, dataItem := range *data {
		result[idx] = strings.Join(r.mapToCSV(dataItem), r.Separator)
	}

	return &result, nil
}

func (r *CSV) mapToCSV(data schemas.Data) []string {
	return []string{
		data.Begin.Format(time2.RFC3339),
		data.End.Format(time2.RFC3339),
		data.Ip,
		strconv.FormatFloat(float64(data.Duration), 'f', -1, 32),
	}
}
