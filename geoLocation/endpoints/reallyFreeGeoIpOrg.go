package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tarpit-analyzer/geoLocation/structs"
	time2 "time"
)

type ReallyFreeGeoIpOrg struct {
	lastExecutionFinished time2.Time
}

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

func (r *ReallyFreeGeoIpOrg) QueryGeoLocationAPI(ips *[]string, bar *progressbar.ProgressBar) ([]structs.GeoLocationItem, error) {
	mappedLocations := make([]structs.GeoLocationItem, 0)
	maxRequests := 1500

	for _, ip := range *ips {
		resp, err := http.Get("https://reallyfreegeoip.org/json/" + ip)
		if err != nil {
			log.Debugln("No response from request: ")
			_ = resp.Body.Close()
			r.lastExecutionFinished = time2.Now()

			return nil, err
		}

		if resp.StatusCode == 200 {
			ipLocation := ReallyFreeGeoIpOrgItem{}
			err = json.NewDecoder(resp.Body).Decode(&ipLocation)
			if err != nil {
				log.Debugln(err)
				r.lastExecutionFinished = time2.Now()

				return nil, err
			}

			mappedLocation, err := r.mapToGeoLocationItem(&ipLocation)
			if err != nil {
				_ = resp.Body.Close()
				r.lastExecutionFinished = time2.Now()

				return nil, err
			}

			_ = bar.Add(1)
			mappedLocations = append(mappedLocations, mappedLocation)
		} else {
			_ = resp.Body.Close()
			r.lastExecutionFinished = time2.Now()

			return nil, errors.New("got response from api: " + resp.Status)
		}

		_ = resp.Body.Close()
		maxRequests--
		if maxRequests == 0 {
			break
		}
	}

	r.lastExecutionFinished = time2.Now()

	return mappedLocations, nil
}

func (r *ReallyFreeGeoIpOrg) Name() string {
	return "ReallyFreeGeoIpOrg"
}

func (r *ReallyFreeGeoIpOrg) mapToGeoLocationItem(item *ReallyFreeGeoIpOrgItem) (structs.GeoLocationItem, error) {
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

func (r *ReallyFreeGeoIpOrg) CanHandle() bool {
	return time2.Now().Sub(r.lastExecutionFinished).Minutes() > 1
}
