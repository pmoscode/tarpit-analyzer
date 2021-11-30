package modules

import (
	"endlessh-analyzer/api"
	"endlessh-analyzer/api/structs"
	cachedb "endlessh-analyzer/cache"
	"endlessh-analyzer/database/schemas"
	geojson "github.com/paulmach/go.geojson"
)

type GEOJSON struct {
	CenterGeoLocationLongitude float64
	CenterGeoLocationLatitude  float64
	Debug                      bool
}

func (r *GEOJSON) Export(data *[]schemas.Data) (*[]string, error) {
	cachedb.Init(api.IpApiCom, r.Debug)
	result := make([]string, 0)

	locations := make([]structs.GeoLocationItem, 0)
	for _, dataItem := range *data {
		location := cachedb.GetLocationFor(dataItem.Ip)
		locations = append(locations, *location)
	}

	locations = uniqueNonEmptyElementsOf(&locations)

	featureCollection := geojson.NewFeatureCollection()

	for _, items := range locations {
		if items.Status == "success" {
			feature := geojson.NewMultiLineStringFeature([][]float64{{items.Longitude, items.Latitude}, {r.CenterGeoLocationLongitude, r.CenterGeoLocationLatitude}})
			featureCollection.AddFeature(feature)
		}
	}

	json, err := featureCollection.MarshalJSON()
	if err != nil {
		return nil, err
	}

	result = append(result, string(json))

	return &result, nil
}
