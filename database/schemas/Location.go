package schemas

import (
	geolocation "endlessh-analyzer/geoLocation/structs"
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	geolocation.GeoLocationItem
}
