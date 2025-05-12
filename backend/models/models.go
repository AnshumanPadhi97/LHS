package models

import "gopkg.in/yaml.v3"

type StackTemplate struct {
	Name     string          `yaml:"name"`
	Services []ServiceConfig `yaml:"services"`
}

type ServiceConfig struct {
	Name            string            `yaml:"name"`
	Image           string            `yaml:"image,omitempty"`
	BuildPath       string            `yaml:"path"`
	BuildDockerfile string            `yaml:"dockerfile"`
	Ports           []string          `yaml:"ports,omitempty"` // "host:container"
	Env             map[string]string `yaml:"env,omitempty"`
	Volumes         []string          `yaml:"volumes,omitempty"` // hostPath:containerPath
}

func ParseStackYAML(data []byte) (*StackTemplate, error) {
	var tmpl StackTemplate
	err := yaml.Unmarshal(data, &tmpl)
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}
