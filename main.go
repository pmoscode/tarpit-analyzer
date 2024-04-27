package main

import (
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
	"os"
	"tarpit-analyzer/cli"
	"tarpit-analyzer/cli/modules"
	"tarpit-analyzer/helper"
)

var cliStruct struct {
	Debug     bool   `short:"d" default:"false" help:"Enable debug mode."`
	Target    string `short:"t" help:"filename where output should be saved" type:"path"`
	StartDate string `default:"unset" help:"Only consider data starting at <yyyy-mm-dd>"`
	EndDate   string `default:"unset" help:"Only consider data ending at <yyyy-mm-dd>"`

	Import  modules.ImportCmd  `cmd:"" help:"ImportCmd logs from different tarpit apps."`
	Resolve modules.ResolveCmd `cmd:"" help:"Resolve all IP's which haven't a corresponding GeoLocation value."`
	Analyze modules.AnalyzeCmd `cmd:"" help:"Analyze file."`
	Export  modules.ExportCmd  `cmd:"" help:"Export to different formats"`
}

func main() {
	file, errLog := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if errLog == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.AddHook(helper.NewConsoleHook(false))

	ctx := kong.Parse(&cliStruct)

	err := ctx.Run(&cli.Context{Debug: cliStruct.Debug, Target: cliStruct.Target, StartDate: cliStruct.StartDate, EndDate: cliStruct.EndDate})
	ctx.FatalIfErrorf(err)
}
