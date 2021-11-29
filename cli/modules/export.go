package modules

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/export"
	log "github.com/sirupsen/logrus"
)

type ExportCmd struct {
	Kml  ExportKmlCmd  `cmd:"" help:"Export KML"`
	Csv  ExportCsvCmd  `cmd:"" help:"Export CSV"`
	Json ExportJsonCmd `cmd:"" help:"Export Json"`
}

type ExportCsvCmd struct {
	Separator string `default:"," help:"Separator to use as delimiter."`
}

type ExportJsonCmd struct{}

type ExportKmlCmd struct {
	CenterGeoLocationLatitude  string `default:"50.840886980084086" help:"Latitude you wish to be the target on the map. Default: Germany"`
	CenterGeoLocationLongitude string `default:"10.276290870120306" help:"Longitude you wish to be the target on the map. Default: Germany"`
}

func (r *ExportCmd) Run(ctx *cli.Context) error {
	log.Infoln("Export command is not implemented. Use the subcommands... See export --help for more information")

	return nil
}

func (r *ExportCsvCmd) Run(ctx *cli.Context) error {
	params := export.Parameters{
		Separator: r.Separator,
	}

	return export.DoExport(export.CSV, params, ctx)
}

func (r *ExportJsonCmd) Run(ctx *cli.Context) error {
	params := export.Parameters{}

	return export.DoExport(export.JSON, params, ctx)
}

func (r *ExportKmlCmd) Run(ctx *cli.Context) error {
	params := export.Parameters{
		CenterGeoLocationLatitude:  r.CenterGeoLocationLatitude,
		CenterGeoLocationLongitude: r.CenterGeoLocationLongitude,
	}

	return export.DoExport(export.KML, params, ctx)
}
