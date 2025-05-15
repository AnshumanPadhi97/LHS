package db

type Stack struct {
	ID        int64
	Name      string
	CreatedAt string
	Status    string
	LastRunAt string
}

type StackService struct {
	ID              int64
	StackID         int64
	ContainerID     string
	Name            string
	Image           string
	BuildPath       string
	BuildDockerfile string
	Ports           string
	Env             string
	Volumes         string
	Status          string
	LastRunAt       string
}
