package export

import (
	"log"
	"strconv"
	"tarpit-analyzer/cli"
	"tarpit-analyzer/database"
	"tarpit-analyzer/export/helper"
	"tarpit-analyzer/export/modules"
	time2 "time"
)

type Parameters struct {
	Separator                  string
	CenterGeoLocationLatitude  string
	CenterGeoLocationLongitude string
	Type                       string
}

type Export interface {
	Export(database *database.Database, start *time2.Time, end *time2.Time) (*[]string, error)
}

type Type int

const (
	CSV Type = iota
	JSON
	KML
	GEOJSON
)

func DoExport(exportType Type, parameters Parameters, context *cli.Context) error {
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
			Type:                       parameters.Type,
		}
	}

	data := helper.PrepareDatabase(context)
	start, end := helper.PrepareTimeBounds(context)

	exportData, err := exporter.Export(data, start, end)
	if err != nil {
		return err
	}

	err = helper.WriteDataToFile(context.Target, exportData)
	if err != nil {
		return err
	}

	return nil
}
