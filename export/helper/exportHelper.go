package helper

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"tarpit-analyzer/cli"
	"tarpit-analyzer/database"
	"tarpit-analyzer/database/schemas"
	"tarpit-analyzer/helper"
	time2 "time"
)

func PrepareDatabase(context *cli.Context) *database.Database {
	db := createDB(context)
	setupLogger(context)

	return db
}

func PrepareTimeBounds(context *cli.Context) (*time2.Time, *time2.Time) {
	start := helper.GetDate(context.StartDate)
	end := helper.GetDate(context.EndDate)

	return start, end
}

func setupLogger(context *cli.Context) {
	if context.Debug {
		log.SetLevel(log.DebugLevel)
	}
}

func createDB(context *cli.Context) *database.Database {
	db, errCreate := database.CreateGenericDatabase(context.Debug)
	if errCreate != nil {
		log.Panicln("Data database could not be loaded.", errCreate)
	}

	return db
}

func QueryDataDB(db *database.Database, start *time2.Time, end *time2.Time) *[]schemas.Data {
	queryParameters := database.QueryParameters{
		StartDate: start,
		EndDate:   end,
	}

	data := make([]schemas.Data, 0)

	_ = db.ExecuteQueryGetList(&schemas.Data{}, &data, queryParameters)

	return &data
}

func WriteDataToFile(path string, exportData *[]string) error {
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
