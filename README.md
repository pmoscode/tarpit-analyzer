![Logo](banner.png)

[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/GPL-3.0)

# Tarpit Analyzer

If you have an ssh tarpit service running on you own, and you want somehow to analyze the logged data...

With Tarpit Analyzer you can dig into the data and do some analysis and generate visual outputs which you can then
import in Google Maps or Openstreetmap.

Currently, supported tarpits:

- Endlessh: https://github.com/skeeto/endlessh
- Python Ssh-tarpit: https://pypi.org/project/ssh-tarpit/

## Features

- Import logs from two sources (Endlessh and ssh-tarpit) into internal database
- Can run an analysis on selected date range of data
- Outputs different visualizations of analyzed data:
- lines (KML and GeoJson format)
- place marks (on attacker country with number of attacks)

## Endpoints used

To get the location of the ips, Tarpit-Analyzer uses following apis:

- https://ip-api.com/
- https://reallyfreegeoip.org/
- https://ipapi.co/
- https://www.geoplugin.com/webservices/json

All endpoint can be used without a token / login.

## Installation

Download binary from https://gitlab.com/pmoscode/tarpit-analyzer/-/releases for your arch. Or clone this repository and
build on your own.

## CLI documentation

### Commands

#### import

    endlessh_analyzer import [<file-source>]

"file-source" is optional. By default, it is set to "tarpit.log"

Flags:

|              | Default  | Description                                                                             |
|--------------|----------|-----------------------------------------------------------------------------------------|
| --type       | endlessh | Import logs from 'endlessh' or 'sshTarpit'                                              |

#### analyze

    endlessh_analyzer analyze 

uses --target flag as destination output

#### export

    endlessh_analyzer export <subcommand>

##### Subcommand: kml
Flags:

|                                 | Default            | Description                                                      |
|---------------------------------|--------------------|------------------------------------------------------------------|
| --center-geo-location-latitude  | 50.840886980084086 | Latitude you wish to be the target on the map. Default: Germany  |
| --center-geo-location-longitude | 10.276290870120306 | Longitude you wish to be the target on the map. Default: Germany |

##### Subcommand: geojson
Flags:

|                                 | Default            | Description                                                                                                                                                           |
|---------------------------------|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| --type                          | point              | 'line': Creates line from attacker source to CenterGeoLocation. <br/> 'point': Places point on attacker country with sum of attacks (prefer for large amount of data) |
| --center-geo-location-latitude  | 50.840886980084086 | Latitude you wish to be the target on the map (for 'line' type). Default: Germany                                                                                     |
| --center-geo-location-longitude | 10.276290870120306 | Longitude you wish to be the target on the map (for 'line' type). Default: Germany                                                                                    |

##### Subcommand: csv
Flags:

|             | Default | Description                   |
|-------------|---------|-------------------------------|
| --separator | ,       | Separator to use as delimiter |

##### Subcommand: json

No special flags.

### General flags

| Short | Long         | Default | Description                                 |
|-------|--------------|---------|---------------------------------------------|
| -h    | --help       |         | Show context-sensitive help.                |
| -d    | --debug      |         | Enable debug mode.                          |
| -t    | --target     | unset   | filename where output should be saved       |
|       | --start-date | unset   | Only consider data starting at <yyyy-mm-dd> |
|       | --end-date   | unset   | Only consider data ending at <yyyy-mm-dd>   |

## Usage/Examples

```shell
./endlessh_analyzer import <path-to>/endlessh.log --type=endlessh # Import Endlessh logs
./endlessh_analyzer analyze --target=analyze.txt # Generate analysis
./endlessh_analyzer export json --start-date=2021-07-16 --end-date=2021-07-18 --target=export.json # Exports a given data range to json format
```

## Known issues

Tests missing O_o
