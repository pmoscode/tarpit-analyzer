package analyze

import (
	"endlessh-analyzer/analyze/statistics"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database"
	"endlessh-analyzer/helper"
	log "github.com/sirupsen/logrus"
)

func DoAnalyze(context *cli.Context) error {
	db, errCreate := database.CreateDbData(context.Debug)
	if errCreate != nil {
		log.Panicln("Data database could not be loaded.", errCreate)
	}

	if context.Debug {
		log.SetLevel(log.DebugLevel)
	}

	start := helper.GetDate(context.StartDate)
	end := helper.GetDate(context.EndDate)

	targetFileWriter := textFileWriter{}
	err := targetFileWriter.openFileForWrite(context.Target)
	if err != nil {
		return err
	}

	targetFileWriter.writeText("\tTarpit Analyzer Statistics")
	targetFileWriter.writeText("==================================")
	targetFileWriter.writeText("")

	headStat := statistics.GetHeadStatistics(&db, start, end, context.Debug)
	targetFileWriter.writeText(headStat)

	topStat, errTopStat := statistics.GetTopStatistics(&db, start, end, context.Debug)
	if errTopStat != nil {
		return errTopStat
	}
	targetFileWriter.writeText(topStat)

	errClose := targetFileWriter.close()
	if errClose != nil {
		return errClose
	}

	return nil
}
