package analyze

import (
	"bufio"
	"encoding/json"
	"endlessh-analyzer/api"
	cachedb "endlessh-analyzer/cache"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database"
	"endlessh-analyzer/helper"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"strconv"
	"strings"
	time2 "time"
)

type Result struct {
	Tarpitted      int
	SumSeconds     int
	Longest        int
	LongestIp      string
	LongestCountry string
	Shortest       int
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
	Shortest:   math.MaxInt,
}

func DoAnalyze(context *cli.Context) error {
	cachedb.Init(context.Debug)

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
	shortest, _ := db.ExecuteQueryGetFirst(getQueryParametersShortestDuration(start, end))

	result.Tarpitted = int(count)
	result.SumSeconds = int(sum)
	result.Longest = int(longest.Duration)
	result.LongestIp = longest.Ip
	result.Shortest = int(shortest.Duration)

	countryLongest, err := getCountryFor(longest.Ip)
	if err == nil {
		result.LongestCountry = " (" + countryLongest + ")"
	}

	log.Debug("Result object: ", result)

	errOutput := writeConvertedDataToFile(context.Target)
	if errOutput != nil {
		return errOutput
	}

	return nil
}

func getCountryFor(ip string) (string, error) {
	location, cacheResult := cachedb.GetLocationFor(ip)
	geolocationApi := api.CreateGeoLocationAPI(api.IpApiCom)

	if cacheResult == cachedb.NoHit || cacheResult == cachedb.RecordOutdated {
		batch := make([]string, 1)
		batch[0] = ip
		resolved, errApi := geolocationApi.QueryGeoLocationAPI(batch)
		if errApi != nil {
			log.Warningln("Could not get Country for ip: ", ip)
			return "", nil
		}
		location = resolved[0]
		err := cachedb.SaveLocations(resolved)
		if err != nil {
			log.Warningln("Could not save Location in cache for: ", ip)
		}
	} else if cacheResult != cachedb.Ok {
		log.Errorln("Something went wrong for ip: ", ip)
	} else {
		log.Infoln("Got ip: ", ip, " from cache: ", location)
	}

	return location.Country, nil
}

func writeConvertedDataToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Errorln(err)
		return err
	}

	dataWriter := bufio.NewWriter(file)

	timeSum := time2.Date(0, 0, 0, 0, 0, result.SumSeconds, 0, time2.Local)
	timeAvg := time2.Date(0, 0, 0, 0, 0, result.SumSeconds/result.Tarpitted, 0, time2.Local)
	timeLongest := time2.Date(0, 0, 0, 0, 0, result.Longest, 0, time2.Local)
	timeShortest := time2.Date(0, 0, 0, 0, 0, result.Shortest, 0, time2.Local)

	writeToDataWriter(dataWriter, "Tarpitted count:", strconv.Itoa(result.Tarpitted))
	writeToDataWriter(dataWriter, "Tarpitted in sec. (Sum):", strconv.Itoa(result.SumSeconds))
	writeToDataWriter(dataWriter, "Tarpitted in hours. (Sum):", timeSum.Format("15:04:05"))
	writeToDataWriter(dataWriter, "Tarpitted in sec. (Avg):", strconv.Itoa(result.SumSeconds/result.Tarpitted))
	writeToDataWriter(dataWriter, "Tarpitted in hours. (Avg):", timeAvg.Format("15:04:05"))
	writeToDataWriter(dataWriter, "Tarpitted in hours. (Longest):", timeLongest.Format("15:04:05"))
	writeToDataWriter(dataWriter, "Tarpitted IP. (Longest):", result.LongestIp+result.LongestCountry)
	writeToDataWriter(dataWriter, "Tarpitted in hours. (Shortest):", timeShortest.Format("15:04:05"))

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

func writeToDataWriter(dataWriter *bufio.Writer, label string, value string) {
	_, _ = dataWriter.WriteString(strings.TrimSpace(label) + " " + value + "\n")
}
