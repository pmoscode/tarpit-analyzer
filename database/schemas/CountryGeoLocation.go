package schemas

import (
	"gorm.io/gorm"
)

type CountryGeoLocation struct {
	gorm.Model
	CountryCode string
	Latitude    float64
	Longitude   float64
	Name        string
}
