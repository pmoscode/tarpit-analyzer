package modules

import (
	"encoding/json"
	"endlessh-analyzer/database/schemas"
	"strconv"
	"strings"
	"time"
)

type JsonItem struct {
	Begin    string `json:"begin"`
	End      string `json:"end"`
	Ip       string `json:"ip"`
	Duration string `json:"duration"`
}

type JSON struct {
}

func (r *JSON) Export(data *[]schemas.Data) (*[]string, error) {
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

func (r JSON) mapToJsonItem(dataItem schemas.Data) *JsonItem {
	return &JsonItem{
		Begin:    dataItem.Begin.Format(time.RFC3339),
		End:      dataItem.End.Format(time.RFC3339),
		Ip:       dataItem.Ip,
		Duration: strconv.FormatFloat(float64(dataItem.Duration), 'f', -1, 32),
	}
}
