package config

import (
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

//TODO: document
func DockerClientInit() (dockerClient *client.Client, dockerInitErr error) {
	dockerClient, initDockerClientErr := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if initDockerClientErr != nil {
		log.Errorf("Cannot initialize Docker client: %v", initDockerClientErr)
		return nil, initDockerClientErr
	}

	return dockerClient, nil
}
