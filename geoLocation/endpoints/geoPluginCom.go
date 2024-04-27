package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tarpit-analyzer/geoLocation/structs"
	"time"
)

type GeoPluginCom struct {
	lastExecutionFinished time.Time
}

type GeoPluginComItem struct {
	GeopluginRequest                string  `json:"geoplugin_request"`
	GeopluginStatus                 int     `json:"geoplugin_status"`
	GeopluginDelay                  string  `json:"geoplugin_delay"`
	GeopluginCredit                 string  `json:"geoplugin_credit"`
	GeopluginCity                   string  `json:"geoplugin_city"`
	GeopluginRegion                 string  `json:"geoplugin_region"`
	GeopluginRegionCode             string  `json:"geoplugin_regionCode"`
	GeopluginRegionName             string  `json:"geoplugin_regionName"`
	GeopluginAreaCode               string  `json:"geoplugin_areaCode"`
	GeopluginDmaCode                string  `json:"geoplugin_dmaCode"`
	GeopluginCountryCode            string  `json:"geoplugin_countryCode"`
	GeopluginCountryName            string  `json:"geoplugin_countryName"`
	GeopluginInEU                   int     `json:"geoplugin_inEU"`
	GeopluginEuVATrate              bool    `json:"geoplugin_euVATrate"`
	GeopluginContinentCode          string  `json:"geoplugin_continentCode"`
	GeopluginContinentName          string  `json:"geoplugin_continentName"`
	GeopluginLatitude               string  `json:"geoplugin_latitude"`
	GeopluginLongitude              string  `json:"geoplugin_longitude"`
	GeopluginLocationAccuracyRadius string  `json:"geoplugin_locationAccuracyRadius"`
	GeopluginTimezone               string  `json:"geoplugin_timezone"`
	GeopluginCurrencyCode           string  `json:"geoplugin_currencyCode"`
	GeopluginCurrencySymbol         string  `json:"geoplugin_currencySymbol"`
	GeopluginCurrencySymbolUTF8     string  `json:"geoplugin_currencySymbol_UTF8"`
	GeopluginCurrencyConverter      float64 `json:"geoplugin_currencyConverter"`
}

func (r *GeoPluginCom) QueryGeoLocationAPI(ips *[]string, bar *progressbar.ProgressBar) ([]structs.GeoLocationItem, error) {
	maxRequests := 500
	mappedLocations := make([]structs.GeoLocationItem, 0)

	for _, ip := range *ips {
		resp, err := http.Get("http://www.geoplugin.net/json.gp?ip=" + ip)
		if err != nil {
			log.Debugln("No response from request")
			r.lastExecutionFinished = time.Now()

			return nil, err
		}

		if resp.StatusCode == 200 {
			ipLocation := GeoPluginComItem{}
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

func (r *GeoPluginCom) Name() string {
	return "GeoPluginCom"
}

func (r *GeoPluginCom) mapToGeoLocationItem(item *GeoPluginComItem) (structs.GeoLocationItem, error) {
	lat, _ := strconv.ParseFloat(item.GeopluginLatitude, 64)
	lon, _ := strconv.ParseFloat(item.GeopluginLongitude, 64)

	return structs.GeoLocationItem{
		Ip:            item.GeopluginRequest,
		Status:        "success",
		Latitude:      lat,
		Longitude:     lon,
		Continent:     item.GeopluginContinentCode,
		ContinentCode: item.GeopluginContinentName,
		Country:       item.GeopluginCountryName,
		CountryCode:   item.GeopluginCountryCode,
		Region:        item.GeopluginRegionCode,
		RegionName:    item.GeopluginRegionName,
		City:          item.GeopluginCity,
		Zip:           "-1",
	}, nil
}

func (r *GeoPluginCom) CanHandle() bool {
	return time.Now().Sub(r.lastExecutionFinished).Minutes() > 60
}
