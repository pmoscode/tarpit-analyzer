package database

import (
	geolocation "endlessh-analyzer/api/structs"
	"endlessh-analyzer/database/schemas"
	log "github.com/sirupsen/logrus"
	"strings"
)

type DbCache struct {
	database
}

func (r *DbCache) GetLocation(ip string) (schemas.Location, DbResult) {
	var location schemas.Location
	result := r.db.Where("ip = ?", ip).First(&location)

	if result.RowsAffected == 0 {
		return schemas.Location{}, DbRecordNotFound
	}

	return location, DbOk
}

func (r *DbCache) AddOrUpdateLocation(location schemas.Location) (DbResult, error) {
	var loc schemas.Location
	resultSelect := r.db.Where("ip = ?", location.Ip).First(&loc)

	if resultSelect.RowsAffected == 0 {
		resultInsert := r.db.Create(&location)
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
		loc.Latitude = location.Latitude
		loc.Longitude = location.Longitude
		loc.City = location.City
		loc.RegionName = location.RegionName
		loc.Region = location.Region
		loc.Continent = location.Continent
		loc.ContinentCode = location.ContinentCode
		loc.CountryCode = location.CountryCode
		loc.Country = location.Country
		loc.Ip = location.Ip

		resultSave := r.db.Save(&loc)
		if resultSave.Error != nil {
			return DbError, resultSave.Error
		}
	}

	return DbOk, nil
}

func (r *DbCache) DeleteLocation(location geolocation.GeoLocationItem) error {
	var loc schemas.Location
	resultSelect := r.db.Where("ip = ?", location.Ip).First(&loc)

	if resultSelect.Error != nil {
		return resultSelect.Error
	}

	if resultSelect.RowsAffected > 0 {
		r.db.Delete(loc)
	}

	return nil
}

func (r *DbCache) MapToLocation(location geolocation.GeoLocationItem) schemas.Location {
	return schemas.Location{
		GeoLocationItem: geolocation.GeoLocationItem{
			Ip:            location.Ip,
			Status:        location.Status,
			Continent:     location.Continent,
			ContinentCode: location.ContinentCode,
			Country:       location.Country,
			CountryCode:   location.CountryCode,
			Region:        location.Region,
			RegionName:    location.RegionName,
			City:          location.City,
			Zip:           location.Zip,
			Latitude:      location.Latitude,
			Longitude:     location.Longitude,
		},
	}
}

func (r *DbCache) MapToGeoLocation(location schemas.Location) geolocation.GeoLocationItem {
	return geolocation.GeoLocationItem{
		Status:        location.Status,
		Continent:     location.Continent,
		ContinentCode: location.ContinentCode,
		Country:       location.Country,
		CountryCode:   location.CountryCode,
		Region:        location.Region,
		RegionName:    location.RegionName,
		City:          location.City,
		Zip:           location.Zip,
		Latitude:      location.Latitude,
		Longitude:     location.Longitude,
		Ip:            location.Ip,
	}
}

func (r *DbCache) Map(vs []geolocation.GeoLocationItem, f func(location geolocation.GeoLocationItem) schemas.Location) []schemas.Location {
	vsm := make([]schemas.Location, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}

	return vsm
}

func CreateDbCache(debug bool) (DbCache, error) {
	db := DbCache{}
	db.initDatabase("data", debug)

	err := db.db.AutoMigrate(&schemas.Location{})
	if err != nil {
		log.Errorln(err)
		return DbCache{}, err
	}

	return db, nil
}
