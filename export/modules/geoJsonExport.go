package modules

import (
	"endlessh-analyzer/database"
	"endlessh-analyzer/database/schemas"
	"endlessh-analyzer/helper"
	geojson "github.com/paulmach/go.geojson"
	log "github.com/sirupsen/logrus"
	"strconv"
	time2 "time"
)

type GEOJSON struct {
	CenterGeoLocationLongitude float64
	CenterGeoLocationLatitude  float64
	Debug                      bool
	Type                       string // line or point
}

type LineDbItem struct {
	Latitude  float64
	Longitude float64
}

type PointDbItem struct {
	Country     string
	CountryCode string
	Attacks     int64
	LineDbItem
}

func (r *GEOJSON) Export(db *database.Database, start *time2.Time, end *time2.Time) (*[]string, error) {
	whereQueries := make([]database.WhereQuery, 1)
	whereQueries[0] = database.WhereQuery{Query: "l.status == ?", Parameters: "success"}

	var lineData string
	var err error
	switch r.Type {
	case "line":
		lineData, err = r.generateLineData(db, start, end, whereQueries)
		if err != nil {
			return nil, err
		}
	case "point":
		lineData, err = r.generatePointData(db, start, end, whereQueries)
		if err != nil {
			return nil, err
		}
	}

	result := []string{lineData}

	return &result, nil
}

func (r *GEOJSON) generateLineData(db *database.Database, start *time2.Time, end *time2.Time, whereQueries []database.WhereQuery) (string, error) {
	parameter := database.QueryParameters{
		StartDate:   start,
		EndDate:     end,
		SelectQuery: helper.String("l.latitude, l.longitude"),
		Distinct:    true,
		JoinQuery:   helper.String("JOIN locations l on data.ip = l.ip"),
		WhereQuery:  &whereQueries,
	}

	lineDbItem := make([]LineDbItem, 0)
	db.ExecuteQueryGetList(schemas.Data{}, &lineDbItem, parameter)

	featureCollection := geojson.NewFeatureCollection()

	for _, items := range lineDbItem {
		feature := geojson.NewMultiLineStringFeature([][]float64{{items.Longitude, items.Latitude}, {r.CenterGeoLocationLongitude, r.CenterGeoLocationLatitude}})
		featureCollection.AddFeature(feature)
	}

	json, err := featureCollection.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(json), nil
}

func (r *GEOJSON) generatePointData(db *database.Database, start *time2.Time, end *time2.Time, whereQueries []database.WhereQuery) (string, error) {

	// SELECT DISTINCT l.country, l.country_code, c.latitude, c.longitude, count(d.id) as attacks
	// FROM locations l JOIN data d on d.ip = l.ip LEFT JOIN country_geo_locations c ON l.country_code = c.country_code
	// where l.status == 'success'
	// GROUP BY l.country, l.country_code, c.latitude, c.longitude
	// ORDER BY attacks DESC

	parameter := database.QueryParameters{
		StartDate:   start,
		EndDate:     end,
		SelectQuery: helper.String("l.country, l.country_code, c.latitude, c.longitude, count(data.id) as attacks"),
		Distinct:    true,
		JoinQuery:   helper.String("JOIN locations l on data.ip = l.ip LEFT JOIN country_geo_locations c ON l.country_code = c.country_code"),
		WhereQuery:  &whereQueries,
		GroupBy:     helper.String("l.country, l.country_code, c.latitude, c.longitude"),
		OrderBy:     helper.String("attacks DESC"),
	}

	pointDbItem := make([]PointDbItem, 0)
	db.ExecuteQueryGetList(schemas.Data{}, &pointDbItem, parameter)

	featureCollection := geojson.NewFeatureCollection()

	for _, item := range pointDbItem {
		if item.Longitude == 0 && item.Latitude == 0 {
			log.Warningln("Could not resolve this country: ", item.Country, " ## Code: ", item.CountryCode)
		} else {
			feature := geojson.NewPointFeature([]float64{item.Longitude, item.Latitude})
			feature.SetProperty("name", item.Country+" ("+strconv.FormatInt(item.Attacks, 10)+")")
			featureCollection.AddFeature(feature)
		}
	}

	json, err := featureCollection.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(json), nil
}
