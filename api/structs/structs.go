package structs

type GeoLocationItem struct {
	Ip            string `gorm:"uniqueIndex"`
	Status        string
	Latitude      float64
	Longitude     float64
	Continent     string
	ContinentCode string
	Country       string
	CountryCode   string
	Region        string
	RegionName    string
	City          string
	Zip           string
}
