package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Store struct {
	Driver   string `mapstructure:"driver" default:"sqlite"`
	Host     string `mapstructure:"host"`
	Port     int64  `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db" default:"kronos"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type Email struct {
	Server   string
	Address  string
	Password string
}

type Alert struct {
	Email Email `mapstructure:"email"`
}

type Config struct {
	Port    int64 `mapstructure:"port"`
	Alert   Alert `mapstructure:"alert"`
	Logging Log   `mapstructure:"logging"`
	Store   Store `mapstructure:"store"`
}

func Parse(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
