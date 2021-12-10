package statistics

import (
	"database/sql"
	"endlessh-analyzer/database"
	"endlessh-analyzer/database/schemas"
	"fmt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
	"strings"
	time2 "time"
)

type TimeStatistic struct {
	label  string
	format string
}

var (
	DAY = TimeStatistic{
		label:  "Weekday",
		format: "%w",
	}
	MONTH = TimeStatistic{
		label:  "Month",
		format: "%m %Y",
	}
	YEAR = TimeStatistic{
		label:  "Year",
		format: "%Y",
	}
)

var weekdays = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var months = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

func GetAttackTimeStatistics(db *database.DbData, start *time2.Time, end *time2.Time, mode TimeStatistic) (string, error) {
	rowsTop, errQuery := db.DbQuery(schemas.Data{}, getQueryParametersDateAttacks(start, end, mode))
	if errQuery != nil {
		return "", errQuery
	}

	dataRows := make([][]string, 0)
	for rowsTop.Next() {
		var modeStr string
		var attacks int64

		err := rowsTop.Scan(&modeStr, &attacks)

		if mode == DAY {
			val, _ := strconv.Atoi(modeStr)
			modeStr = weekdays[val]
		}

		if mode == MONTH {
			split := strings.Split(modeStr, " ")
			val, _ := strconv.Atoi(split[0])
			modeStr = months[val-1] + " " + split[1]
		}

		if err != nil {
			return "", err
		}
		dataRows = append(dataRows, []string{modeStr, fmt.Sprint(attacks)})
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Errorln("Could not close query in DB.")
		}
	}(rowsTop)

	builder := new(strings.Builder)

	table := tablewriter.NewWriter(builder)
	table.SetHeader([]string{mode.label, "Attacks"})
	table.SetAutoWrapText(false)

	p := message.NewPrinter(language.English)
	for _, item := range dataRows {
		line := []string{item[0], p.Sprint(item[1])}
		table.Append(line)
	}

	table.Render()

	return "  " + strings.ToUpper(mode.label) + " ATTACKER STATISTICS\n" + builder.String(), nil
}

//  DAY:
//SELECT strftime('%w', begin) as weekday, count(id)
//FROM data
//where success == '1'
//GROUP BY weekday

//  Month:
//SELECT strftime('%m', begin) as month, count(id)
//FROM data
//where success == '1'
//GROUP BY month

//  Year:
//SELECT strftime('%Y', begin) as year, count(id)
//FROM data
//where success == '1'
//GROUP BY year
