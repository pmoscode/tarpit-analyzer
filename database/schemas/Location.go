package schemas

import (
	"gorm.io/gorm"
	geolocation "tarpit-analyzer/geoLocation/structs"
)

type Location struct {
	gorm.Model
	geolocation.GeoLocationItem
}
