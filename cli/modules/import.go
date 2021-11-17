package modules

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/importData"
	"errors"
)

type ImportCmd struct {
	FileSource string `arg:"" default:"tarpit.log" help:"Log file to import." type:"path"`
	Type       string `default:"endlessh" enum:"endlessh,tarpit" help:"For now only endlessh is implemented"`
}

func (r *ImportCmd) Run(ctx *cli.Context) error {
	switch r.Type {
	case "endlessh":
		return importData.DoImport(importData.Endlessh, r.FileSource, ctx)
	default:
		return errors.New("import type '" + r.Type + "' not implemented")
	}

}
