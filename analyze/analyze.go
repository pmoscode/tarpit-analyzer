package analyze

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"endlessh-analyzer/api"
	cachedb "endlessh-analyzer/cache"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database"
	"endlessh-analyzer/helper"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"strconv"
	time2 "time"
)

type Chart struct {
	Country    string `gorm:"column:country"`
	SumAttacks int64  `gorm:"column:sum_attacks"`
	SumTime    int64  `gorm:"column:sum_time"`
	AvgTime    int64  `gorm:"column:avg_time"`
}

type Result struct {
	Tarpitted      int
	SumSeconds     int
	Longest        int
	LongestIp      string
	LongestCountry string
	Charts         []Chart
}

func (r Result) String() string {
	s, _ := json.MarshalIndent(r, "", "\t")
	return string(s)
}

var debug = false
var result = Result{
	Tarpitted:  0,
	SumSeconds: 0,
	Longest:    0,
}

func DoAnalyze(context *cli.Context) error {
	db, errCreate := database.CreateDbData(context.Debug)
	if errCreate != nil {
		log.Panicln("Data database could not be loaded.", errCreate)
	}

	debug = context.Debug
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	start := helper.GetDate(context.StartDate)
	end := helper.GetDate(context.EndDate)

	count, _ := db.ExecuteQueryGetAggregator(getQueryParametersCountAll(start, end))
	sum, _ := db.ExecuteQueryGetAggregator(getQueryParametersSumAll(start, end))
	longest, _ := db.ExecuteQueryGetFirst(getQueryParametersLongestDuration(start, end))
	rowsTop, errQuery := db.DbRawQuery(getRawTopCountriesAttacks(start, end))
	if errQuery != nil {
		return errQuery
	}

	charts := make([]Chart, 0)
	for rowsTop.Next() {
		var chart Chart
		err := db.ScanToStruct(rowsTop, &chart)
		if err != nil {
			return err
		}
		charts = append(charts, chart)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Errorln("Could not close query in DB.")
		}
	}(rowsTop)

	result.Tarpitted = int(count)
	result.SumSeconds = int(sum)
	result.Longest = int(longest.Duration)
	result.LongestIp = longest.Ip
	result.Charts = charts

	countryLongest := getCountryFor(longest.Ip, context.Debug)
	if countryLongest != "" {
		result.LongestCountry = " (" + countryLongest + ")"
	}

	log.Debug("Result object: ", result)

	errOutput := writeConvertedDataToFile(context.Target)
	if errOutput != nil {
		return errOutput
	}

	return nil
}

func getCountryFor(ip string, debug bool) string {
	cachedb.Init(api.IpApiCom, debug)
	location := cachedb.GetLocationFor(ip)

	if location == nil {
		return ""
	}

	return location.Country
}

func writeConvertedDataToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Errorln(err)
		return err
	}

	dataWriter := bufio.NewWriter(file)
	p := message.NewPrinter(language.English)

	timeSum, _ := time2.ParseDuration(strconv.Itoa(result.SumSeconds) + "s")
	timeAvg, _ := time2.ParseDuration(strconv.Itoa(result.SumSeconds/result.Tarpitted) + "s")
	timeLongest, _ := time2.ParseDuration(strconv.Itoa(result.Longest) + "s")

	timeSumFormat, _ := durafmt.ParseString(timeSum.String())
	timeAvgFormat, _ := durafmt.ParseString(timeAvg.String())
	timeLongestFormat, _ := durafmt.ParseString(timeLongest.String())

	writeToDataWriter(dataWriter, "Overall tarpitted count: "+p.Sprint(result.Tarpitted))
	writeToDataWriter(dataWriter, "Sum of tarpitted: "+timeSumFormat.String())
	writeToDataWriter(dataWriter, "Average tarpitted: "+timeAvgFormat.String())
	writeToDataWriter(dataWriter, "Longest tarpitted IP: "+result.LongestIp+result.LongestCountry+" => "+timeLongestFormat.String())

	for idx, chart := range result.Charts {
		sumT, _ := time2.ParseDuration(strconv.Itoa(int(chart.SumTime)) + "s")
		avgT, _ := time2.ParseDuration(strconv.Itoa(int(chart.AvgTime)) + "s")

		sumTFormat, _ := durafmt.ParseString(sumT.String())
		avgTFormat, _ := durafmt.ParseString(avgT.String())

		writeToDataWriter(dataWriter, "TOP "+strconv.Itoa(idx+1)+" attacker from "+chart.Country+":")
		writeToDataWriter(dataWriter, "\tAttacks: "+p.Sprint(chart.SumAttacks))
		writeToDataWriter(dataWriter, "\tOverall attack time: "+sumTFormat.String())
		writeToDataWriter(dataWriter, "\tAverage attack time: "+avgTFormat.String())
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

func writeToDataWriter(dataWriter *bufio.Writer, value string) {
	_, _ = dataWriter.WriteString(value + "\n")
}
