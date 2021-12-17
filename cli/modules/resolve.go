package modules

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/resolve"
)

type ResolveCmd struct{}

func (r *ResolveCmd) Run(ctx *cli.Context) error {
	return resolve.DoResolve(ctx)
}
