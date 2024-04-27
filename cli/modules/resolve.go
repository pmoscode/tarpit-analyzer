package modules

import (
	"tarpit-analyzer/cli"
	"tarpit-analyzer/resolve"
)

type ResolveCmd struct{}

func (r *ResolveCmd) Run(ctx *cli.Context) error {
	return resolve.DoResolve(ctx)
}
