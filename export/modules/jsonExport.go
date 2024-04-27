package modules

import (
	"encoding/json"
	"strconv"
	"strings"
	"tarpit-analyzer/database"
	"tarpit-analyzer/database/schemas"
	"tarpit-analyzer/export/helper"
	time2 "time"
)

type JsonItem struct {
	Begin    string `json:"begin"`
	End      string `json:"end"`
	Ip       string `json:"ip"`
	Duration string `json:"duration"`
}

type JSON struct {
}

func (r *JSON) Export(database *database.Database, start *time2.Time, end *time2.Time) (*[]string, error) {
	data := helper.QueryDataDB(database, start, end)

	result := make([]string, 0)

	result = append(result, "[")

	for _, dataItem := range *data {
		mappedData := r.mapToJsonItem(dataItem)
		jsonData, err := json.Marshal(*mappedData)
		if err != nil {
			return nil, err
		}
		result = append(result, string(jsonData)+",")
	}

	result[len(result)-1] = strings.TrimRight(result[len(result)-1], ",")

	result = append(result, "]")

	return &result, nil
}

func (r *JSON) mapToJsonItem(dataItem schemas.Data) *JsonItem {
	return &JsonItem{
		Begin:    dataItem.Begin.Format(time2.RFC3339),
		End:      dataItem.End.Format(time2.RFC3339),
		Ip:       dataItem.Ip,
		Duration: strconv.FormatFloat(float64(dataItem.Duration), 'f', -1, 32),
	}
}
