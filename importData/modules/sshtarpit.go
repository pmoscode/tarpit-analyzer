package modules

import (
	"bufio"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/helper"
	"endlessh-analyzer/importData/structs"
	"github.com/schollz/progressbar/v3"
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

func (r SshTarpit) Import(sourcePath string, start *time2.Time, end *time2.Time, context *cli.Context) (*[]structs.ImportItem, int, int, error) {
	r.debug = context.Debug
	if r.debug {
		log.SetLevel(log.DebugLevel)
	}

	items := make([]structs.ImportItem, 0)
	tempMap := make(map[int64]SshTarpitItem)

	file, err := os.Open(sourcePath)
	if err != nil {
		log.Errorln(err)
		return nil, 0, 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	sc := bufio.NewScanner(file)
	processedLine := 1
	skipLines := 0
	bar := progressbar.Default(-1)

	for sc.Scan() {
		text := sc.Text()
		item, err := r.processLine(tempMap, text)
		if err != nil {
			log.Warningln(err)
		} else {
			processLine := helper.IsAfter(item.Begin, start) && helper.IsBefore(item.End, end)
			if item.Success && processLine {
				items = append(items, item)
				processedLine++
			} else if !processLine {
				skipLines++
			}
			bar.Add(1)
		}
	}

	if len(items) == 0 {
		log.Debugln("No data found... :(")
	}

	if err := sc.Err(); err != nil {
		log.Errorln(err)
	}

	return &items, processedLine, skipLines, nil
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
		if helper.CheckPrivateNetwork(ip) {
			return structs.ImportItem{Success: false}, nil
		}

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
