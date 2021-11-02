package cli

import (
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

func (r *ConvertCmd) Run(ctx *Context) error {
	return convert.DoConvert(r.FileSource, r.FileTarget, r.StartDate, r.EndDate, ctx.Debug)
}
