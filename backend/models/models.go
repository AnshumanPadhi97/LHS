package models

type StackTemplate struct {
	Code        string            `yaml:"code"`
	Name        string            `yaml:"name"`
	Label       string            `yaml:"label"`
	Description string            `yaml:"description"`
	Version     int               `yaml:"version"`
	Tags        []string          `yaml:"tags"`
	Volumes     map[string]string `yaml:"volumes,omitempty"`
	Services    []ServiceConfig   `yaml:"services"`
}

type ServiceConfig struct {
	Name      string            `yaml:"name"`
	Image     string            `yaml:"image,omitempty"` // optional if build is provided
	Build     *BuildConfig      `yaml:"build,omitempty"` // optional if image is provided
	Ports     []string          `yaml:"ports,omitempty"` // "host:container"
	Env       map[string]string `yaml:"env,omitempty"`
	Volumes   []string          `yaml:"volumes,omitempty"`    // hostPath:containerPath
	DependsOn []string          `yaml:"depends_on,omitempty"` // names of services
	Tunnel    bool              `yaml:"tunnel,omitempty"`     // expose to public via tunnel
}

type BuildConfig struct {
	Context    string `yaml:"context"`
	Dockerfile string `yaml:"dockerfile"`
}
