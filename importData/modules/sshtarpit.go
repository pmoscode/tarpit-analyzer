package modules

import (
	"bufio"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/importData/structs"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
	"strings"
	time2 "time"
)

type SshTarpitItem struct {
	start    time2.Time
	end      time2.Time
	duration float32
	ip       string
	port     int64
}

type SshTarpit struct {
	debug bool
}

func (r SshTarpit) Import(sourcePath string, context *cli.Context) (*[]structs.ImportItem, error) {
	r.debug = context.Debug
	if r.debug {
		log.SetLevel(log.DebugLevel)
	}

	items := make([]structs.ImportItem, 0)
	tempMap := make(map[int64]SshTarpitItem)

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
		item, err := r.processLine(tempMap, text)
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

func (r SshTarpit) processLine(tempMap map[int64]SshTarpitItem, line string) (structs.ImportItem, error) {
	if strings.Contains(line, "connected") || strings.Contains(line, "disconnected") {
		space := regexp.MustCompile(`\s+`)
		cleanedString := space.ReplaceAllString(line, " ")
		chunks := strings.Split(cleanedString, " ")

		dateStr := chunks[0] + " " + chunks[1]
		date, errDate := time2.Parse("2006-01-02 15:04:05", dateStr)
		if errDate != nil {
			return structs.ImportItem{Success: false}, errDate
		}

		ip := strings.Split(chunks[5], "'")[1]

		port, errPort := strconv.ParseInt(strings.TrimRight(chunks[6], ")"), 10, 64)
		if errPort != nil {
			log.Errorln(errPort)
			return structs.ImportItem{Success: false}, errPort
		}

		if value, exist := tempMap[port]; exist {
			value.end = date
			value.duration = float32(date.Sub(value.start).Seconds())

			delete(tempMap, port)

			return r.mapToImportItem(value)
		} else {
			value := SshTarpitItem{
				start: date,
				ip:    ip,
				port:  port,
			}
			tempMap[port] = value
		}
	}

	return structs.ImportItem{Success: false}, nil
}

func (r SshTarpit) mapToImportItem(value SshTarpitItem) (structs.ImportItem, error) {
	return structs.ImportItem{
		Begin:    value.start,
		End:      value.end,
		Duration: value.duration,
		Ip:       value.ip,
		Success:  true,
	}, nil
}
