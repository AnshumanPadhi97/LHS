package backend

import (
	docker "LHS/backend/Processors/Docker"
	"LHS/backend/Processors/db"
	"LHS/backend/models"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

//STACK CONTROLS

func CleanupAll() error {
	// Step 1: Stop and remove all stack containers
	services, err := db.GetAllServices()
	if err != nil {
		return fmt.Errorf("failed to retrieve services: %w", err)
	}

	for _, svc := range services {
		fmt.Printf("🧹 Removing container: %s\n", svc.ContainerID)

		// Stop
		_ = docker.StopContainer(svc.ContainerID, container.StopOptions{})

		// Remove
		err := docker.DeleteContainer(svc.ContainerID, container.RemoveOptions{Force: true, RemoveVolumes: true, RemoveLinks: true})
		if err != nil {
			fmt.Printf("⚠️ Failed to remove container %s: %v\n", svc.ContainerID, err)
		}
	}

	// Step 2: Remove images
	imageIDs := map[string]bool{}
	for _, svc := range services {
		if svc.Image != "" && !imageIDs[svc.Image] {
			fmt.Printf("🧹 Removing image: %s\n", svc.Image)
			_, err := docker.DeleteImage(svc.Image, image.RemoveOptions{Force: true, PruneChildren: true})
			if err != nil {
				fmt.Printf("⚠️ Failed to remove image %s: %v\n", svc.Image, err)
			}
			imageIDs[svc.Image] = true
		}
	}

	// Step 3: Clear DB
	if err := db.DeleteAllStacks(); err != nil {
		return fmt.Errorf("failed to delete stacks: %w", err)
	}

	fmt.Println("🧼 Cleanup complete: All containers, services, stacks removed.")
	return nil
}

func BuildStack(tmpl *models.StackTemplate) error {
	//save stack data
	stackID, err := db.CreateStack(tmpl.Name)
	if err != nil {
		return fmt.Errorf("failed to create stack: %w", err)
	}

	for _, svc := range tmpl.Services {
		containerName := "Container-" + svc.Name + "-" + uuid.NewString()

		imageName, err := resolveImage(svc)
		if err != nil {
			return fmt.Errorf("failed to resolve image for %s: %w", svc.Name, err)
		}

		envVars := formatEnvVars(svc.Env)

		portSet, portMap, err := mapPorts(svc.Ports)
		if err != nil {
			return fmt.Errorf("port mapping error in service %s: %w", svc.Name, err)
		}

		mounts, err := mapVolumes(svc.Volumes)
		if err != nil {
			return fmt.Errorf("volume mapping error in service %s: %w", svc.Name, err)
		}

		resp, err := docker.CreateContainer(
			&container.Config{
				Image:        imageName,
				Env:          envVars,
				ExposedPorts: portSet,
			},
			&container.HostConfig{
				PortBindings: portMap,
				Mounts:       mounts,
			},
			&network.NetworkingConfig{}, nil, containerName)

		if err != nil {
			return fmt.Errorf("error creating container %s: %w", svc.Name, err)
		}

		_, err = db.DB.Exec(`
			INSERT INTO stack_services 
			(stack_id, container_id, name, image, build_path, build_dockerfile, ports, env, volumes)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			stackID, containerName, svc.Name, imageName, svc.BuildPath, svc.BuildDockerfile,
			strings.Join(svc.Ports, ","), encodeEnv(svc.Env), strings.Join(svc.Volumes, ","))

		if err != nil {
			return fmt.Errorf("failed to insert service %s: %w", svc.Name, err)
		}

		fmt.Printf("✅ Stack '%s' | Service '%s' created as container '%s' (ID: '%s')\n", tmpl.Name, svc.Name, containerName, resp.ID)
	}

	fmt.Printf("✅ Stack '%s', successfully built.", tmpl.Name)

	return nil
}

func BuildStackFromDB(stackID int64) error {
	stack, err := db.GetStackById(stackID)
	if err != nil {
		return fmt.Errorf("failed to retrieve stack: %w", err)
	}

	services, err := db.GetServicesByStackID(stackID)
	if err != nil {
		return fmt.Errorf("failed to retrieve services for stack: %w", err)
	}

	for _, svc := range services {
		containerName := "Container-" + svc.Name + "-" + uuid.NewString()

		imageName, err := resolveImageFromDB(svc)
		if err != nil {
			return fmt.Errorf("failed to resolve image for %s: %w", svc.Name, err)
		}

		envVars := strings.Split(svc.Env, ",")

		portList := strings.Split(svc.Ports, ",")
		portSet, portMap, err := mapPorts(portList)
		if err != nil {
			return fmt.Errorf("port mapping error in service %s: %w", svc.Name, err)
		}

		volList := strings.Split(svc.Volumes, ",")
		mounts, err := mapVolumes(volList)
		if err != nil {
			return fmt.Errorf("volume mapping error in service %s: %w", svc.Name, err)
		}

		resp, err := docker.CreateContainer(
			&container.Config{
				Image:        imageName,
				Env:          envVars,
				ExposedPorts: portSet,
			},
			&container.HostConfig{
				PortBindings: portMap,
				Mounts:       mounts,
			},
			&network.NetworkingConfig{}, nil, containerName)

		if err != nil {
			return fmt.Errorf("error creating container %s: %w", svc.Name, err)
		}

		// Update container ID and container name in DB
		svc.ContainerID = containerName
		svc.Image = imageName
		if err := db.UpdateService(svc); err != nil {
			return fmt.Errorf("failed to update service record: %w", err)
		}

		fmt.Printf("✅ Stack '%s' | Service '%s' created as container '%s' (ID: '%s')\n", stack.Name, svc.Name, containerName, resp.ID)
	}

	fmt.Printf("✅ Stack '%s', successfully built.\n", stack.Name)
	return nil
}

func RunStack(stackID int64) error {
	query := `
		SELECT s.name AS stack_name, ss.name AS service_name, ss.container_id
		FROM stack_services ss
		JOIN stacks s ON ss.stack_id = s.id
		WHERE ss.stack_id = ?
	`
	rows, err := db.DB.Query(query, stackID)
	if err != nil {
		return fmt.Errorf("failed to query services for stack %d: %w", stackID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var stackName, serviceName, containerName string
		if err := rows.Scan(&stackName, &serviceName, &containerName); err != nil {
			return fmt.Errorf("failed to scan service data: %w", err)
		}

		err := docker.RunContainer(containerName, container.StartOptions{})
		if err != nil {
			return fmt.Errorf("failed to start container %s: %w", containerName, err)
		}

		fmt.Printf("🟢 Stack: %s | Service: %s | Container: %s started successfully\n", stackName, serviceName, containerName)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error reading rows: %w", err)
	}

	return nil
}

func resolveImage(svc models.ServiceConfig) (string, error) {
	if svc.Image != "" {
		if _, err := docker.PullImage(svc.Image, image.PullOptions{}); err != nil {
			return "", fmt.Errorf("failed to pull image %s: %w", svc.Image, err)
		}
		return svc.Image, nil
	}

	if svc.BuildPath != "" && svc.BuildDockerfile != "" {
		imageTag := svc.Name + "-" + uuid.NewString()
		opts := types.ImageBuildOptions{
			Dockerfile: svc.BuildDockerfile,
			Tags:       []string{imageTag},
			Remove:     true,
		}
		if err := docker.BuildImage(svc, opts); err != nil {
			return "", fmt.Errorf("build failed: %w", err)
		}
		return imageTag, nil
	}

	return "", fmt.Errorf("no image or build configuration provided")
}

func resolveImageFromDB(svc db.StackService) (string, error) {
	if svc.Image != "" {
		if _, err := docker.PullImage(svc.Image, image.PullOptions{}); err != nil {
			return "", fmt.Errorf("failed to pull image %s: %w", svc.Image, err)
		}
		return svc.Image, nil
	}

	if svc.BuildPath != "" && svc.BuildDockerfile != "" {
		imageTag := svc.Name + "-" + uuid.NewString()
		opts := types.ImageBuildOptions{
			Dockerfile: svc.BuildDockerfile,
			Tags:       []string{imageTag},
			Remove:     true,
		}
		modelSvc := models.ServiceConfig{
			Name:            svc.Name,
			BuildPath:       svc.BuildPath,
			BuildDockerfile: svc.BuildDockerfile,
		}
		if err := docker.BuildImage(modelSvc, opts); err != nil {
			return "", fmt.Errorf("build failed: %w", err)
		}
		return imageTag, nil
	}

	return "", fmt.Errorf("no image or build configuration provided")
}

func formatEnvVars(env map[string]string) []string {
	var vars []string
	for k, v := range env {
		vars = append(vars, fmt.Sprintf("%s=%s", k, v))
	}
	return vars
}

func mapPorts(ports []string) (nat.PortSet, nat.PortMap, error) {
	portSet := nat.PortSet{}
	portMap := nat.PortMap{}
	for _, p := range ports {
		parts := strings.Split(p, ":")
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid port format: %s", p)
		}
		hostPort, containerPort := parts[0], parts[1]
		portKey := nat.Port(containerPort + "/tcp")
		portSet[portKey] = struct{}{}
		portMap[portKey] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: hostPort,
		}}
	}
	return portSet, portMap, nil
}

func mapVolumes(vols []string) ([]mount.Mount, error) {
	var mounts []mount.Mount
	for _, v := range vols {
		parts := strings.Split(v, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid volume format: %s", v)
		}
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: parts[0],
			Target: parts[1],
		})
	}
	return mounts, nil
}

func encodeEnv(env map[string]string) string {
	var parts []string
	for k, v := range env {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ",")
}
