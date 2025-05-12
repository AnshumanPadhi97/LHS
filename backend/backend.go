package backend

import (
	docker "LHS/backend/Processors/Docker"
	"LHS/backend/Processors/db"
	"LHS/backend/models"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

//STACK CONTROLS

func BuildStack(tmpl *models.StackTemplate) error {
	//save stack data
	res, err := db.DB.Exec("INSERT INTO stacks (name) VALUES (?)", tmpl.Name)
	if err != nil {
		return fmt.Errorf("failed to insert stack: %w", err)
	}

	stackID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get stack ID: %w", err)
	}

	for _, svc := range tmpl.Services {
		imageName := svc.Image

		//Build image
		if svc.BuildPath != "" && svc.BuildDockerfile != "" {
			imgName := svc.Name + "-" + uuid.New().String()
			buildOptions := types.ImageBuildOptions{
				Dockerfile: svc.BuildDockerfile,
				Tags:       []string{imgName},
				Remove:     true,
			}
			err := docker.BuildImage(svc, buildOptions)
			if err != nil {
				return fmt.Errorf("build failed for %s: %w", svc.Name, err)
			}
			imageName = imgName
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
			&network.NetworkingConfig{}, nil, "Container-"+svc.Name+"-"+uuid.NewString())

		if err != nil {
			return fmt.Errorf("error creating container %s: %w", svc.Name, err)
		}

		// Save service to DB
		_, err = db.DB.Exec(`
			INSERT INTO stack_services 
			(stack_id, container_id, name, image, build_path, build_dockerfile, ports, env, volumes)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			stackID, resp.ID, svc.Name, imageName, svc.BuildPath, svc.BuildDockerfile,
			strings.Join(svc.Ports, ","),
			encodeEnv(svc.Env),
			strings.Join(svc.Volumes, ","),
		)
		if err != nil {
			return fmt.Errorf("failed to insert service %s: %w", svc.Name, err)
		}

		fmt.Printf("Service %s started as container %s\n", svc.Name, resp.ID)
	}

	return nil
}
