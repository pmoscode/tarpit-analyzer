package database

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type DbResult int

const (
	DbOk DbResult = iota
	DbRecordNotFound
	DbError
)

type database struct {
	db *gorm.DB
}

func (r *database) initDatabase(dbFilename string, debug bool) {

	logLevel := logger.Silent
	if debug {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
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

	r.db = db
}

func (r *database) DbRawQuery(model interface{}, query string, parameters ...interface{}) (*sql.Rows, error) {
	rows, err := r.db.Model(model).Raw(query, parameters).Rows()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
