package importData

import (
	"endlessh-analyzer/api"
	cachedb "endlessh-analyzer/cache"
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database"
	"endlessh-analyzer/database/schemas"
	"endlessh-analyzer/importData/modules"
	"endlessh-analyzer/importData/structs"
	log "github.com/sirupsen/logrus"
)

type Import interface {
	Import(sourcePath string, context *cli.Context) (*[]structs.ImportItem, error)
}

type ImportSource int

const (
	Endlessh ImportSource = iota
)

func createImportSource(source ImportSource) Import {
	switch source {
	case Endlessh:
		return modules.Endlessh{}
	}

	return nil
}

func DoImport(source ImportSource, sourcePath string, batchSize int, context *cli.Context) error {
	if context.Target != "" {
		log.Infoln("--target was set to '" + context.Target + "', but is actually unused for import command...")
	}

	importAction := createImportSource(source)

	importItems, errCreate := importAction.Import(sourcePath, context)
	if errCreate != nil {
		return errCreate
	}

	cachedb.Init(context.Debug)
	db, errCreate := database.CreateDbData(context.Debug)
	if errCreate != nil {
		log.Panicln("Data database could not be loaded.", errCreate)
	}

	result, errSave := db.SaveData(db.Map(importItems, db.MapToData))
	if errSave != nil {
		return errSave
	}

	if result == database.DbOk {
		log.Debugln("Imported data saved to database")
	}

	rows, errQuery := db.DbRawQuery(schemas.Location{}, getQueryParametersUnlocalizedIps())
	if errQuery != nil {
		return errQuery
	}

	ips := make([]string, 0)
	for rows.Next() {
		var ip string
		rows.Scan(&ip)
		ips = append(ips, ip)
	}
	defer rows.Close()

	err := processIps(ips, batchSize)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func processIps(ips []string, batchSize int) error {
	if len(ips) == 0 {
		return nil
	}

	batchCount := len(ips)
	geolocationApi := api.CreateGeoLocationAPI(api.IpApiCom)
	for i := 0; i < batchCount; i += batchSize {
		if i+batchSize >= batchCount {
			batchSize = batchCount - i
		}

		ipBatch := ips[i : i+batchSize]
		resolved, _ := geolocationApi.QueryGeoLocationAPI(ipBatch)
		err := cachedb.SaveLocations(resolved)
		if err != nil {
			return err
		}
	}

	return nil
}
