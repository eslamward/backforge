package parser

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Schema struct {
	Configuration Configuration `yaml:"configuration"`
	Models        []Model       `yaml:"models"`
}

type Configuration struct {
	ServerConfig   ServerConfig   `yaml:"server_config"`
	DatabaseConfig DatabaseConfig `yaml:"database_config"`
}
type ServerConfig struct {
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type Model struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
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
type ForeignKey struct {
	Model    string `yaml:"model"`
	Field    string `yaml:"field"`
	OnDelete string `yaml:"on_delete"`
	OnUpdate string `yaml:"on_update"`
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
