package statistics

import (
	"tarpit-analyzer/database"
	"tarpit-analyzer/helper"
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

func getQueryParametersTopCountriesAttacks(start *time2.Time, end *time2.Time) database.QueryParameters {
	return database.QueryParameters{
		StartDate:   start,
		EndDate:     end,
		SelectQuery: helper.String("l.country, count(data.id) as sum_attacks, CAST(sum(data.duration) as INT) as sum_time, CAST(round(avg(data.duration), 0) as INT) as avg_time"),
		JoinQuery:   helper.String("JOIN locations l ON data.ip = l.ip"),
		GroupBy:     helper.String("l.country"),
		OrderBy:     helper.String("sum_attacks desc"),
		Limit:       helper.Int(5),
	}
}

func getQueryParametersDateAttacks(start *time2.Time, end *time2.Time, mode TimeStatistic) database.QueryParameters {
	return database.QueryParameters{
		StartDate:   start,
		EndDate:     end,
		SelectQuery: helper.String("strftime('" + mode.format + "', begin) as mode_str, count(id) as attacks"),
		WhereQuery: &[]database.WhereQuery{{
			Query:      "success = ?",
			Parameters: 1,
		}},
		GroupBy: helper.String("mode_str"),
		OrderBy: helper.String("attacks DESC"),
	}
}
