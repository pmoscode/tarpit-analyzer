package cache

import (
	"endlessh-analyzer/database"
	"endlessh-analyzer/geoLocation"
	geolocationStructs "endlessh-analyzer/geoLocation/structs"
	"errors"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	time2 "time"
)

var (
	db database.DbCache
)

type Result int

const (
	Ok Result = iota
	RecordOutdated
	NoHit
)

func Init(debug bool) {
	dbInit, err := database.CreateDbCache(debug)
	if err != nil {
		log.Panicln("Cache database could not be loaded.", err)
	}

	db = dbInit
}

func GetLocationFor(ip string) *geolocationStructs.GeoLocationItem {
	ResoleLocationsFor([]string{ip})
	geoLocationItem, _ := getSavedLocation(ip)

	return &geoLocationItem
}

func ResoleLocationsFor(ips []string) int {
	ipsArraySize := len(ips)
	batch := make([]string, 0)
	cacheHits := 0
	processedIps := 0

	bar := progressbar.Default(int64(ipsArraySize), "Checking IP's...")

	for i := 0; i < ipsArraySize; i++ {
		geoLocationItem, cacheResult := getSavedLocation(ips[i])

		switch cacheResult {
		case Ok:
			cacheHits++
		case RecordOutdated:
			err := db.DeleteLocation(geoLocationItem)
			if err != nil {
				log.Errorln("Could not delete outdated geoLocationItem: ", geoLocationItem)
			}
			fallthrough
		case NoHit:
			batch = append(batch, ips[i])
		}

		_ = bar.Add(1)
	}

	_ = bar.Finish()

	if len(batch) > 0 {
		geolocation := geoLocation.CreateGeoLocation()

		locations, err := geolocation.ResolveLocations(batch)
		if err != nil {
			return 0
		}
		processedIps = len(*locations)
		errDb := saveLocations(locations)
		if errDb != nil {
			log.Errorln("Save Location to DB error: ", errDb)
		}
	}

	log.Debugln("Cache Hits: ", cacheHits)

	return processedIps + cacheHits
}

func getSavedLocation(ip string) (geolocationStructs.GeoLocationItem, Result) {
	location, dbResult := db.GetLocation(ip)

	if dbResult == database.DbOk {
		geoLocationItem := db.MapToGeoLocation(location)

		var maxAge float64 = 4 * 24
		if time2.Now().Sub(location.UpdatedAt).Hours() > maxAge {
			log.Infoln("Record for ip: ", ip, " is older than ", maxAge, " hours... ")
			return geoLocationItem, RecordOutdated
		}

		return geoLocationItem, Ok
	}

	return geolocationStructs.GeoLocationItem{}, NoHit
}

func saveLocations(locations *[]geolocationStructs.GeoLocationItem) error {
	locs := db.Map(locations, db.MapToLocation)

	bar := progressbar.Default(int64(len(*locations)), "Saving resolved IP's...")
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
			_ = bar.Add(1)
		} else {
			log.Warningln("Geo Location info for: ", loc.Ip, " failed: ", loc)
		}
	}
	_ = bar.Finish()

	return nil
}
