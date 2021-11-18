package database

import (
	"endlessh-analyzer/database/schemas"
	"endlessh-analyzer/importData/structs"
	log "github.com/sirupsen/logrus"
	time2 "time"
)

type DbData struct {
	database
}

func (r DbData) GetData(startDate time2.Time, endDate time2.Time) ([]schemas.Data, DbResult) {
	var data []schemas.Data

	result := r.db.Where("begin >= ? AND end <= ?", startDate, endDate).Find(&data)

	if result.RowsAffected == 0 {
		return []schemas.Data{}, DbRecordNotFound
	}

	return data, DbOk
}

func (r DbData) SaveData(data *[]schemas.Data) (DbResult, error) {
	result := r.db.CreateInBatches(data, 100)

	if result.Error != nil {
		return DbError, result.Error
	}

	return DbOk, nil
}

func (r *DbData) MapToData(importItem structs.ImportItem) schemas.Data {
	return schemas.Data{
		ImportItem: structs.ImportItem{
			Begin:    importItem.Begin,
			End:      importItem.End,
			Duration: importItem.Duration,
			Ip:       importItem.Ip,
			Success:  importItem.Success,
		},
	}
}

func (r *DbData) MapToImportItem(data schemas.Data) structs.ImportItem {
	return structs.ImportItem{
		Begin:    data.Begin,
		End:      data.End,
		Duration: data.Duration,
		Ip:       data.Ip,
		Success:  data.Success,
	}
}

func (r DbData) Map(vs *[]structs.ImportItem, f func(importItem structs.ImportItem) schemas.Data) *[]schemas.Data {
	vsm := make([]schemas.Data, len(*vs))
	for i, v := range *vs {
		vsm[i] = f(v)
	}

	return &vsm
}

func CreateDbData() (DbData, error) {
	db := DbData{}
	db.initDatabase("data")

	err := db.db.AutoMigrate(&schemas.Data{})
	if err != nil {
		log.Errorln(err)
		return DbData{}, err
	}

	return db, nil
}
