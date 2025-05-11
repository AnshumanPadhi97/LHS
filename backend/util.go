package backend

import (
	model "LHS/backend/models"
	"os"

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
