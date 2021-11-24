package export

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database/schemas"
	"endlessh-analyzer/export/modules"
)

var debug = false

type Export interface {
	Export(data *[]schemas.Data) (*[]string, error)
}

func CSV(separator string, context *cli.Context) error {
	data := getData(context)

	var exporter Export
	exporter = &modules.CSV{Separator: separator}
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

func JSON(context *cli.Context) error {
	data := getData(context)

	var exporter Export
	exporter = &modules.JSON{}
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

func KML(context *cli.Context) error {
	data := getData(context)

	var exporter Export
	exporter = &modules.KML{}
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
