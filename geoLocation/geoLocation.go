package geoLocation

import (
	"endlessh-analyzer/geoLocation/endpoints"
	"endlessh-analyzer/geoLocation/structs"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type QueryGeoLocationAPI interface {
	QueryGeoLocationAPI(ips *[]string, bar *progressbar.ProgressBar) ([]structs.GeoLocationItem, error)
	Name() string
}

type GeoLocation struct {
	apis []QueryGeoLocationAPI
}

func (r *GeoLocation) ResolveLocations(ips []string) (*[]structs.GeoLocationItem, error) {
	resolvedGeoLocationItems := make([]structs.GeoLocationItem, 0)

	bar := progressbar.Default(int64(len(ips)), "Geo location of IP's...")
	for _, api := range r.apis {
		geoLocationItems, err := api.QueryGeoLocationAPI(&ips, bar)
		if err == nil {
			resolvedItemsCount := len(geoLocationItems)
			ips = ips[resolvedItemsCount:]
			resolvedGeoLocationItems = append(resolvedGeoLocationItems, geoLocationItems...)
			if len(ips) == 0 {
				break
			}
		} else {
			log.Errorln("Api '", api.Name(), "' got error: ", err)
		}
	}
	_ = bar.Finish()

	return &resolvedGeoLocationItems, nil
}

func CreateGeoLocation() *GeoLocation {
	return &GeoLocation{apis: []QueryGeoLocationAPI{&endpoints.IpApiCom{}, &endpoints.ReallyFreeGeoIpOrg{}, &endpoints.IpapiCo{}, &endpoints.GeoPluginCom{}}}
}
