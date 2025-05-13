package models

import "gopkg.in/yaml.v3"

type StackTemplate struct {
	Name     string          `yaml:"name"`
	Services []ServiceConfig `yaml:"services"`
}

type ServiceConfig struct {
	Name            string            `yaml:"name"`
	Image           string            `yaml:"image"`
	BuildPath       string            `yaml:"path"`
	BuildDockerfile string            `yaml:"dockerfile"`
	Ports           []string          `yaml:"ports"`
	Env             map[string]string `yaml:"env"`
	Volumes         []string          `yaml:"volumes"`
}

func ParseStackYAML(data []byte) (*StackTemplate, error) {
	var tmpl StackTemplate
	err := yaml.Unmarshal(data, &tmpl)
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}
