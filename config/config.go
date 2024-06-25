package config

import (
	"datapointbackend/pkg/database"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Http     Http            `yaml:"http"`
	Database database.Config `yaml:"database"`
}

type Http struct {
	Addr string `yaml:"addr"`
}

func New() (*Config, error) {
	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, err
	}

	cfg := new(Config)

	return cfg, yaml.Unmarshal(data, cfg)
}
