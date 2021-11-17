package importData

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/importData/modules"
	"endlessh-analyzer/importData/structs"
	log "github.com/sirupsen/logrus"
)

type ImportData interface {
	Import(sourcePath string, context *cli.Context) (*[]structs.ImportItem, error)
}

type ImportSource int

const (
	Endlessh ImportSource = iota
)

func createImportSource(source ImportSource) ImportData {
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

	importItems, err := importAction.Import(sourcePath, context)
	if err != nil {
		return err
	}

	log.Debug(len(*importItems))

	// Save to DB

	return nil
}
