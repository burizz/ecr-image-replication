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
	var dockerImages = []string{"hello-world:latest"}
	var ecrRegistry = ""

	// Pull images
	for _, imageTag := range dockerImages {
		log.Infof("Pulling image [%v]", imageTag)
		if imagePullErr := docker.PullImage(imageTag); imagePullErr != nil {
			log.Errorf("Error: %v", imagePullErr)
		}

		ok, imageTagErr := docker.TagImage(imageTag, ecrRegistry)
		if !ok || imageTagErr != nil {
			log.Errorf("Error: %v", imageTagErr)
		}
	}

	// List local images
	if listImageErr := docker.ListImages(); listImageErr != nil {
		log.Errorf("Error: %v", listImageErr)
	}

	// TODO: cleanup local images
}
