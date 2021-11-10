package api

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type IpLocation struct {
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

func DoQuery(ips []string) ([]IpLocation, error) {

	body := "[\"" + strings.Join(ips, "\",\"") + "\"]"

	resp, err := http.Post("http://ip-api.com/batch", "application/json", bytes.NewBufferString(body))
	if err != nil {
		log.Warningln("No response from request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		ipLocation := make([]IpLocation, 0)
		err = json.NewDecoder(resp.Body).Decode(&ipLocation)
		if err != nil {
			log.Errorln(err)
		}

		return ipLocation, nil
	}

	return nil, errors.New(resp.Status)
}
