package main

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/cli/modules"
	"github.com/alecthomas/kong"
)

var cliStruct struct {
	Debug  bool   `short:"d" default:"false" help:"Enable debug mode."`
	Target string `short:"t" help:"filename where output should be saved" type:"path"`

	Import  modules.ImportCmd  `cmd:"" help:"ImportCmd logs from different tarpit apps."`
	Analyze modules.AnalyzeCmd `cmd:"" help:"Analyze file."`
	Export  modules.ExportCmd  `cmd:"" help:"Export to different formats"`
}

func main() {
	ctx := kong.Parse(&cliStruct)

	err := ctx.Run(&cli.Context{Debug: cliStruct.Debug, Target: cliStruct.Target})
	ctx.FatalIfErrorf(err)
}
