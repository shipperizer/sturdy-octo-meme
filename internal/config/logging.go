package config

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// SetupLogger initializes the log level of the global logger
func SetupLogger(lvl string) {
	level, err := log.ParseLevel(lvl)

	if err != nil {
		level = log.ErrorLevel
	}
	// set global log level
	log.SetLevel(level)

	// Set Log type
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "@timestamp",
			log.FieldKeyLevel: "severity",
			log.FieldKeyMsg:   "message",
		},
	})
}
