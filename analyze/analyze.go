package analyze

import log "github.com/sirupsen/logrus"

var debug = false

func DoAnalyze(pathSource string, debugParam bool) error {
	debug = debugParam
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	return nil
}
