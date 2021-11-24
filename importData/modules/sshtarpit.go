package modules

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/importData/structs"
	"errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SshTarpit struct {
	debug bool
}

func (r SshTarpit) Import(sourcePath string, context *cli.Context) (*[]structs.ImportItem, error) {
	r.debug = context.Debug
	if r.debug {
		log.SetLevel(log.DebugLevel)
	}

	items := make([]structs.ImportItem, 0)

	return &items, nil
}

func (r SshTarpit) processLine(line string) (structs.ImportItem, error) {
	return structs.ImportItem{Success: false}, nil
}

func (r SshTarpit) getValue(source string) (string, error) {
	if source == "" {
		return "", errors.New("source is empty")
	}

	split := strings.Split(source, "=")

	if len(split) == 1 {
		return "", errors.New("source has no '=' char")
	}

	return split[1], nil
}
