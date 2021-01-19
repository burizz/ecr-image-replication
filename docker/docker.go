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
	//defer reader.Close()
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
	//log.Infof("Test: %v", string(contents))

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

// TagImage - changes image tag, e.g. registry.example.com/myimage:latest
func TagImage(sourceImageTag string, targetImageTag string) (ok bool, imageTagErr error) {
	dockerCtx := context.Background()

	dockerClient, initClientErr := config.DockerClientInit()
	if initClientErr != nil {
		return false, initClientErr
	}

	if imageTagErr := dockerClient.ImageTag(dockerCtx, sourceImageTag, targetImageTag); imageTagErr != nil {
		return false, fmt.Errorf("cannot tag image; source tag [%v] - target tag [%v]; Error: %v", sourceImageTag, targetImageTag, imageTagErr)
	}

	log.Infof("Image tagged successfully; source tag [%v] - target tag [%v]", sourceImageTag, targetImageTag)
	return true, nil
}

// PushImage - to private registry; image must already have tag which references the registry, e.g. registry.example.com/myimage:latest
func PushImage(imageTag string) (imagePushErr error) {
	dockerCtx := context.Background()

	dockerClient, initClientErr := config.DockerClientInit()
	if initClientErr != nil {
		return initClientErr
	}

	reader, imagePushErr := dockerClient.ImagePush(dockerCtx, imageTag, types.ImagePushOptions{})
	if imagePushErr != nil {
		return fmt.Errorf("cannot push image: %v, %v", imageTag, imagePushErr)
	}

	// TODO: change this to logrus as well
	// Sends ImagePush output via reader - to show download progress
	io.Copy(os.Stdout, reader)

	return nil
}

// DockerLogin - to registry
func DockerLogin(registryName string) (dockerLoginErr error) {
	dockerCtx := context.Background()

	dockerClient, initClientErr := config.DockerClientInit()
	if initClientErr != nil {
		return initClientErr
	}

	auth, dockerLoginErr := dockerClient.RegistryLogin(dockerCtx, types.AuthConfig{})
	if dockerLoginErr != nil {
		return fmt.Errorf("cannot login to docker registry %v ; %v", registryName, dockerLoginErr)
	}

	//TODO: finish this
	fmt.Println(auth)
	return nil
}
