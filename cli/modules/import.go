package modules

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/importData"
	"errors"
)

type ImportCmd struct {
	FileSource string `arg:"" default:"tarpit.log" help:"Log file to import." type:"path"`
	Type       string `default:"endlessh" enum:"endlessh,tarpit" help:"For now only endlessh is implemented"`
	BatchSize  int    `default:"50" help:"Amount of ips which are send in one batch to GeoLocalization API. (If supported in API)"`
}

func (r *ImportCmd) Run(ctx *cli.Context) error {
	switch r.Type {
	case "endlessh":
		return importData.DoImport(importData.Endlessh, r.FileSource, r.BatchSize, ctx)
	case "sshTarpit":
		return importData.DoImport(importData.SshTarpit, r.FileSource, r.BatchSize, ctx)
	default:
		return errors.New("import type '" + r.Type + "' not implemented")
	}

}
