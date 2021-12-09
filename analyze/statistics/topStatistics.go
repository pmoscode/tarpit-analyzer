package statistics

import (
	"database/sql"
	"endlessh-analyzer/database"
	"endlessh-analyzer/database/schemas"
	"fmt"
	"github.com/hako/durafmt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
	"strings"
	time2 "time"
)

func GetTopStatistics(db *database.DbData, start *time2.Time, end *time2.Time, debug bool) (string, error) {
	rowsTop, errQuery := db.DbQuery(schemas.Data{}, getRawTopCountriesAttacks(start, end))
	if errQuery != nil {
		return "", errQuery
	}

	dataRows := make([][]string, 0)
	for rowsTop.Next() {
		var country string
		var sumAttacks, sumTime, avgTime int64

		err := rowsTop.Scan(&country, &sumAttacks, &sumTime, &avgTime)
		if err != nil {
			return "", err
		}
		dataRows = append(dataRows, []string{country, fmt.Sprint(sumAttacks), fmt.Sprint(sumTime), fmt.Sprint(avgTime)})
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Errorln("Could not close query in DB.")
		}
	}(rowsTop)

	builder := new(strings.Builder)

	table := tablewriter.NewWriter(builder)
	table.SetHeader([]string{"Rank", "Country", "Attacks", "Sum attack time", "Avg attack time"})
	table.SetAutoWrapText(false)

	p := message.NewPrinter(language.English)
	for idx, item := range dataRows {
		sumT, _ := time2.ParseDuration(item[2] + "s")
		avgT, _ := time2.ParseDuration(item[3] + "s")

		sumTFormat, _ := durafmt.ParseString(sumT.String())
		avgTFormat, _ := durafmt.ParseString(avgT.String())

		line := []string{strconv.Itoa(idx + 1), item[0], p.Sprint(item[1]), sumTFormat.String(), avgTFormat.String()}
		table.Append(line)
	}

	table.Render()

	return builder.String(), nil
}
