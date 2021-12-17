package modules

import (
	"endlessh-analyzer/analyze"
	"endlessh-analyzer/cli"
)

type AnalyzeCmd struct{}

func (r *AnalyzeCmd) Run(ctx *cli.Context) error {
	return analyze.DoAnalyze(ctx)
}
