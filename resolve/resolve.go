package resolve

import (
	"database/sql"
	cachedb "endlessh-analyzer/cache"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database"
	log "github.com/sirupsen/logrus"
)

func DoResolve(context *cli.Context) error {
	cachedb.Init(context.Debug)
	db, errCreate := database.CreateDbData(context.Debug)
	if errCreate != nil {
		log.Panicln("Data database could not be loaded.", errCreate)
	}

	rows, errQuery := db.DbRawQuery(getQueryParametersUnlocalizedIps())
	if errQuery != nil {
		return errQuery
	}

	ips := make([]string, 0)
	for rows.Next() {
		var ip string
		err := rows.Scan(&ip)
		if err != nil {
			return err
		}
		ips = append(ips, ip)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Errorln("Could not close query in DB.")
		}
	}(rows)

	log.Infoln("####### [Start] Processing IP's #######")
	processIps(ips)
	log.Infoln("####### [End] Processing IP's #######")

	return nil
}

func processIps(ips []string) {
	if len(ips) == 0 {
		log.Infoln("No IP's to process...")
		return
	}

	processedIps := cachedb.ResolveLocationsFor(ips)
	log.Infoln("Processed IP's: ", processedIps)
}
