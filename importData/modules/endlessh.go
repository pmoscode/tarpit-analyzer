package modules

import (
	"bufio"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/importData/structs"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	time2 "time"
)

type Endlessh struct {
	debug bool
}

func (r Endlessh) Import(sourcePath string, context *cli.Context) (*[]structs.ImportItem, error) {
	r.debug = context.Debug
	if r.debug {
		log.SetLevel(log.DebugLevel)
	}

	items := make([]structs.ImportItem, 0)

	file, err := os.Open(sourcePath)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		text := sc.Text()
		item, err := r.processLine(text)
		if err != nil {
			log.Warningln(err)
		} else {
			if item.Success {
				items = append(items, item)
			}
		}
	}

	if len(items) == 0 {
		log.Debugln("No data found... :(")
	}

	if err := sc.Err(); err != nil {
		log.Errorln(err)
	}

	return &items, nil
}

func (r Endlessh) processLine(line string) (structs.ImportItem, error) {
	if strings.Contains(line, "CLOSE") {
		chunks := strings.Split(line, " ")

		date, errDate := time2.Parse(time2.RFC3339Nano, chunks[0])
		if errDate != nil {
			return structs.ImportItem{Success: false}, errDate
		}
		ip, errIp := r.getValue(chunks[2])
		time, errTime := r.getValue(chunks[5])

		if errTime != nil && errIp != nil {
			log.Errorln(errIp, errTime)
			return structs.ImportItem{Success: false}, errors.New("line '" + line + "' could not be cleaned")
		}

		timeFloat, errParse := strconv.ParseFloat(time, 32)
		if errParse != nil {
			return structs.ImportItem{Success: false}, errParse
		}

		durationTime, errDuration := time2.ParseDuration(time + "s")
		if errDuration != nil {
			return structs.ImportItem{Success: false}, errDuration
		}

		item := structs.ImportItem{
			Begin:    date.Add(-durationTime),
			End:      date,
			Duration: float32(timeFloat),
			Ip:       ip,
			Success:  true,
		}

		return item, nil
	}

	return structs.ImportItem{Success: false}, nil
}

func (r Endlessh) getValue(source string) (string, error) {
	if source == "" {
		return "", errors.New("source is empty")
	}

	split := strings.Split(source, "=")

	if len(split) == 1 {
		return "", errors.New("source has no '=' char")
	}

	return split[1], nil
}
