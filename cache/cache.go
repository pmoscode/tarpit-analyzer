package cache

import (
	geolocation "endlessh-analyzer/api/structs"
	"endlessh-analyzer/database"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	time2 "time"
)

type Result int

const (
	Ok Result = iota
	RecordOutdated
	NoHit
)

var db database.DbCache

func Init() {
	dbInit, err := database.CreateDbCache()
	if err != nil {
		log.Panicln("Cache database could not be loaded.", err)
	}

	db = dbInit
}

func GetLocationFor(ip string) (geolocation.GeoLocationItem, Result) {
	location, dbResult := db.GetLocation(ip)
	if dbResult == database.DbRecordNotFound {
		log.Infoln("Record for ip: ", ip, " not found...")
		return geolocation.GeoLocationItem{}, NoHit
	}

	var maxAge float64 = 4 * 24
	if time2.Now().Sub(location.UpdatedAt).Hours() > maxAge {
		log.Infoln("Record for ip: ", ip, " is older than ", maxAge, " hours... Query again from Geo-Location API...")
		return db.MapToGeoLocation(location), RecordOutdated
	}

	return db.MapToGeoLocation(location), Ok
}

func SaveLocations(locations []geolocation.GeoLocationItem) error {
	locs := db.Map(locations, db.MapToLocation)

	for i := 0; i < len(locs); i++ {
		dbResult, err := db.AddOrUpdateLocation(locs[i])
		if err != nil {
			return err
		}

		fmt.Print(i, " of ", len(locs), "\r")

		if dbResult != database.DbOk {
			log.Errorln("something went wrong for location: ", locs[i])
			return errors.New("something went wrong for location IP: " + locs[i].Ip)
		}
	}

	return nil
}