package settings

import (
	log "github.com/sirupsen/logrus"
)

var Log *log.Logger

func init() {
	Log = log.New()
	Log.SetLevel(log.TraceLevel)
}
