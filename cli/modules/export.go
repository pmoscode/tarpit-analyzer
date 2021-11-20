package modules

import (
	"endlessh-analyzer/cli"
	"endlessh-analyzer/export"
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

//func (r *ExportCmd) Run(ctx *cli.Context) error {
//	return nil
//}

func (r *ExportCsvCmd) Run(ctx *cli.Context) error {
	return export.CSV(r.Separator, ctx)
}

func (r *ExportJsonCmd) Run(ctx *cli.Context) error {
	return nil
}

func (r *ExportKmlCmd) Run(ctx *cli.Context) error {
	return nil
}
