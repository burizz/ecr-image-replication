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
	var dockerImages = map[string]string{
		"hello-world": "latest",
		"alpine":      "3.12.3",
	}
	var ecrRegistry = "235694435776.dkr.ecr.us-east-1.amazonaws.com"
	var awsRegion = "us-east-1"

	ecrAuthToken, getEcrTokenErr := aws.GetECRAuthToken(awsRegion)
	if getEcrTokenErr != nil {
		log.Errorf("Error: %v", getEcrTokenErr)
	}

	for image, tag := range dockerImages {
		dockerImageConfig := docker.Image{Image: image, Tag: tag, AuthToken: ecrAuthToken}

		// Pull images
		log.Infof("Pulling image [%v]", image)
		if imagePullErr := dockerImageConfig.PullImage(); imagePullErr != nil {
			log.Errorf("Error: %v", imagePullErr)
		}

		ecrImageTag := ecrRegistry + "/" + image + ":" + tag

		// Change image tag
		imageTagErr := dockerImageConfig.TagImage(ecrImageTag)
		if imageTagErr != nil {
			log.Errorf("Error: %v", imageTagErr)
		}

		// Create ECR repo if it doesn't exist; skip if it does
		if ecrRepoCreateErr := aws.CreateECRRepo(awsRegion, image); ecrRepoCreateErr != nil {
			log.Errorf("Error: %v", ecrRepoCreateErr)
		}

		// Push image to ECR
		if pushImageErr := dockerImageConfig.PushImage(ecrImageTag); pushImageErr != nil {
			log.Errorf("Error: %v", pushImageErr)
		}

		// Cleanup local image after push
		if removeImageErr := dockerImageConfig.RemoveImage(ecrImageTag); removeImageErr != nil {
			log.Errorf("Error: %v", removeImageErr)
		}
	}

	// List local images
	if listImageErr := docker.ListImages(); listImageErr != nil {
		log.Errorf("Error: %v", listImageErr)
	}
}
