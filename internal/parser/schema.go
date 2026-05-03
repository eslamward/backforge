package parser

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ForeignKey struct {
	Model    string `yaml:"model"`
	Field    string `yaml:"field"`
	OnDelete string `yaml:"on_delete"`
	OnUpdate string `yaml:"on_update"`
}

type Field struct {
	Name          string      `yaml:"name"`
	Type          string      `yaml:"type"`
	MinLength     int         `yaml:"min_length"`
	MaxLength     int         `yaml:"max_length"`
	MinValue      int         `yaml:"min_value"`
	MaxValue      int         `yaml:"max_value"`
	Primary       bool        `yaml:"primary"`
	AutoIncrement bool        `yaml:"auto_increment"`
	NotNull       bool        `yaml:"not_null"`
	Unique        bool        `yaml:"unique"`
	Default       string      `yaml:"default"`
	Check         string      `yaml:"check"`
	ForeignKey    *ForeignKey `yaml:"foreign_key"`
}

type Model struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

type Schema struct {
	Models []Model `yaml:"models"`
}

func Parse(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var schema Schema
	err = yaml.Unmarshal(data, &schema)
	return &schema, err
}
