package modules

import (
	"tarpit-analyzer/cli"
	"tarpit-analyzer/export"
)

type ExportCmd struct {
	Kml     ExportKmlCmd     `cmd:"" help:"Export KML"`
	Geojson ExportGeoJsonCmd `cmd:"" help:"Export GeoJson"`
	Csv     ExportCsvCmd     `cmd:"" help:"Export CSV"`
	Json    ExportJsonCmd    `cmd:"" help:"Export Json"`
}

type ExportCsvCmd struct {
	Separator string `default:"," help:"Separator to use as delimiter."`
}

type ExportJsonCmd struct{}

type ExportKmlCmd struct {
	CenterGeoLocationLatitude  string `default:"50.840886980084086" help:"Latitude you wish to be the target on the map. Default: Germany"`
	CenterGeoLocationLongitude string `default:"10.276290870120306" help:"Longitude you wish to be the target on the map. Default: Germany"`
}

type ExportGeoJsonCmd struct {
	Type                       string `default:"point" enum:"line,point" help:"'line': Creates line from attacker source to CenterGeoLocation. ## 'point': Places point on attacker country with sum of attacks (prefer for large amount of data)"`
	CenterGeoLocationLatitude  string `default:"50.840886980084086" help:"Latitude you wish to be the target on the map (for 'line' type). Default: Germany"`
	CenterGeoLocationLongitude string `default:"10.276290870120306" help:"Longitude you wish to be the target on the map (for 'line' type). Default: Germany"`
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

func (r *ExportGeoJsonCmd) Run(ctx *cli.Context) error {
	params := export.Parameters{
		CenterGeoLocationLatitude:  r.CenterGeoLocationLatitude,
		CenterGeoLocationLongitude: r.CenterGeoLocationLongitude,
		Type:                       r.Type,
	}

	return export.DoExport(export.GEOJSON, params, ctx)
}
