package api

import (
	"endlessh-analyzer/api/endpoints"
	"endlessh-analyzer/api/structs"
)

type GeoLocation interface {
	QueryGeoLocationApi(ips []string) ([]structs.GeoLocationItem, error)
}

type Api int

const (
	IpApiCom Api = iota
)

func CreateGeoLocationApi(api Api) GeoLocation {
	switch api {
	case IpApiCom:
		return endpoints.IpApiCom{}
	}

	return nil
}
