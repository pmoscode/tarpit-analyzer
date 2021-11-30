package modules

import "endlessh-analyzer/api/structs"

func uniqueNonEmptyElementsOf(items *[]structs.GeoLocationItem) []structs.GeoLocationItem {
	unique := make(map[string]bool, len(*items))
	us := make([]structs.GeoLocationItem, len(unique))
	for _, elem := range *items {
		if len(elem.Ip) != 0 {
			if !unique[elem.Ip] {
				us = append(us, elem)
				unique[elem.Ip] = true
			}
		}
	}

	return us
}
