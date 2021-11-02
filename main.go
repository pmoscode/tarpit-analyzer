package main

import (
	"endlessh-analyzer/cli"
	"github.com/alecthomas/kong"
)

var cliStruct struct {
	Debug bool `default:"false" help:"Enable debug mode."`

	Convert cli.ConvertCmd `cmd:"" help:"Convert file to generic format."`
}

func main() {
	ctx := kong.Parse(&cliStruct)

	err := ctx.Run(&cli.Context{Debug: cliStruct.Debug})
	ctx.FatalIfErrorf(err)
}
