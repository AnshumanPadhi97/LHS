package backend

import (
	model "LHS/backend/models"
	"fmt"
	"os"
	"strings"

	"github.com/stretchr/testify/assert/yaml"
)

func LoadTemplateFromFile(path string) (*model.StackTemplate, error) {
	var tmpl model.StackTemplate
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &tmpl)
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func encodeEnv(env map[string]string) string {
	var parts []string
	for k, v := range env {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ",")
}
