package analyze

import (
	"endlessh-analyzer/database"
	"endlessh-analyzer/helper"
	time2 "time"
)

func getQueryParametersCountAll(start *time2.Time, end *time2.Time) database.QueryParameters {
	return database.QueryParameters{
		StartDate:   start,
		EndDate:     end,
		SelectQuery: helper.String("count(*) as count"),
	}
}

func getQueryParametersSumAll(start *time2.Time, end *time2.Time) database.QueryParameters {
	return database.QueryParameters{
		StartDate:   start,
		EndDate:     end,
		SelectQuery: helper.String("SUM(duration) as duration"),
	}
}

func getQueryParametersLongestDuration(start *time2.Time, end *time2.Time) database.QueryParameters {
	return database.QueryParameters{
		StartDate: start,
		EndDate:   end,
		OrderBy:   helper.String("duration desc"),
	}
}

func getQueryParametersShortestDuration(start *time2.Time, end *time2.Time) database.QueryParameters {
	return database.QueryParameters{
		StartDate: start,
		EndDate:   end,
		OrderBy:   helper.String("duration"),
	}
}
