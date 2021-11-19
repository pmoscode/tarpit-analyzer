package cache

import (
	geolocation "endlessh-analyzer/api/structs"
	"endlessh-analyzer/database"
	"errors"
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

func Init(debug bool) {
	dbInit, err := database.CreateDbCache(debug)
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

	for _, loc := range locs {
		if loc.Status == "success" {
			dbResult, err := db.AddOrUpdateLocation(loc)
			if err != nil {
				return err
			}

			if dbResult != database.DbOk {
				log.Errorln("something went wrong for location: ", loc)
				return errors.New("something went wrong for location IP: " + loc.Ip)
			}
		} else {
			log.Warningln("Geo Location info for: ", loc.Ip, " failed: ", loc)
		}
	}

	return nil
}
