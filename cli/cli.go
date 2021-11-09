package cli

import (
	"endlessh-analyzer/analyze"
	"endlessh-analyzer/convert"
	geolocation "endlessh-analyzer/geo-location"
)

type Context struct {
	Debug bool
}

type ConvertCmd struct {
	FileSource string `arg:"" default:"tarpit.log" help:"Log file to analyze." type:"path"`
	FileTarget string `arg:"" default:"tarpit-converted.log" help:"Write converted data to" type:"path"`

	StartDate string `default:"unset" help:"Only consider data starting at <yyyy-mm-dd>"`
	EndDate   string `default:"unset" help:"Only consider data ending at <yyyy-mm-dd> (including that day)"`
}

type AnalyzeCmd struct {
	FileSource string `arg:"" default:"tarpit-converted.log" help:"Converted log file to analyze." type:"path"`
	FileTarget string `arg:"" default:"tarpit-analyzed.txt" help:"Write analyzed data to" type:"path"`
}

type GeoCmd struct {
	FileSource string `arg:"" default:"tarpit-converted.log" help:"Converted log file to get Geo-Locations for." type:"path"`
	FileTarget string `arg:"" default:"tarpit-geo.kml" help:"Write Geo-Location data to (ex. for Google Maps)" type:"path"`

	BatchSize                  int    `default:"50" help:"Query batch size for GEO IP API."`
	CenterGeoLocationLatitude  string `default:"50.840886980084086" help:"Latitude you wish to be the target on the map. Default: Germany"`
	CenterGeoLocationLongitude string `default:"10.276290870120306" help:"Longitude you wish to be the target on the map. Default: Germany"`
}

func (r *ConvertCmd) Run(ctx *Context) error {
	return convert.DoConvert(r.FileSource, r.FileTarget, r.StartDate, r.EndDate, ctx.Debug)
}

func (r *AnalyzeCmd) Run(ctx *Context) error {
	return analyze.DoAnalyze(r.FileSource, r.FileTarget, ctx.Debug)
}

func (r *GeoCmd) Run(ctx *Context) error {
	return geolocation.DoLocalization(r.FileSource, r.FileTarget, r.BatchSize, r.CenterGeoLocationLongitude, r.CenterGeoLocationLatitude,
		ctx.Debug)
}
