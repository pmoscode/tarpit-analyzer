
![Logo](banner.png)


[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/GPL-3.0)

# Tarpit Analyzer

If you have an ssh tarpit service running on you own, and you want somehow to analyze the logged data...

With Tarpit Analyzer you can dig into the data and do some analysis and generate visual outputs which you can then import in Google Maps or Openstreetmap.

Currently, supported tarpits:

- Endlessh: https://github.com/skeeto/endlessh
- Python Ssh-tarpit: https://pypi.org/project/ssh-tarpit/

## Features

- Kong cli parser: https://github.com/alecthomas/kong
- Progressbar: https://github.com/schollz/progressbar
- Duration format: https://golangexample.com/better-time-duration-formatting-in-go/


## Installation

Download binary from https://gitlab.com/pmoscode/tarpit-analyzer/-/releases for your arch.
Or clone this repository and build on your own.
    
## Usage/Examples

```shell
./endlessh_analyzer import <path-to>/endlessh.log --type=endlessh # Import Enlessh logs
./endlessh_analyzer analyze --target=analyze.txt # Generate analysis
./endlessh_analyzer export json --start-date=2021-07-16 --end-date=2021-07-18 --target=export.json # # Exports for a given data range to json format
```
