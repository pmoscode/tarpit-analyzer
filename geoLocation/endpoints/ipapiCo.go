package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tarpit-analyzer/geoLocation/structs"
	"time"
)

type IpapiCo struct {
	lastExecutionFinished time.Time
}

type IpapiCoItem struct {
	Ip                 string  `json:"ip"`
	Version            string  `json:"version"`
	City               string  `json:"city"`
	Region             string  `json:"region"`
	RegionCode         string  `json:"region_code"`
	CountryCode        string  `json:"country_code"`
	CountryCodeIso3    string  `json:"country_code_iso3"`
	CountryName        string  `json:"country_name"`
	CountryCapital     string  `json:"country_capital"`
	CountryTld         string  `json:"country_tld"`
	ContinentCode      string  `json:"continent_code"`
	InEu               bool    `json:"in_eu"`
	Postal             string  `json:"postal"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Timezone           string  `json:"timezone"`
	UtcOffset          string  `json:"utc_offset"`
	CountryCallingCode string  `json:"country_calling_code"`
	Currency           string  `json:"currency"`
	CurrencyName       string  `json:"currency_name"`
	Languages          string  `json:"languages"`
	CountryArea        float64 `json:"country_area"`
	CountryPopulation  int     `json:"country_population"`
	Asn                string  `json:"asn"`
	Org                string  `json:"org"`
}

func (r *IpapiCo) QueryGeoLocationAPI(ips *[]string, bar *progressbar.ProgressBar) ([]structs.GeoLocationItem, error) {
	maxRequests := 999 // per 24 hours and max 30000 per month
	mappedLocations := make([]structs.GeoLocationItem, 0)

	for _, ip := range *ips {
		resp, err := http.Get("https://ipapi.co/" + ip + "/json/")
		if err != nil {
			log.Debugln("No response from request")
			r.lastExecutionFinished = time.Now()

			return nil, err
		}

		if resp.StatusCode == 200 {
			ipLocation := IpapiCoItem{}
			err = json.NewDecoder(resp.Body).Decode(&ipLocation)
			if err != nil {
				log.Debugln(err)
				r.lastExecutionFinished = time.Now()

				return nil, err
			}

			mappedLocation, err := r.mapToGeoLocationItem(&ipLocation)
			if err != nil {
				_ = resp.Body.Close()
				r.lastExecutionFinished = time.Now()

				return nil, err
			}

			_ = bar.Add(1)
			mappedLocations = append(mappedLocations, mappedLocation)
		} else {
			_ = resp.Body.Close()
			log.Debugln("Done requests: ", 200-maxRequests)
			r.lastExecutionFinished = time.Now()

			return nil, errors.New("got response from api: " + resp.Status)
		}

		_ = resp.Body.Close()
		maxRequests--
		if maxRequests == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	r.lastExecutionFinished = time.Now()

	return mappedLocations, nil
}

func (r *IpapiCo) Name() string {
	return "IpapiCo"
}

func (r *IpapiCo) mapToGeoLocationItem(item *IpapiCoItem) (structs.GeoLocationItem, error) {
	return structs.GeoLocationItem{
		Ip:            item.Ip,
		Status:        "success",
		Latitude:      item.Latitude,
		Longitude:     item.Longitude,
		Continent:     "not available",
		ContinentCode: item.ContinentCode,
		Country:       item.CountryName,
		CountryCode:   item.CountryCode,
		Region:        item.RegionCode,
		RegionName:    item.Region,
		City:          item.City,
		Zip:           item.Postal,
	}, nil
}

func (r *IpapiCo) CanHandle() bool {
	return time.Now().Sub(r.lastExecutionFinished).Hours() > 24
}
