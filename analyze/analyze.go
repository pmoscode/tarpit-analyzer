package analyze

import (
	"endlessh-analyzer/analyze/statistics"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database"
	"endlessh-analyzer/helper"
	log "github.com/sirupsen/logrus"
	time2 "time"
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
	targetFileWriter.writeTextWithBottomPadding("==================================", 1)
	header := addDateHeader(start, end)
	if header != "" {
		targetFileWriter.writeTextWithBottomPadding("Selected date range: "+header, 2)
	} else {
		targetFileWriter.writeTextWithBottomPadding("Selected date range: All data", 2)
	}

	headStat := statistics.GetHeadStatistics(&db, start, end, context.Debug)
	targetFileWriter.writeTextWithBottomPadding(headStat, 1)

	topStat, errTopStat := statistics.GetTopStatistics(&db, start, end, context.Debug)
	if errTopStat != nil {
		return errTopStat
	}
	targetFileWriter.writeTextWithBottomPadding(topStat, 1)

	weekdayStat, errWeekdayStat := statistics.GetAttackTimeStatistics(&db, start, end, statistics.DAY)
	if errWeekdayStat != nil {
		return errWeekdayStat
	}
	targetFileWriter.writeTextWithBottomPadding(weekdayStat, 1)

	monthStat, errMonthStat := statistics.GetAttackTimeStatistics(&db, start, end, statistics.MONTH)
	if errMonthStat != nil {
		return errMonthStat
	}
	targetFileWriter.writeTextWithBottomPadding(monthStat, 1)

	yearStat, errYearStat := statistics.GetAttackTimeStatistics(&db, start, end, statistics.YEAR)
	if errYearStat != nil {
		return errYearStat
	}
	targetFileWriter.writeTextWithBottomPadding(yearStat, 1)

	errClose := targetFileWriter.close()
	if errClose != nil {
		return errClose
	}

	return nil
}

func addDateHeader(start *time2.Time, end *time2.Time) string {
	header := ""

	if start != nil {
		header = start.Format("2006-01-02")
	}

	if end != nil {
		if header != "" {
			header += " - "
		}
		header += end.Format("2006-01-02")
	}

	return header
}
