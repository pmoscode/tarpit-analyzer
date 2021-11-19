package database

import (
	"endlessh-analyzer/database/schemas"
	"endlessh-analyzer/importData/structs"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	time2 "time"
)

type QueryParameters struct {
	StartDate   *time2.Time
	EndDate     *time2.Time
	SelectQuery *string
	WhereQuery  *WhereQuery
	OrderBy     *string
	Limit       *int
}

type WhereQuery struct {
	query      string
	parameters []interface{}
}

type DbData struct {
	database
}

func (r *DbData) timeRange(startDate *time2.Time, endDate *time2.Time) func(db *gorm.DB) *gorm.DB {
	if startDate == nil && endDate == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	if startDate == nil && endDate != nil {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("end <= ?", &endDate)
		}
	}

	if startDate != nil && endDate == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("begin >= ?", &endDate)
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where("begin >= ? AND end <= ?", &startDate, &endDate)
	}

}

func (r *DbData) query(queryParameter WhereQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(queryParameter.query, queryParameter.parameters)
	}
}

func (r *DbData) dbQuery(queryParameters QueryParameters) *gorm.DB {
	dbSub := r.db.Model(&schemas.Data{})
	if queryParameters.SelectQuery != nil {
		dbSub = dbSub.Select(*queryParameters.SelectQuery)
	}
	dbSub = dbSub.Scopes(r.timeRange(queryParameters.StartDate, queryParameters.EndDate))
	if queryParameters.WhereQuery != nil {
		dbSub = dbSub.Scopes(r.query(*queryParameters.WhereQuery))
	}
	if queryParameters.OrderBy != nil {
		dbSub = dbSub.Order(*queryParameters.OrderBy)
	}

	return dbSub
}

func (r *DbData) ExecuteQueryGetList(queryParameters QueryParameters) ([]schemas.Data, DbResult) {
	var data []schemas.Data

	dbSub := r.dbQuery(queryParameters)
	result := dbSub.Find(&data)

	if result.RowsAffected == 0 {
		return data, DbRecordNotFound
	}

	return data, DbOk
}

func (r *DbData) ExecuteQueryGetFirst(queryParameters QueryParameters) (schemas.Data, DbResult) {
	var data schemas.Data

	dbSub := r.dbQuery(queryParameters)
	result := dbSub.First(&data)

	if result.RowsAffected == 0 {
		return data, DbRecordNotFound
	}

	return data, DbOk
}

func (r *DbData) ExecuteQueryGetAggregator(queryParameters QueryParameters) (float64, DbResult) {
	var data float64

	dbSub := r.dbQuery(queryParameters)
	dbSub.First(&data)

	return data, DbOk
}

func (r *DbData) SaveData(data *[]schemas.Data) (DbResult, error) {
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

func CreateDbData(debug bool) (DbData, error) {
	db := DbData{}
	db.initDatabase("data", debug)

	err := db.db.AutoMigrate(&schemas.Data{})
	if err != nil {
		log.Errorln(err)
		return DbData{}, err
	}

	return db, nil
}
