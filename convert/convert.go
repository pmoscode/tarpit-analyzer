package convert

import (
	"bufio"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	time2 "time"
)

var debug = false

func DoConvert(pathSource string, pathTarget string, startDate string, endDate string, debugParam bool) error {
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
	lines := make([]string, 0)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		text := sc.Text()
		line, err := processLine(text, startDate, endDate)
		if err != nil {
			log.Warningln(err)
		} else {
			if line != "" {
				lines = append(lines, line)
			}
		}
	}

	if debug {
		for i := 0; i < len(lines); i++ {
			log.Debugln(lines[i])
		}

		if len(lines) == 0 {
			log.Debugln("Lines are empty... :(")
		}
	}

	if err := sc.Err(); err != nil {
		log.Errorln(err)
	}

	errOutput := writeConvertedDataToFile(pathTarget, lines)
	if errOutput != nil {
		return errOutput
	}

	return nil
}

func processLine(line string, startDate string, endDate string) (string, error) {
	if strings.Contains(line, "CLOSE") {
		chunks := strings.Split(line, " ")
		log.Debugln("Chunks: ", chunks)

		if checkStartDate(startDate, chunks[0]) && checkEndDate(endDate, chunks[0]) {
			date, errDate := time2.Parse(time2.RFC3339Nano, chunks[0])
			if errDate != nil {
				return "", errDate
			}
			ip, errIp := getValue(chunks[2])
			time, errTime := getValue(chunks[5])

			if errTime != nil && errIp != nil {
				log.Errorln(errIp, errTime)
				return "", errors.New("line '" + line + "' could not be cleaned")
			}
			dataString := date.Format("2006-01-02 15:04:05") + "," + ip + "," + time

			return dataString, nil
		}
	}

	return "", nil
}

func checkStartDate(startDateParam string, startDateLine string) bool {
	if startDateParam == "unset" {
		return true
	}

	startingDateParam, startingDateLine, err := prepareDates(startDateParam, startDateLine)

	if err != nil {
		return true
	}

	return startingDateLine.After(startingDateParam)
}

func checkEndDate(endDateParam string, endDateLine string) bool {
	if endDateParam == "unset" {
		return true
	}

	endingDateParam, endingDateLine, err := prepareDates(endDateParam, endDateLine)

	if err != nil {
		return true
	}

	return endingDateLine.Before(endingDateParam.AddDate(0, 0, 1))
}

func prepareDates(dateParam string, dateLine string) (time2.Time, time2.Time, error) {
	startingDateParam, errParam := time2.Parse("2006-01-02", dateParam)
	startingDateLine, errLine := time2.Parse(time2.RFC3339Nano, dateLine)

	if errLine != nil && errParam != nil {
		log.Errorln("Start dates cannot be parsed: parameter -- ", dateParam, " ## data -- ", dateLine)
		log.Errorln(errParam, errLine)

		return time2.Now(), time2.Now(), errors.New("nothing")
	}

	if debug {
		log.Debugln("Dates to check: ", startingDateParam, " and ", startingDateLine)
	}

	return startingDateParam, startingDateLine, nil
}

func getValue(source string) (string, error) {
	if source == "" {
		return "", errors.New("source is empty")
	}

	split := strings.Split(source, "=")

	if len(split) == 1 {
		return "", errors.New("source has no '=' char")
	}

	return split[1], nil
}

func writeConvertedDataToFile(path string, dataLines []string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Errorln(err)
		return err
	}

	dataWriter := bufio.NewWriter(file)

	for _, data := range dataLines {
		_, _ = dataWriter.WriteString(data + "\n")
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
