package database

import (
	"database/sql"
	"endlessh-analyzer/database/countryGeoLocationData"
	"endlessh-analyzer/database/schemas"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	time2 "time"
)

type DbResult int

const (
	DbOk DbResult = iota
	DbRecordNotFound
	DbError
)

type QueryParameters struct {
	StartDate   *time2.Time
	EndDate     *time2.Time
	Distinct    bool
	SelectQuery *string
	WhereQuery  *[]WhereQuery
	JoinQuery   *string
	GroupBy     *string
	OrderBy     *string
	Limit       *int
}

type WhereQuery struct {
	Query      string
	Parameters interface{}
}

type Database struct {
	db *gorm.DB
}

func (r *Database) initDatabase(dbFilename string, debug bool) {

	logLevel := logger.Silent
	if debug {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time2.Second, // Slow SQL threshold
			LogLevel:                  logLevel,     // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,        // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(dbFilename+".db"), &gorm.Config{
		Logger:      newLogger,
		PrepareStmt: true,
	})
	if err != nil {
		logrus.Errorln(err)
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&schemas.Location{}, &schemas.Data{}, &schemas.CountryGeoLocation{})
	if err != nil {
		logrus.Errorln(err)
		panic("Could not migrate schema for 'Location' / 'Data'")
	}

	r.db = db

	r.initCountryCodes()
}

func (r *Database) DbRawQuery(query string, parameters ...interface{}) (*sql.Rows, error) {
	rows, err := r.db.Raw(query, parameters).Rows()
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *Database) DbQuery(schema interface{}, queryParameters QueryParameters) (*sql.Rows, error) {
	dbSub := r.internalQuery(schema, queryParameters)
	rows, err := dbSub.Rows()
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *Database) ScanToStruct(rows *sql.Rows, model interface{}) error {
	err := r.db.ScanRows(rows, model)
	if err != nil {
		return err
	}

	return nil
}

func (r *Database) ExecuteQueryGetList(model interface{}, target interface{}, queryParameters QueryParameters) DbResult {
	dbSub := r.internalQuery(model, queryParameters)
	result := dbSub.Find(target)

	if result.Error != nil {
		logrus.Errorln(result.Error)
		return DbError
	}

	if result.RowsAffected == 0 {
		return DbRecordNotFound
	}

	return DbOk
}

func (r *Database) timeRange(startDate *time2.Time, endDate *time2.Time) func(db *gorm.DB) *gorm.DB {
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

func (r *Database) whereQuery(queryParameter WhereQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(queryParameter.Query, queryParameter.Parameters)
	}
}

func (r *Database) internalQuery(model interface{}, queryParameters QueryParameters) *gorm.DB {
	dbSub := r.db.Model(&model)
	if queryParameters.SelectQuery != nil {
		dbSub = dbSub.Select(*queryParameters.SelectQuery)
	}
	if queryParameters.Distinct {
		dbSub = dbSub.Distinct()
	}
	dbSub = dbSub.Scopes(r.timeRange(queryParameters.StartDate, queryParameters.EndDate))
	if queryParameters.JoinQuery != nil {
		dbSub = dbSub.Joins(*queryParameters.JoinQuery)
	}
	if queryParameters.WhereQuery != nil && len(*queryParameters.WhereQuery) > 0 {
		for _, item := range *queryParameters.WhereQuery {
			dbSub = dbSub.Scopes(r.whereQuery(item))
		}
	}
	if queryParameters.GroupBy != nil {
		dbSub = dbSub.Group(*queryParameters.GroupBy)
	}
	if queryParameters.OrderBy != nil {
		dbSub = dbSub.Order(*queryParameters.OrderBy)
	}

	return dbSub
}

func (r *Database) initCountryCodes() {
	test := schemas.CountryGeoLocation{}
	result := r.db.Model(schemas.CountryGeoLocation{}).First(&test)
	if result.RowsAffected == 0 {
		countryGeoLocationData.BuildCountryGeoLocationData()
		r.db.Model(schemas.CountryGeoLocation{}).CreateInBatches(countryGeoLocationData.CountryGeoLocationData, 50)
	}
}

func CreateGenericDatabase(debug bool) (*Database, error) {
	db := Database{}
	db.initDatabase("data", debug)

	return &db, nil
}
