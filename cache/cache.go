package cache

import (
	"endlessh-analyzer/api"
	geolocation "endlessh-analyzer/api/structs"
	"endlessh-analyzer/database"
	"errors"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	time2 "time"
)

var (
	db             database.DbCache
	geolocationApi api.Api
)

type Result int

const (
	Ok Result = iota
	RecordOutdated
	NoHit
)

func Init(api api.Api, debug bool) {
	geolocationApi = api
	dbInit, err := database.CreateDbCache(debug)
	if err != nil {
		log.Panicln("Cache database could not be loaded.", err)
	}

	db = dbInit
}

func GetLocationFor(ip string) *geolocation.GeoLocationItem {
	geoLocationItem, result := getLocation(ip)

	switch result {
	case Ok:
		return &geoLocationItem
	case RecordOutdated:
		err := db.DeleteLocation(geoLocationItem)
		if err != nil {
			return nil
		}
	}

	geolocationApi := api.CreateGeoLocationAPI(geolocationApi)
	ipBatch := []string{ip}
	resolved, errApi := geolocationApi.QueryGeoLocationAPI(&ipBatch)

	if errApi != nil {
		log.Errorln("Geolocation Api error: ", errApi)
		return nil
	} else {
		errDb := saveLocations(resolved)
		if errDb != nil {
			log.Errorln("Save Location to DB error: ", errDb)
		}
	}

	return &resolved[0]
}

func ResoleLocationsFor(ips []string, batchSize int) int {
	processedIps := 0
	ipsArraySize := len(ips)
	batch := make([]string, 0)
	cacheHits := 0
	bar := progressbar.Default(int64(ipsArraySize), "Checking IP's...")

	for i := 0; i < ipsArraySize; i++ {
		_, cacheResult := getLocation(ips[i])

		if cacheResult == NoHit || cacheResult == RecordOutdated {
			batch = append(batch, ips[i])
		} else if cacheResult == Ok {
			cacheHits++
		} else {
			log.Errorln("Something went wrong for ip: ", ips[i])
		}

		bar.Add(1)
	}
	bar.Finish()

	geolocationApi := api.CreateGeoLocationAPI(geolocationApi)
	batchCount := len(batch)
	if batchCount > 0 {
		bar = progressbar.Default(int64(batchCount), "Processing IP's...")

		for i := 0; i < batchCount; i += batchSize {
			if i+batchSize >= batchCount {
				batchSize = batchCount - i
			}

			ipBatch := batch[i : i+batchSize]
			resolved, errApi := geolocationApi.QueryGeoLocationAPI(&ipBatch)

			processedIps = processedIps + len(resolved)
			bar.Add(len(resolved))

			if errApi != nil {
				log.Errorln("Geolocation Api error: ", errApi)
			} else {
				errDb := saveLocations(resolved)
				if errDb != nil {
					log.Errorln("Save Location to DB error: ", errDb)
				}
			}
			err := saveLocations(resolved)
			if err != nil {
				log.Warningln("Could not save location: ", resolved)
			}
		}
	}
	bar.Finish()
	log.Debugln("Cache Hits: ", cacheHits)

	return processedIps + cacheHits
}

func getLocation(ip string) (geolocation.GeoLocationItem, Result) {
	location, dbResult := db.GetLocation(ip)

	if dbResult == database.DbOk {
		geoLocation := db.MapToGeoLocation(location)

		var maxAge float64 = 4 * 24
		if time2.Now().Sub(location.UpdatedAt).Hours() > maxAge {
			log.Infoln("Record for ip: ", ip, " is older than ", maxAge, " hours... ")
			return geoLocation, RecordOutdated
		}

		return geoLocation, Ok
	}

	return geolocation.GeoLocationItem{}, NoHit
}

func saveLocations(locations []geolocation.GeoLocationItem) error {
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
