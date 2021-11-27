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

func getRawTopCountriesAttacks(start *time2.Time, end *time2.Time) string {
	whereStringStart := ""
	if start != nil {
		whereStringStart = "d.begin >= \"" + start.Format("2006-01-02") + "\" "
	}

	whereStringEnd := ""
	if end != nil {
		if whereStringStart != "" {
			whereStringEnd = " and "
		}
		whereStringEnd = whereStringEnd + "d.end <= \"" + end.Format("2006-01-02") + "\" "
	}

	whereString := ""
	if whereStringStart != "" || whereStringEnd != "" {
		whereString = "WHERE " + whereStringStart + whereStringEnd + " "
	}

	return "SELECT l.country, count(d.id) as sum_attacks, CAST(sum(d.duration) as INT) as sum_time, CAST(round(avg(d.duration), 0) as INT) as avg_time " +
		"from data d JOIN locations l ON d.ip = l.ip " +
		whereString +
		"GROUP BY l.country " +
		"ORDER BY sum_attacks DESC " +
		"LIMIT 5"
}
