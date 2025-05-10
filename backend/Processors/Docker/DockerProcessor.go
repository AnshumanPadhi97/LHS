package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

var cli *client.Client

func InitDockerClient() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("InitDockerClient : " + err.Error())
	}
}

func CloseClient() {
	cli.Close()
}

func ListContainers(options container.ListOptions) []container.Summary {
	if cli == nil {
		return nil
	}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		fmt.Println("ListContainers : " + err.Error())
		return nil
	}
	return containers
}

func ListImages(options image.ListOptions) []image.Summary {
	if cli == nil {
		return nil
	}
	images, err := cli.ImageList(context.Background(), options)
	if err != nil {
		fmt.Println("ListImages : " + err.Error())
		return nil
	}
	return images
}

func CreateContainer(options container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) string {
	if cli == nil {
		return ""
	}
	res, err := cli.ContainerCreate(context.Background(), &options, hostConfig, networkingConfig, platform, containerName)
	if err != nil {
		fmt.Println("CreateContainer : " + err.Error())
		return ""
	}
	return res.ID
}

func RunContainer(containerId string, options container.StartOptions) string {
	if cli == nil {
		return ""
	}
	err := cli.ContainerStart(context.Background(), containerId, options)
	if err != nil {
		fmt.Println("RunContainer : " + err.Error())
		return ""
	}
	return "container running"
}

func StopContainer(containerId string, options container.StopOptions) string {
	if cli == nil {
		return ""
	}
	err := cli.ContainerStop(context.Background(), containerId, options)
	if err != nil {
		fmt.Println("StopContainer : " + err.Error())
		return ""
	}
	return "container stopped"
}
