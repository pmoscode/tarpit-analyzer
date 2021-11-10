package cache_db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

type Location struct {
	gorm.Model
	Ip          string `gorm:"uniqueIndex"`
	Status      string
	Country     string
	CountryCode string
	Region      string
	RegionName  string
	City        string
	Zip         string
	Lat         float64
	Lon         float64
	Timezone    string
}

type DbResult int

const (
	DbOk DbResult = iota
	DbRecordNotFound
	DbError
)

var db *gorm.DB

func initDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	dbT, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db = dbT

	// Migrate the schema
	err = db.AutoMigrate(&Location{})
	if err != nil {
		return
	}
}

func getLocation(ip string) (Location, DbResult) {
	var location Location
	result := db.Where("ip = ?", ip).First(&location)

	if result.RowsAffected == 0 {
		return Location{}, DbRecordNotFound
	}

	return location, DbOk
}

func addOrUpdateLocation(location Location) (DbResult, error) {
	var loc Location
	resultSelect := db.Where("ip = ?", location.Ip).First(&loc)

	if resultSelect.RowsAffected == 0 {
		resultInsert := db.Create(&location)
		if resultInsert.Error != nil {
			if strings.Contains(resultInsert.Error.Error(), "UNIQUE constraint failed") { // Theoretically should not happen
				return DbOk, nil
			} else {
				return DbError, resultInsert.Error
			}
		}
	} else {
		loc.Status = location.Status
		loc.Zip = location.Zip
		loc.Lat = location.Lat
		loc.Lon = location.Lon
		loc.City = location.City
		loc.Timezone = location.Timezone
		loc.RegionName = location.RegionName
		loc.Region = location.Region
		loc.CountryCode = location.CountryCode
		loc.Country = location.Country
		loc.Ip = location.Ip

		db.Save(&loc)

		return DbOk, nil
	}

	return DbError, nil
}
