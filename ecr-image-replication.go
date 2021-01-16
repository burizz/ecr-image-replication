package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/burizz/ecr-image-replication/config"
	"github.com/burizz/ecr-image-replication/docker"
)

func init() {
	logFormat := os.Getenv("LOG_FORMAT")
	logOutput := os.Getenv("LOG_OUTPUT")
	logLevel := os.Getenv("LOG_LEVEL")
	config.LoggingConfig(logFormat, logLevel, logOutput)
}

func main() {
	// TODO: Get Docker images from env variable passed into Docker
	//dockerImages = os.Getenv("DOCKER_IMAGES")
	var dockerImages = []string{"alpine:latest", "scratch", "alpine:3.12.3"}
	for _, image := range dockerImages {
		log.Infof("Test pull Image: %s", image)
		docker.PullImage(image, "")
	}
}
