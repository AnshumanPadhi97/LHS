package backend

import (
	docker "LHS/backend/Processors/Docker"
	"LHS/backend/models"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

//STACK CONTROLS

func BuildStack(tmpl *models.StackTemplate) error {
	for _, svc := range tmpl.Services {
		imageName := svc.Image

		//Build or pull image
		if svc.Build != nil {
			//TODO
			//make - image build options settings available to user
			buildOptions := types.ImageBuildOptions{
				Dockerfile: svc.Build.Dockerfile,
				Tags:       []string{svc.Code},
				Remove:     true,
			}
			err := docker.BuildImage(svc, buildOptions)
			if err != nil {
				return fmt.Errorf("build failed for %s: %w", svc.Name, err)
			}
			imageName = svc.Code
		}

		// Env
		var envVars []string
		for k, v := range svc.Env {
			envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
		}

		//Ports
		portSet := nat.PortSet{}
		portMap := nat.PortMap{}
		for _, port := range svc.Ports {
			parts := strings.Split(port, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid port format in %s", port)
			}
			hostPort := parts[0]
			containerPort := parts[1]
			portKey := nat.Port(containerPort + "/tcp")
			portSet[portKey] = struct{}{}
			portMap[portKey] = []nat.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: hostPort,
			}}
		}

		//Volumes
		var mountsList []mount.Mount
		for _, vol := range svc.Volumes {
			parts := strings.Split(vol, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid volume format in %s", vol)
			}
			mountsList = append(mountsList, mount.Mount{
				Type:   mount.TypeBind,
				Source: parts[0],
				Target: parts[1],
			})
		}

		// Create Container
		resp, err := docker.CreateContainer(
			&container.Config{
				Image:        imageName,
				Env:          envVars,
				ExposedPorts: portSet,
			},
			&container.HostConfig{
				PortBindings: portMap,
				Mounts:       mountsList,
			},
			&network.NetworkingConfig{}, nil, svc.Code)

		if err != nil {
			return fmt.Errorf("error creating container %s: %w", svc.Name, err)
		}

		//save created container id in db
		//return resp.ID
	}

	return nil
}
