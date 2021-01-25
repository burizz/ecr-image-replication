package config

import (
	"fmt"

	"github.com/docker/docker/client"
)

// DockerClientInit() - returns an initialized Docker client
func DockerClientInit() (dockerClient *client.Client, dockerInitErr error) {
	dockerClient, initDockerClientErr := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if initDockerClientErr != nil {
		return nil, fmt.Errorf("cannot inititialize docker client: %v", initDockerClientErr)
	}

	return dockerClient, nil
}
