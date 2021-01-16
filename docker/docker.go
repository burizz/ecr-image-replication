package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"

	"github.com/burizz/ecr-image-replication/config"
)

func PullImage(imageName string, imageTag string) (imagePullErr error) {
	dockerCtx := context.Background()

	dockerClient, initClientErr := config.DockerClientInit()
	if initClientErr != nil {
		//TODO: Fix proper exit if cannot init Docker client
		log.Fatalf("Err")
	}

	// TODO: move to docker package
	reader, dockerPullErr := dockerClient.ImagePull(dockerCtx, "apline", types.ImagePullOptions{})
	if dockerPullErr != nil {
		log.Errorf("Err: Cannot download image: %v", dockerPullErr)
	}
	// TODO: figure this out
	io.Copy(os.Stdout, reader)
}

func ListImages() error {
	dockerCtx := context.Background()

	dockerClient, initClientErr := config.DockerClientInit()
	if initClientErr != nil {
		//TODO: Fix proper exit if cannot init Docker client
		log.Fatalf("Err")
	}

	images, listImageErr := dockerClient.ImageList(dockerCtx, types.ImageListOptions{})
	if listImageErr != nil {
		log.Errorf("Err: Cannot list images: %v", listImageErr)
	}
	for _, image := range images {
		fmt.Println(image.ID)
	}
}
