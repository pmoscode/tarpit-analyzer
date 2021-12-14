package endpoints

import (
	"bytes"
	"encoding/json"
	"endlessh-analyzer/geoLocation/structs"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type IpApiCom struct{}

type IpApiComItem struct {
	Status        string  `json:"status"`
	Country       string  `json:"country"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Query         string  `json:"query"`
}

func (r IpApiCom) QueryGeoLocationAPI(ips *[]string) ([]structs.GeoLocationItem, error) {
	batchSize := 100
	maxRequests := 15
	batchCount := len(*ips)
	data := *ips
	mappedLocations := make([]structs.GeoLocationItem, 0)

	if batchCount > 0 {
		for i := 0; i < batchCount; i += batchSize {
			if i+batchSize >= batchCount {
				batchSize = batchCount - i
			}

			ipBatch := data[i : i+batchSize]

			body := "[\"" + strings.Join(ipBatch, "\",\"") + "\"]"

			resp, errRequest := http.Post("http://ip-api.com/batch?fields=status,continent,continentCode,country,countryCode,region,regionName,city,zip,lat,lon,query", "application/json", bytes.NewBufferString(body))
			if errRequest != nil {
				log.Warningln("No response from request")
				return nil, errRequest
			}

			if resp.StatusCode == 429 {
				log.Debugln("Max requests (15) per minute reached!")
				_ = resp.Body.Close()
				return nil, errors.New("max requests reached")
			}

			if resp.StatusCode != 200 {
				_ = resp.Body.Close()
				return nil, errors.New("got response from api: " + resp.Status)
			}

			ipLocation := make([]IpApiComItem, 0)
			errJson := json.NewDecoder(resp.Body).Decode(&ipLocation)
			if errJson != nil {
				log.Errorln(errJson)
				return nil, errJson
			}

			mappedLocationsLocal, errMap := r.mapToGeoLocationItem(&ipLocation)
			if errMap != nil {
				return nil, errMap
			}

			_ = resp.Body.Close()
			mappedLocations = append(mappedLocations, mappedLocationsLocal...)

			maxRequests--
			if maxRequests == 0 {
				break
			}
		}
	}

	return mappedLocations, nil
}

func (r IpApiCom) Name() string {
	return "IpApiCom"
}

func (r *IpApiCom) mapToGeoLocationItem(items *[]IpApiComItem) ([]structs.GeoLocationItem, error) {
	count := len(*items)
	locations := make([]structs.GeoLocationItem, count)

	if count != 0 {
		for idx, item := range *items {
			target := structs.GeoLocationItem{
				Ip:            item.Query,
				Status:        item.Status,
				Latitude:      item.Lat,
				Longitude:     item.Lon,
				Continent:     item.Continent,
				ContinentCode: item.ContinentCode,
				Country:       item.Country,
				CountryCode:   item.CountryCode,
				Region:        item.Region,
				RegionName:    item.RegionName,
				City:          item.City,
				Zip:           item.Zip,
			}
			locations[idx] = target
		}
	}

	return locations, nil
}
