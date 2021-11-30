package modules

import (
	"endlessh-analyzer/api"
	"endlessh-analyzer/api/structs"
	cachedb "endlessh-analyzer/cache"
	"endlessh-analyzer/database/schemas"
	"fmt"
)

type KML struct {
	CenterGeoLocationLongitude string
	CenterGeoLocationLatitude  string
	Debug                      bool
}

func (r *KML) Export(data *[]schemas.Data) (*[]string, error) {
	cachedb.Init(api.IpApiCom, r.Debug)
	result := make([]string, 0)

	locations := make([]structs.GeoLocationItem, 0)
	for _, dataItem := range *data {
		location := cachedb.GetLocationFor(dataItem.Ip)
		locations = append(locations, *location)
	}

	locations = uniqueNonEmptyElementsOf(&locations)

	r.generateKMLContent(&result, &locations)

	return &result, nil
}

func (r *KML) generateKMLContent(result *[]string, data *[]structs.GeoLocationItem) {
	*result = append(*result, `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
	<Document>
		<Style id="transBluePoly">
			<LineStyle>
				<width>1.5</width>
				<color>501400E6</color>
			</LineStyle>
		</Style>`)

	for _, items := range *data {
		if items.Status == "success" {
			*result = append(*result, `
		<Placemark>
			<name>`+items.Country+`</name>
			<extrude>1</extrude>
			<tessellate>1</tessellate>
			<styleUrl>#transBluePoly</styleUrl>
			<LineString>
				<coordinates>
					`+fmt.Sprintf("%f", items.Longitude)+`, `+fmt.Sprintf("%f", items.Latitude)+`
					`+r.CenterGeoLocationLongitude+`,`+r.CenterGeoLocationLatitude+`
				</coordinates>
			</LineString>
		</Placemark>`)
		}
	}

	*result = append(*result, `
	</Document>
</kml>`)
}
