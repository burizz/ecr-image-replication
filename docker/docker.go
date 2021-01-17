package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"

	"github.com/burizz/ecr-image-replication/config"
	log "github.com/sirupsen/logrus"
)

// PullImage - download remote Docker image
func PullImage(imageTag string) (imagePullErr error) {
	dockerCtx := context.Background()

	dockerClient, initClientErr := config.DockerClientInit()
	if initClientErr != nil {
		return initClientErr
	}

	// ImagePull makes an API request to Docker daemon to pull the image, we don't actually pull it and store it ourselves here
	reader, dockerPullErr := dockerClient.ImagePull(dockerCtx, imageTag, types.ImagePullOptions{})
	defer reader.Close()
	if dockerPullErr != nil {
		return fmt.Errorf("cannot download image: %v", dockerPullErr)
	}
	// Sends ImagePull output via reader - to show download progress
	io.Copy(os.Stdout, reader)

	// TODO: convert this output line by line to a logrus event
	//contents, readOutputErr := ioutil.ReadAll(reader)
	//if readOutputErr != nil {
	//return fmt.Errorf("cannot read docker pull output: %v", readOutputErr)
	//}

	log.Infof("Test: %v", string(contents))

	return nil
}

// ListImages - display local Docker images
func ListImages() (listImagesErr error) {
	dockerCtx := context.Background()

	dockerClient, initClientErr := config.DockerClientInit()
	if initClientErr != nil {
		return initClientErr
	}

	images, listImageErr := dockerClient.ImageList(dockerCtx, types.ImageListOptions{})
	if listImageErr != nil {
		return fmt.Errorf("cannot list images: %v", listImageErr)
	}

	for _, image := range images {
		log.Debugf("Image id: %v", image.ID)
		log.Infof("Image tag: %v", image.RepoTags)
	}

	return nil
}
