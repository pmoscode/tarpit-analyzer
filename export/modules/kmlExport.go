package modules

import (
	"endlessh-analyzer/database"
	"endlessh-analyzer/database/schemas"
	"endlessh-analyzer/helper"
	"fmt"
	time2 "time"
)

type KML struct {
	CenterGeoLocationLongitude string
	CenterGeoLocationLatitude  string
	Debug                      bool
}

type KmlDbItem struct {
	Country   string
	Latitude  float64
	Longitude float64
}

func (r *KML) Export(db *database.Database, start *time2.Time, end *time2.Time) (*[]string, error) {
	whereQueries := make([]database.WhereQuery, 0)
	whereQueries = append(whereQueries, database.WhereQuery{Query: "l.status == ?", Parameters: "success"})
	if start != nil {
		whereQueries = append(whereQueries, database.WhereQuery{Query: "d.begin >= ?", Parameters: start})
	}
	if end != nil {
		whereQueries = append(whereQueries, database.WhereQuery{Query: "d.end <= ?", Parameters: end})
	}
	parameter := database.QueryParameters{
		SelectQuery: helper.String("l.country, l.latitude, l.longitude"),
		Distinct:    true,
		JoinQuery:   helper.String("JOIN locations l on data.ip = l.ip"),
		WhereQuery:  &whereQueries,
	}

	result := make([]string, 0)
	kmlDbItems := make([]KmlDbItem, 0)

	db.ExecuteQueryGetList(schemas.Data{}, &kmlDbItems, parameter)

	r.generateKMLContent(&result, &kmlDbItems)

	return &result, nil
}

func (r *KML) generateKMLContent(result *[]string, data *[]KmlDbItem) {
	*result = append(*result, `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
	<Document>
		<Style id="transBluePoly">
			<LineStyle>
				<width>1.5</width>
				<color>501400E6</color>
			</LineStyle>
		</Style>`)

	for _, item := range *data {
		*result = append(*result, `
		<Placemark>
			<name>`+item.Country+`</name>
			<extrude>1</extrude>
			<tessellate>1</tessellate>
			<styleUrl>#transBluePoly</styleUrl>
			<LineString>
				<coordinates>
					`+fmt.Sprintf("%f", item.Longitude)+`, `+fmt.Sprintf("%f", item.Latitude)+`
					`+r.CenterGeoLocationLongitude+`,`+r.CenterGeoLocationLatitude+`
				</coordinates>
			</LineString>
		</Placemark>`)
	}

	*result = append(*result, `
	</Document>
</kml>`)
}
