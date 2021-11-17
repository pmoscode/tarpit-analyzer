package modules

import (
	"endlessh-analyzer/analyze"
	"endlessh-analyzer/cli"
)

type AnalyzeCmd struct {
	FileSource string `arg:"" default:"tarpit-converted.log" help:"Converted log file to analyze." type:"path"`
	FileTarget string `arg:"" default:"tarpit-analyzed.txt" help:"Write analyzed data to" type:"path"`
	BatchSize  int    `default:"50" help:"Query batch size for GEO IP API (if supported by GEO IP API)."`
}

func (r *AnalyzeCmd) Run(ctx *cli.Context) error {
	return analyze.DoAnalyze(r.FileSource, r.FileTarget, ctx.Debug)
}
