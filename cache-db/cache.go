package cache_db

import (
	geolocation "endlessh-analyzer/api/structs"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	time2 "time"
)

type CacheResult int

const (
	CacheOk CacheResult = iota
	CacheRecordOutdated
	CacheNoHit
)

func Init() {
	initDatabase()
}

func GetLocationFor(ip string) (geolocation.GeoLocationItem, CacheResult) {
	location, dbResult := getLocation(ip)
	if dbResult == DbRecordNotFound {
		log.Infoln("Record for ip: ", ip, " not found...")
		return geolocation.GeoLocationItem{}, CacheNoHit
	}

	var maxAge float64 = 4 * 24
	if time2.Now().Sub(location.UpdatedAt).Hours() > maxAge {
		log.Infoln("Record for ip: ", ip, " is older than ", maxAge, " hours... Query again from Geo-Location API...")
		return mapToGeoLocation(location), CacheRecordOutdated
	}

	return mapToGeoLocation(location), CacheOk
}

func SaveLocations(locations []geolocation.GeoLocationItem) error {
	locs := Map(locations, mapToLocation)

	for i := 0; i < len(locs); i++ {
		dbResult, err := addOrUpdateLocation(locs[i])
		if err != nil {
			return err
		}

		fmt.Print(i, " of ", len(locs), "\r")

		if dbResult != DbOk {
			log.Errorln("something went wrong for location: ", locs[i])
			return errors.New("something went wrong for location IP: " + locs[i].Ip)
		}
	}

	return nil
}

func mapToLocation(location geolocation.GeoLocationItem) Location {
	return Location{
		GeoLocationItem: geolocation.GeoLocationItem{
			Ip:          location.Ip,
			Status:      location.Status,
			Country:     location.Country,
			CountryCode: location.CountryCode,
			Region:      location.Region,
			RegionName:  location.RegionName,
			City:        location.City,
			Zip:         location.Zip,
			Latitude:    location.Latitude,
			Longitude:   location.Longitude,
		},
	}
}

func mapToGeoLocation(location Location) geolocation.GeoLocationItem {
	return geolocation.GeoLocationItem{
		Status:      location.Status,
		Country:     location.Country,
		CountryCode: location.CountryCode,
		Region:      location.Region,
		RegionName:  location.RegionName,
		City:        location.City,
		Zip:         location.Zip,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		Ip:          location.Ip,
	}
}

func Map(vs []geolocation.GeoLocationItem, f func(location geolocation.GeoLocationItem) Location) []Location {
	vsm := make([]Location, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}

	return vsm
}
