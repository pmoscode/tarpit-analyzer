package geoLocation

import (
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"tarpit-analyzer/geoLocation/endpoints"
	"tarpit-analyzer/geoLocation/structs"
)

type QueryGeoLocationAPI interface {
	QueryGeoLocationAPI(ips *[]string, bar *progressbar.ProgressBar) ([]structs.GeoLocationItem, error)
	Name() string
	CanHandle() bool
}

type GeoLocation struct {
	apis []QueryGeoLocationAPI
}

func (r *GeoLocation) ResolveLocations(ips []string) (*[]structs.GeoLocationItem, error) {
	resolvedGeoLocationItems := make([]structs.GeoLocationItem, 0)

	bar := progressbar.Default(int64(len(ips)), "Geo location of IP's...")

	for {
		geoApi := r.nextApi()
		if geoApi == nil {
			log.Warningln("No more endpoint found which can resolve IP's... Try it in a few minutes, hours or days...")
			break
		}

		geoLocationItems, err := geoApi.QueryGeoLocationAPI(&ips, bar)
		if err == nil {
			resolvedItemsCount := len(geoLocationItems)
			ips = ips[resolvedItemsCount:]
			resolvedGeoLocationItems = append(resolvedGeoLocationItems, geoLocationItems...)
			if len(ips) == 0 {
				break
			}
		} else {
			log.Errorln("Api '", geoApi.Name(), "' got error: ", err)
		}
	}

	_ = bar.Finish()

	return &resolvedGeoLocationItems, nil
}

func (r *GeoLocation) nextApi() QueryGeoLocationAPI {
	for _, api := range r.apis {
		if api.CanHandle() {
			return api
		}
	}

	return nil
}

func CreateGeoLocation() *GeoLocation {
	return &GeoLocation{apis: []QueryGeoLocationAPI{&endpoints.IpApiCom{}, &endpoints.ReallyFreeGeoIpOrg{},
		&endpoints.IpWhoIsIo{}, &endpoints.IpapiCo{}, &endpoints.GeoPluginCom{}}}
}
