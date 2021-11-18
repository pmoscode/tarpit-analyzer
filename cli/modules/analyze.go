package modules

import (
	"endlessh-analyzer/analyze"
	"endlessh-analyzer/cli"
)

type AnalyzeCmd struct {
	BatchSize int `default:"50" help:"Query batch size for GEO IP API (if supported by GEO IP API)."`

	StartDate string `default:"unset" help:"Only consider data starting at <yyyy-mm-dd>"`
	EndDate   string `default:"unset" help:"Only consider data ending at <yyyy-mm-dd> (including that day)"`
}

func (r *AnalyzeCmd) Run(ctx *cli.Context) error {
	return analyze.DoAnalyze(r.StartDate, r.EndDate, ctx)
}
