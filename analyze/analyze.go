package analyze

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"strconv"
	"strings"
	time2 "time"
)

type Result struct {
	Tarpitted  int
	SumSeconds int
	Longest    int
	Shortest   int
}

var debug = false
var result = Result{
	Tarpitted:  0,
	SumSeconds: 0,
	Longest:    0,
	Shortest:   math.MaxInt,
}

func DoAnalyze(pathSource string, pathTarget string, debugParam bool) error {
	debug = debugParam
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	file, err := os.Open(pathSource)
	if err != nil {
		log.Errorln(err)
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		text := sc.Text()
		err := processLine(text)
		if err != nil {
			log.Warningln(err)
		}
	}

	if err := sc.Err(); err != nil {
		log.Errorln(err)
	}

	log.Debug("Result object: ", result)

	errOutput := writeConvertedDataToFile(pathTarget)
	if errOutput != nil {
		return errOutput
	}

	return nil
}

func processLine(line string) error {
	chunks := strings.Split(line, ",")

	result.Tarpitted = result.Tarpitted + 1

	time, err := strconv.ParseFloat(chunks[2], 32)
	if err != nil {
		return err
	}
	timeInt := int(math.Trunc(time))

	result.SumSeconds = result.SumSeconds + timeInt

	if timeInt < result.Shortest {
		result.Shortest = timeInt
		log.Debugln("New Shortest: " + strconv.Itoa(timeInt))
	} else if timeInt > result.Longest {
		result.Longest = timeInt
		log.Debugln("New Longest: " + strconv.Itoa(timeInt))
	}

	return nil
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

	_, _ = dataWriter.WriteString("Tarpitted count: " + strconv.Itoa(result.Tarpitted) + "\n")
	_, _ = dataWriter.WriteString("Tarpitted in sec. (Sum): " + strconv.Itoa(result.SumSeconds) + "\n")
	_, _ = dataWriter.WriteString("Tarpitted in hours. (Sum): " + timeSum.Format("15:04:05") + "\n")
	_, _ = dataWriter.WriteString("Tarpitted in sec. (Avg): " + strconv.Itoa(result.SumSeconds/result.Tarpitted) + "\n")
	_, _ = dataWriter.WriteString("Tarpitted in hours. (Avg): " + timeAvg.Format("15:04:05") + "\n")
	_, _ = dataWriter.WriteString("Tarpitted in hours. (Longest): " + timeLongest.Format("15:04:05") + "\n")
	_, _ = dataWriter.WriteString("Tarpitted in hours. (Shortest): " + timeShortest.Format("15:04:05") + "\n")

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
