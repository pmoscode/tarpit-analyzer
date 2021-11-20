package modules

import (
	"endlessh-analyzer/analyze"
	"endlessh-analyzer/cli"
)

type AnalyzeCmd struct {
	BatchSize int `default:"50" help:"Query batch size for GEO IP API (if supported by GEO IP API)."`
}

func (r *AnalyzeCmd) Run(ctx *cli.Context) error {
	return analyze.DoAnalyze(ctx)
}
