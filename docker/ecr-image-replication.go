package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/burizz/ecr-image-replication/config"
)

func main() {
	logFormat := os.Getenv("LOG_FORMAT")
	logOutput := os.Getenv("LOG_OUTPUT")
	//logLevel := os.Getenv("LOG_LEVEL")
	logLevel := "DEBUG"

	config.LoggingConfig(logFormat, logLevel, logOutput)

	log.Debug("End")
}
