package api

import (
	"endlessh-analyzer/api/endpoints"
	"endlessh-analyzer/api/structs"
)

type QueryGeoLocationAPI interface {
	QueryGeoLocationAPI(ips []string) ([]structs.GeoLocationItem, error)
}

type Api int

const (
	IpApiCom Api = iota
)

func CreateGeoLocationAPI(api Api) QueryGeoLocationAPI {
	switch api {
	case IpApiCom:
		return endpoints.IpApiCom{}
	}

	return nil
}
