package export

import (
	"bufio"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database"
	"endlessh-analyzer/database/schemas"
	"endlessh-analyzer/export/modules"
	"endlessh-analyzer/helper"
	log "github.com/sirupsen/logrus"
	"os"
	time2 "time"
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

func getData(context *cli.Context) []schemas.Data {
	// Get Data DB connection
	db := createDB(context)

	// Setup Debug logger (or not)
	setupLogger(context)

	// Get start and end date from CLI params
	start := helper.GetDate(context.StartDate)
	end := helper.GetDate(context.EndDate)

	// Query data
	data := queryDB(start, end, db)

	return data
}

func setupLogger(context *cli.Context) {
	debug = context.Debug
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func queryDB(start *time2.Time, end *time2.Time, db database.DbData) []schemas.Data {
	queryParameters := database.QueryParameters{
		StartDate: start,
		EndDate:   end,
	}
	data, _ := db.ExecuteQueryGetList(queryParameters)

	return data
}

func createDB(context *cli.Context) database.DbData {
	// Load data with parameters
	db, errCreate := database.CreateDbData(context.Debug)
	if errCreate != nil {
		log.Panicln("Data database could not be loaded.", errCreate)
	}

	return db
}

func writeDataToFile(path string, exportData *[]string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Errorln(err)
		return err
	}

	dataWriter := bufio.NewWriter(file)

	for _, str := range *exportData {
		_, _ = dataWriter.WriteString(str + "\n")
	}

	errWriter := dataWriter.Flush()
	if errWriter != nil {
		return errWriter
	}

	errFile := file.Close()
	if errFile != nil {
		return errFile
	}

	return nil
}
