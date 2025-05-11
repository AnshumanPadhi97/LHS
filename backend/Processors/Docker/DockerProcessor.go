package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

var cli *client.Client

//MANAGE CLIENT

func InitDockerClient() error {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return err
}

func CloseClient() error {
	return cli.Close()
}

//CONTAINERS HANDLING

func CreateContainer(options container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error) {
	return cli.ContainerCreate(context.Background(), &options, hostConfig, networkingConfig, platform, containerName)
}

func DeleteContainer(containerID string, options container.RemoveOptions) error {
	return cli.ContainerRemove(context.Background(), containerID, options)
}

func ListContainers(options container.ListOptions) ([]container.Summary, error) {
	return cli.ContainerList(context.Background(), options)
}

func RunContainer(containerId string, options container.StartOptions) error {
	return cli.ContainerStart(context.Background(), containerId, options)
}

func StopContainer(containerId string, options container.StopOptions) error {
	return cli.ContainerStop(context.Background(), containerId, options)
}

//IMAGES HANDLING

func ListImages(options image.ListOptions) ([]image.Summary, error) {
	return cli.ImageList(context.Background(), options)
}

func PullImage(refStr string, options image.PullOptions) (io.ReadCloser, error) {
	return cli.ImagePull(context.Background(), refStr, options)
}

func DeleteImage(imageId string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
	return cli.ImageRemove(context.Background(), imageId, options)
}
