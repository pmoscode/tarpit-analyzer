package structs

type GeoLocationItem struct {
	Ip          string `gorm:"uniqueIndex"`
	Status      string
	Latitude    float64
	Longitude   float64
	Country     string
	CountryCode string
	Region      string
	RegionName  string
	City        string
	Zip         string
}
