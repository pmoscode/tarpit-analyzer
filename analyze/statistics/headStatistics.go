package statistics

import (
	"endlessh-analyzer/database"
	"github.com/hako/durafmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
	"strings"
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

	print(p.Sprint(strconv.Itoa(resultTarpitted)))
	data := [][]string{
		{"Attacks count:", p.Sprint(strconv.Itoa(resultTarpitted))},
		{"Attacks sum:", timeSumFormat.String()},
		{"Attacks Avg:", timeAvgFormat.String()},
		{"Attack max from:", resultLongestIp + " => " + timeLongestFormat.String()},
	}

	builder := new(strings.Builder)

	table := tablewriter.NewWriter(builder)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	return builder.String()
}
