package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Store struct {
	Driver   string `yaml:"driver" default:"sqlite"`
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db" default:"kronos"`
}

type Log struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Email struct {
	Server   string
	Address  string
	Password string
}

type Alert struct {
	Email Email `yaml:"email"`
}

type Service struct {
	Port    int64 `yaml:"port"`
	Alert   Alert `yaml:"alert"`
	Logging Log   `yaml:"logging"`
}

type Config struct {
	Service Service `yaml:"service"`
	Store   Store   `yaml:"store"`
}

func Parse(filename string) (*Config, error) {
	var c Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
