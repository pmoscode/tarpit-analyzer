package cli

import (
	"endlessh-analyzer/analyze"
	"endlessh-analyzer/convert"
)

type Context struct {
	Debug bool
}

type ConvertCmd struct {
	FileSource string `arg:"" default:"tarpit.log" help:"Log file to analyze." type:"path"`
	FileTarget string `arg:"" default:"tarpit-converted.log" help:"Write converted data to" type:"path"`

	StartDate string `default:"unset" help:"Only consider data starting at <yyyy-mm-dd>"`
	EndDate   string `default:"unset" help:"Only consider data ending at <yyyy-mm-dd> (including that day)"`
}

type AnalyzeCmd struct {
	FileSource string `arg:"" default:"tarpit-converted.log" help:"Converted log file to analyze." type:"path"`
	FileTarget string `arg:"" default:"tarpit-analyzed.txt" help:"Write analyzed data to" type:"path"`
}

func (r *ConvertCmd) Run(ctx *Context) error {
	return convert.DoConvert(r.FileSource, r.FileTarget, r.StartDate, r.EndDate, ctx.Debug)
}

func (r *AnalyzeCmd) Run(ctx *Context) error {
	return analyze.DoAnalyze(r.FileSource, r.FileTarget, ctx.Debug)
}
