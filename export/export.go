package export

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database/schemas"
	"endlessh-analyzer/export/modules"
	"log"
	"strconv"
)

var debug = false

type Parameters struct {
	Separator                  string
	CenterGeoLocationLatitude  string
	CenterGeoLocationLongitude string
}

type Export interface {
	Export(data *[]schemas.Data) (*[]string, error)
}

type Type int

const (
	CSV Type = iota
	JSON
	KML
	GEOJSON
)

func DoExport(exportType Type, parameters Parameters, context *cli.Context) error {
	data := getData(context)

	var exporter Export
	switch exportType {
	case CSV:
		exporter = &modules.CSV{Separator: parameters.Separator}
	case JSON:
		exporter = &modules.JSON{}
	case KML:
		exporter = &modules.KML{
			CenterGeoLocationLongitude: parameters.CenterGeoLocationLongitude,
			CenterGeoLocationLatitude:  parameters.CenterGeoLocationLatitude,
			Debug:                      context.Debug,
		}
	case GEOJSON:
		lon, errLon := strconv.ParseFloat(parameters.CenterGeoLocationLongitude, 64)
		lat, errLat := strconv.ParseFloat(parameters.CenterGeoLocationLatitude, 64)

		if errLon != nil || errLat != nil {
			log.Fatalln("CenterGeoLocation parameter must be a valid number")
		}

		exporter = &modules.GEOJSON{
			CenterGeoLocationLongitude: lon,
			CenterGeoLocationLatitude:  lat,
			Debug:                      context.Debug,
		}
	}

	exportData, err := exporter.Export(&data)
	if err != nil {
		return err
	}

	err = writeDataToFile(context.Target, exportData)
	if err != nil {
		return err
	}

	return nil
}
