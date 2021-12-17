package modules

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/importData"
	"errors"
)

type ImportCmd struct {
	FileSource      string `arg:"" default:"tarpit.log" help:"Log file to import." type:"path"`
	Type            string `default:"endlessh" enum:"endlessh,sshTarpit" help:"Import logs from 'endlessh' or 'sshTarpit'"`
	SkipIpResolving bool   `default:"false" help:"Skip the resolving of IP's during the import."`
}

func (r *ImportCmd) Run(ctx *cli.Context) error {
	switch r.Type {
	case "endlessh":
		return importData.DoImport(importData.Endlessh, r.FileSource, r.SkipIpResolving, ctx)
	case "sshTarpit":
		return importData.DoImport(importData.SshTarpit, r.FileSource, r.SkipIpResolving, ctx)
	default:
		return errors.New("import type '" + r.Type + "' not implemented")
	}

}
