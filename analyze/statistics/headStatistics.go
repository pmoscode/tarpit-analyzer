package statistics

import (
	"github.com/hako/durafmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
	"strings"
	"tarpit-analyzer/database"
	time2 "time"
)

func GetHeadStatistics(db *database.DbData, start *time2.Time, end *time2.Time, debug bool) string {
	count, _ := db.ExecuteQueryGetAggregator(getQueryParametersCountAll(start, end))
	sum, _ := db.ExecuteQueryGetAggregator(getQueryParametersSumAll(start, end))
	longest, _ := db.ExecuteQueryGetFirst(getQueryParametersLongestDuration(start, end))

	// Count of tarpitted attacker
	resultTarpitted := int(count)

	// Time sum of attacks
	resultSumSeconds := int(sum)
	timeSum, _ := time2.ParseDuration(strconv.Itoa(resultSumSeconds) + "s")
	timeSumFormat, _ := durafmt.ParseString(timeSum.String())

	// Time avg of attacks
	timeAvg, _ := time2.ParseDuration(strconv.Itoa(resultSumSeconds/resultTarpitted) + "s")
	timeAvgFormat, _ := durafmt.ParseString(timeAvg.String())

	// Time max of attack
	resultLongest := int(longest.Duration)
	timeLongest, _ := time2.ParseDuration(strconv.Itoa(resultLongest) + "s")
	timeLongestFormat, _ := durafmt.ParseString(timeLongest.String())

	// Ip of max attack
	resultLongestIp := longest.Ip
	countryLongest := getCountryFor(longest.Ip, debug)
	if countryLongest != "" {
		resultLongestIp = resultLongestIp + " (" + countryLongest + ")"
	}

	p := message.NewPrinter(language.English)

	data := [][]string{
		{"Count", p.Sprint(strconv.Itoa(resultTarpitted))},
		{"Time Sum", timeSumFormat.String()},
		{"Time Avg", timeAvgFormat.String()},
		{"IP max time", resultLongestIp + " => " + timeLongestFormat.String()},
	}

	builder := new(strings.Builder)

	table := tablewriter.NewWriter(builder)
	table.SetHeader([]string{"Attack", ""})
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(data)
	table.Render()

	return "  GLOBAL STATISTICS\n" + builder.String()
}
