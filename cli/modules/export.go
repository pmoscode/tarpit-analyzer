package modules

import "endlessh-analyzer/cli"

type ExportCmd struct {
	Kml  ExportKmlCmd  `cmd:"" help:"Test KML"`
	Csv  ExportCsvCmd  `cmd:"" help:"Test CSV"`
	Json ExportJsonCmd `cmd:"" help:"Test Json"`
}

type ExportCsvCmd struct {
	FileTarget string `arg:"" required:"" help:"Write Geo-Location data to (ex. for Google Maps)" type:"path"`
	Separator  string `default:"," help:"Seperator to use as delimiter."`
}

type ExportJsonCmd struct {
	FileTarget string `arg:"" required:"" help:"Write Geo-Location data to (ex. for Google Maps)" type:"path"`
}

type ExportKmlCmd struct {
	FileTarget                 string `arg:"" required:"" help:"Write Geo-Location data to (ex. for Google Maps)" type:"path"`
	CenterGeoLocationLatitude  string `default:"50.840886980084086" help:"Latitude you wish to be the target on the map. Default: Germany"`
	CenterGeoLocationLongitude string `default:"10.276290870120306" help:"Longitude you wish to be the target on the map. Default: Germany"`
}

func (r *ExportCmd) Run(ctx *cli.Context) error {
	return nil
}

func (r *ExportCsvCmd) Run(ctx *cli.Context) error {
	return nil
}

func (r *ExportJsonCmd) Run(ctx *cli.Context) error {
	return nil
}

func (r *ExportKmlCmd) Run(ctx *cli.Context) error {
	return nil
}
