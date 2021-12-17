package importData

import (
	"database/sql"
	cachedb "endlessh-analyzer/cache"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database"
	"endlessh-analyzer/helper"
	"endlessh-analyzer/importData/modules"
	"endlessh-analyzer/importData/structs"
	"fmt"
	log "github.com/sirupsen/logrus"
	time2 "time"
)

type Import interface {
	Import(sourcePath string, start *time2.Time, end *time2.Time, context *cli.Context) (*[]structs.ImportItem, int, int, error) // int, int == processed and skipped lines
}

type ImportSource int

const (
	Endlessh ImportSource = iota
	SshTarpit
)

func createImportSource(source ImportSource) Import {
	switch source {
	case Endlessh:
		return modules.Endlessh{}
	case SshTarpit:
		return modules.SshTarpit{}
	}

	return nil
}

func DoImport(source ImportSource, sourcePath string, skipIpResolving bool, context *cli.Context) error {
	if context.Target != "" {
		log.Infoln("--target was set to '" + context.Target + "', but is actually unused for import command...")
	}

	start := helper.GetDate(context.StartDate)
	end := helper.GetDate(context.EndDate)

	if start != nil {
		log.Infoln("Starting at date: ", context.StartDate)
	}
	if end != nil {
		log.Infoln("Stopping at date: ", context.EndDate)
	}

	importAction := createImportSource(source)

	log.Infoln("####### [Start] Reading file #######")
	importItems, processedLines, skippedLines, errCreate := importAction.Import(sourcePath, start, end, context)
	if errCreate != nil {
		return errCreate
	}
	fmt.Println()
	log.Infoln("Processed lines: ", processedLines)
	log.Infoln("Skipped lines: ", skippedLines)
	log.Infoln("####### [End] Reading file #######")

	db, errCreate := database.CreateDbData(context.Debug)
	if errCreate != nil {
		log.Panicln("Data database could not be loaded.", errCreate)
	}

	log.Infoln("####### [Start] Saving to database #######")
	result, errSave := db.SaveData(db.Map(importItems, db.MapToData))
	if errSave != nil {
		return errSave
	}
	log.Infoln("####### [End] Saving to database #######")

	if result == database.DbOk {
		log.Debugln("Imported data saved to database")
	}

	cachedb.Init(context.Debug)
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

	if !skipIpResolving {
		log.Infoln("####### [Start] Processing IP's #######")
		processIps(ips)
		log.Infoln("####### [End] Processing IP's #######")
	} else {
		log.Infoln("####### IP resolving skipped #######")
	}

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
