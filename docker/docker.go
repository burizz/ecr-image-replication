package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/burizz/ecr-image-replication/config"
	log "github.com/sirupsen/logrus"
)

//TODO: implement interface and struct to make these methods

type DockerClient interface {
	//PullImage(image string, tag string) error
	PullImage() error
	ListImages() error
	TagImage(sourceImageTag string, targetImageTag string) (bool, error)
	//TODO: figure out variadic stuff
	PushImage(image string, tag string, options ...string) error
}

type Image struct {
	Image     string
	Tag       string
	AuthToken string
}

// PullImage - download remote Docker image
//func (d Docker) PullImage(image string) (imagePullErr error) {
func (i Image) PullImage() (imagePullErr error) {
	dockerCtx := context.Background()

	dockerClient, initClientErr := config.DockerClientInit()
	if initClientErr != nil {
		return initClientErr
	}

	image := i.Image + ":" + i.Tag

	// ImagePull makes an API request to Docker daemon to pull the image, we don't actually pull it and store it ourselves here
	reader, dockerPullErr := dockerClient.ImagePull(dockerCtx, image, types.ImagePullOptions{})
	// TODO: figure out a proper way to handle this closer
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
func (i Image) ListImages() (listImagesErr error) {
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
//TODO: make authToken optional with variadic function
//func PushImage(imageTag string, authToken string) (imagePushErr error) {
func PushImage(imageTag string, authToken string) (imagePushErr error) {
	dockerCtx := context.Background()

	// Init Docker client
	dockerClient, initClientErr := config.DockerClientInit()
	if initClientErr != nil {
		return initClientErr
	}

	// Decode authToken from base64 to string
	decodedToken, base64DecodeErr := base64.StdEncoding.DecodeString(authToken)
	if base64DecodeErr != nil {
		return fmt.Errorf("cannot decode token value: %v", base64DecodeErr)
	}

	// Separate user and token values
	token := strings.Split(string(decodedToken), ":")

	// Build auth config with user and token values
	authConfig := types.AuthConfig{
		Username: token[0],
		Password: token[1],
	}

	// Convert authConfig to JSON
	encodedJSON, jsonMarshallErr := json.Marshal(authConfig)
	if jsonMarshallErr != nil {
		return fmt.Errorf("cannot marshal AuthConfig json: %v", jsonMarshallErr)
	}

	// Encode JSON back to base64 - docker client expects it in base64 encoded json format
	authString := base64.URLEncoding.EncodeToString(encodedJSON)

	// Push image with authentication
	reader, imagePushErr := dockerClient.ImagePush(dockerCtx, imageTag, types.ImagePushOptions{RegistryAuth: authString})
	if imagePushErr != nil {
		return fmt.Errorf("cannot push image: %v, %v", imageTag, imagePushErr)
	}

	// TODO: change this to logrus as well
	// Sends ImagePush output via reader - to show download progress
	io.Copy(os.Stdout, reader)

	return nil
}
