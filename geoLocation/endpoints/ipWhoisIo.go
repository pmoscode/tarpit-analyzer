package endpoints

import (
	"encoding/json"
	"endlessh-analyzer/geoLocation/structs"
	"errors"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type IpWhoIsIo struct{}

type IpWhoIsIoItem struct {
	Ip                string  `json:"ip"`
	Success           bool    `json:"success"`
	Type              string  `json:"type"`
	Continent         string  `json:"continent"`
	ContinentCode     string  `json:"continent_code"`
	Country           string  `json:"country"`
	CountryCode       string  `json:"country_code"`
	CountryFlag       string  `json:"country_flag"`
	CountryCapital    string  `json:"country_capital"`
	CountryPhone      string  `json:"country_phone"`
	CountryNeighbours string  `json:"country_neighbours"`
	Region            string  `json:"region"`
	City              string  `json:"city"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Asn               string  `json:"asn"`
	Org               string  `json:"org"`
	Isp               string  `json:"isp"`
	Timezone          string  `json:"timezone"`
	TimezoneName      string  `json:"timezone_name"`
	TimezoneDstOffset int     `json:"timezone_dstOffset"`
	TimezoneGmtOffset int     `json:"timezone_gmtOffset"`
	TimezoneGmt       string  `json:"timezone_gmt"`
	Currency          string  `json:"currency"`
	CurrencyCode      string  `json:"currency_code"`
	CurrencySymbol    string  `json:"currency_symbol"`
	CurrencyRates     int     `json:"currency_rates"`
	CurrencyPlural    string  `json:"currency_plural"`
	CompletedRequests int     `json:"completed_requests"`
}

func (r IpWhoIsIo) QueryGeoLocationAPI(ips *[]string, bar *progressbar.ProgressBar) ([]structs.GeoLocationItem, error) {
	mappedLocations := make([]structs.GeoLocationItem, 0)
	maxRequests := 10000 // per month

	for _, ip := range *ips {
		resp, err := http.Get("http://ipwhois.app/json/" + ip)
		if err != nil {
			log.Debugln("No response from request: ")
			_ = resp.Body.Close()

			return nil, err
		}

		if resp.StatusCode == 200 {
			ipLocation := IpWhoIsIoItem{}
			err = json.NewDecoder(resp.Body).Decode(&ipLocation)
			if err != nil {
				log.Debugln(err)
				return nil, err
			}

			mappedLocation, err := r.mapToGeoLocationItem(&ipLocation)
			if err != nil {
				_ = resp.Body.Close()

				return nil, err
			}

			_ = bar.Add(1)
			mappedLocations = append(mappedLocations, mappedLocation)
		} else {
			_ = resp.Body.Close()

			return nil, errors.New("got response from api: " + resp.Status)
		}

		_ = resp.Body.Close()
		maxRequests--
		if maxRequests == 0 {
			break
		}
	}

	return mappedLocations, nil
}

func (r IpWhoIsIo) Name() string {
	return "IpWhoIsIo"
}

func (r *IpWhoIsIo) mapToGeoLocationItem(item *IpWhoIsIoItem) (structs.GeoLocationItem, error) {
	return structs.GeoLocationItem{
		Ip:            item.Ip,
		Status:        strconv.FormatBool(item.Success),
		Latitude:      item.Latitude,
		Longitude:     item.Longitude,
		Continent:     item.Continent,
		ContinentCode: item.ContinentCode,
		Country:       item.Country,
		CountryCode:   item.CountryCode,
		Region:        "not available",
		RegionName:    item.Region,
		City:          item.City,
		Zip:           "-1",
	}, nil
}
