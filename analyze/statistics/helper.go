package statistics

import (
	cachedb "tarpit-analyzer/cache"
)

func getCountryFor(ip string, debug bool) string {
	cachedb.Init(debug)
	location := cachedb.GetLocationFor(ip)

	if location == nil {
		return ""
	}

	return location.Country
}
