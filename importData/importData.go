package importData

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/database"
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

func DoImport(source ImportSource, sourcePath string, context *cli.Context) error {
	if context.Target != "" {
		log.Infoln("--target was set to '" + context.Target + "', but is actually unused for import command...")
	}

	importAction := createImportSource(source)

	importItems, errCreate := importAction.Import(sourcePath, context)
	if errCreate != nil {
		return errCreate
	}

	db, errCreate := database.CreateDbData()
	if errCreate != nil {
		log.Panicln("Cache database could not be loaded.", errCreate)
	}

	result, errSave := db.SaveData(db.Map(importItems, db.MapToData))
	if errSave != nil {
		return errSave
	}

	if result == database.DbOk {
		log.Debugln("Imported data saved to database")
	}

	return nil
}
