package database

import (
	"crypto"
	"encoding/hex"
	"endlessh-analyzer/database/schemas"
	"endlessh-analyzer/importData/structs"
	"fmt"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
	"strings"
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
	Query      string
	Parameters interface{}
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
		return db.Where(queryParameter.Query, queryParameter.Parameters)
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
	*data = removeDuplicateValues(data)
	sum := len(*data)
	bar := progressbar.Default(int64(sum), "Saving...")

	errTx := r.db.Transaction(func(tx *gorm.DB) error {
		for _, d := range *data {
			res := tx.Create(&d)
			bar.Add(1)
			if res.Error != nil {
				if !strings.Contains(res.Error.Error(), "UNIQUE constraint failed") { // Theoretically should not happen
					return res.Error
				}
			}
		}
		bar.Finish()

		return nil
	})
	if errTx != nil {
		return DbError, errTx
	}

	return DbOk, nil
}

func (r *DbData) MapToData(importItem structs.ImportItem) schemas.Data {
	return schemas.Data{
		ID: hex.EncodeToString(hash(importItem.Begin, importItem.End, importItem.Ip)),
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

func (r *DbData) Map(vs *[]structs.ImportItem, f func(importItem structs.ImportItem) schemas.Data) *[]schemas.Data {
	vsm := make([]schemas.Data, len(*vs))
	for i, v := range *vs {
		vsm[i] = f(v)
	}

	return &vsm
}

func hash(objs ...interface{}) []byte {
	digester := crypto.MD5.New()
	for _, ob := range objs {
		_, err1 := fmt.Fprint(digester, reflect.TypeOf(ob))
		if err1 != nil {
			return nil
		}
		_, err2 := fmt.Fprint(digester, ob)
		if err2 != nil {
			return nil
		}
	}
	return digester.Sum(nil)
}

func removeDuplicateValues(dataSlice *[]schemas.Data) []schemas.Data {
	keys := make(map[schemas.Data]bool)
	var list []schemas.Data

	for _, entry := range *dataSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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
