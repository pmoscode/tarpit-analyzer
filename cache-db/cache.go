package cache_db

import (
	geolocation "endlessh-analyzer/api"
	"errors"
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

func GetLocationFor(ip string) (geolocation.IpLocation, CacheResult) {
	location, dbResult := getLocation(ip)
	if dbResult == DbRecordNotFound {
		log.Infoln("Record for ip: ", ip, " not found...")
		return geolocation.IpLocation{}, CacheNoHit
	}

	var maxAge float64 = 4 * 24
	if time2.Now().Sub(location.UpdatedAt).Hours() > maxAge {
		log.Infoln("Record for ip: ", ip, " is older than ", maxAge, " hours... Query again from Geo-Location API...")
		return mapToIpLocation(location), CacheRecordOutdated
	}

	return mapToIpLocation(location), CacheOk
}

func SaveLocations(locations []geolocation.IpLocation) error {
	locs := Map(locations, mapToLocation)

	for i := 0; i < len(locs); i++ {
		dbResult, err := addOrUpdateLocation(locs[i])
		if err != nil {
			return err
		}

		if dbResult != DbOk {
			log.Errorln("something went wrong for location: ", locs[i])
			return errors.New("something went wrong for location IP: " + locs[i].Ip)
		}
	}

	return nil
}

func mapToLocation(location geolocation.IpLocation) Location {
	return Location{
		Ip:          location.Query,
		Status:      location.Status,
		Country:     location.Country,
		CountryCode: location.CountryCode,
		Region:      location.Region,
		RegionName:  location.RegionName,
		City:        location.City,
		Zip:         location.Zip,
		Lat:         location.Lat,
		Lon:         location.Lon,
		Timezone:    location.Timezone,
	}
}

func mapToIpLocation(location Location) geolocation.IpLocation {
	return geolocation.IpLocation{
		Status:      location.Status,
		Country:     location.Country,
		CountryCode: location.CountryCode,
		Region:      location.Region,
		RegionName:  location.RegionName,
		City:        location.City,
		Zip:         location.Zip,
		Lat:         location.Lat,
		Lon:         location.Lon,
		Timezone:    location.Timezone,
		Query:       location.Ip,
	}
}

func Map(vs []geolocation.IpLocation, f func(location geolocation.IpLocation) Location) []Location {
	vsm := make([]Location, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}

	return vsm
}
