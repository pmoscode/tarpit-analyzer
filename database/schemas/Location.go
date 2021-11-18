package schemas

import (
	geolocation "endlessh-analyzer/api/structs"
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	geolocation.GeoLocationItem
}
