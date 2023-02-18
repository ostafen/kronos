package config

import (
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
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
	Server   string `mapstructure:"server"`
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
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

func viperDefaults() {
	viper.SetDefault("port", 9175)
	viper.SetDefault("store.driver", "sqlite3")
}

func envKeys(m map[string]any) []string {
	keys := make([]string, 0)
	for k, v := range m {
		prefix := k

		if vm, isMap := v.(map[string]any); isMap {
			subkeys := envKeys(vm)
			for _, sk := range subkeys {
				keys = append(keys, prefix+"."+sk)
			}
		} else {
			keys = append(keys, prefix)
		}
	}
	return keys
}

func bindEnv(v any) error {
	envKeysMap := map[string]any{}
	if err := mapstructure.Decode(v, &envKeysMap); err != nil {
		return err
	}

	keys := envKeys(envKeysMap)
	for _, key := range keys {
		if err := viper.BindEnv(key); err != nil {
			return err
		}
	}
	return nil
}

func Read() (*Config, error) {
	viperDefaults()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if len(os.Args) > 1 {
		viper.SetConfigFile(os.Args[1])

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var c Config
	if err := bindEnv(c); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
