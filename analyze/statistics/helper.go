package statistics

import (
	"endlessh-analyzer/api"
	cachedb "endlessh-analyzer/cache"
)

func getCountryFor(ip string, debug bool) string {
	cachedb.Init(api.IpApiCom, debug)
	location := cachedb.GetLocationFor(ip)

	if location == nil {
		return ""
	}

	return location.Country
}
