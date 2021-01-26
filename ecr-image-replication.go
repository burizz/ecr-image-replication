package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/burizz/ecr-image-replication/aws"
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
	var ecrRegistry = "235694435776.dkr.ecr.us-east-1.amazonaws.com/image-replication"

	if loginToEcrErr := aws.LoginToECR(); loginToEcrErr != nil {
		log.Errorf("Error: %v", loginToEcrErr)
	}

	for _, imageTag := range dockerImages {
		// Pull images
		log.Infof("Pulling image [%v]", imageTag)
		if imagePullErr := docker.PullImage(imageTag); imagePullErr != nil {
			log.Errorf("Error: %v", imagePullErr)
		}

		ecrImageTag := ecrRegistry + "/" + imageTag

		// Change image tag
		ok, imageTagErr := docker.TagImage(imageTag, ecrImageTag)
		if !ok || imageTagErr != nil {
			log.Errorf("Error: %v", imageTagErr)
		}

		// Push image to ECR
		if pushImageErr := docker.PushImage(ecrImageTag); pushImageErr != nil {
			log.Errorf("Error: %v", pushImageErr)
		}
	}

	// List local images
	if listImageErr := docker.ListImages(); listImageErr != nil {
		log.Errorf("Error: %v", listImageErr)
	}

	// TODO: cleanup local images
}
