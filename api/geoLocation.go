package api

import (
	"endlessh-analyzer/api/endpoints"
	"endlessh-analyzer/api/structs"
)

type QueryGeoLocationAPI interface {
	QueryGeoLocationAPI(ips *[]string) ([]structs.GeoLocationItem, error)
}

type Api int

const (
	IpApiCom Api = iota
	ReallyFreeGeoIpOrg
)

func CreateGeoLocationAPI(api Api) QueryGeoLocationAPI {
	switch api {
	case IpApiCom:
		return endpoints.IpApiCom{}
	case ReallyFreeGeoIpOrg:
		return endpoints.ReallyFreeGeoIpOrg{}
	}

	return nil
}
