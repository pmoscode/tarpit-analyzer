package main

import (
	"endlessh-analyzer/cli"
	"github.com/alecthomas/kong"
)

var cliStruct struct {
	Debug bool `default:"false" help:"Enable debug mode."`

	Convert cli.ConvertCmd `cmd:"" help:"Convert file to generic format. Only already closed SSH sessions are processed!"`

	Analyze cli.AnalyzeCmd `cmd:"" help:"Analyze file."`

	Geo cli.GeoCmd `cmd:"" help:"Generate KML file for visualization."`
}

func main() {
	ctx := kong.Parse(&cliStruct)

	err := ctx.Run(&cli.Context{Debug: cliStruct.Debug})
	ctx.FatalIfErrorf(err)
}
