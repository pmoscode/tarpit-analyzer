package endpoints

import (
	"bytes"
	"encoding/json"
	"endlessh-analyzer/api/structs"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type IpApiCom struct{}

type IpApiComItem struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func (c IpApiCom) QueryGeoLocationApi(ips []string) ([]structs.GeoLocationItem, error) {
	body := "[\"" + strings.Join(ips, "\",\"") + "\"]"

	resp, err := http.Post("http://ip-api.com/batch", "application/json", bytes.NewBufferString(body))
	if err != nil {
		log.Warningln("No response from request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		ipLocation := make([]IpApiComItem, 0)
		err = json.NewDecoder(resp.Body).Decode(&ipLocation)
		if err != nil {
			log.Errorln(err)
		}

		mappedLocations, err := mapToGeoLocationItem(&ipLocation)
		if err != nil {
			return nil, err
		}

		return mappedLocations, nil
	}

	return nil, errors.New("got response from api: " + resp.Status)
}

func mapToGeoLocationItem(items *[]IpApiComItem) ([]structs.GeoLocationItem, error) {
	count := len(*items)
	locations := make([]structs.GeoLocationItem, count)

	if count != 0 {
		for idx, item := range *items {
			target := structs.GeoLocationItem{
				Ip:          item.Query,
				Status:      item.Status,
				Latitude:    item.Lat,
				Longitude:   item.Lon,
				Country:     item.Country,
				CountryCode: item.CountryCode,
				Region:      item.Region,
				RegionName:  item.RegionName,
				City:        item.City,
				Zip:         item.Zip,
			}
			locations[idx] = target
		}
	}

	return locations, nil
}
