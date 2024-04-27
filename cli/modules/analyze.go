package modules

import (
	"tarpit-analyzer/analyze"
	"tarpit-analyzer/cli"
)

type AnalyzeCmd struct{}

func (r *AnalyzeCmd) Run(ctx *cli.Context) error {
	return analyze.DoAnalyze(ctx)
}
