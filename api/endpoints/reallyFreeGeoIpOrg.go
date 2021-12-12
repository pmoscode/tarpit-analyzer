package endpoints

import (
	"encoding/json"
	"endlessh-analyzer/api/structs"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ReallyFreeGeoIpOrg struct{}

type ReallyFreeGeoIpOrgItem struct {
	Ip          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
}

func (r ReallyFreeGeoIpOrg) QueryGeoLocationAPI(ips *[]string) ([]structs.GeoLocationItem, error) {
	mappedLocations := make([]structs.GeoLocationItem, len(*ips))

	for idx, ip := range *ips {
		resp, err := http.Get("https://reallyfreegeoip.org/json/" + ip)
		if err != nil {
			log.Warningln("No response from request")
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			ipLocation := ReallyFreeGeoIpOrgItem{}
			err = json.NewDecoder(resp.Body).Decode(&ipLocation)
			if err != nil {
				log.Errorln(err)
			}

			mappedLocation, err := r.mapToGeoLocationItem(&ipLocation)
			if err != nil {
				return nil, err
			}

			mappedLocations[idx] = mappedLocation
		} else {
			return nil, errors.New("got response from api: " + resp.Status)
		}
	}

	return mappedLocations, nil
}

func (r ReallyFreeGeoIpOrg) mapToGeoLocationItem(item *ReallyFreeGeoIpOrgItem) (structs.GeoLocationItem, error) {
	return structs.GeoLocationItem{
		Ip:            item.Ip,
		Status:        "success",
		Latitude:      item.Latitude,
		Longitude:     item.Longitude,
		Continent:     "not available",
		ContinentCode: "not available",
		Country:       item.CountryName,
		CountryCode:   item.CountryCode,
		Region:        item.RegionCode,
		RegionName:    item.RegionName,
		City:          item.City,
		Zip:           item.ZipCode,
	}, nil
}
