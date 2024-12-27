package config

import (
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Store struct {
	Path string `mapstructure:"path"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type Config struct {
	Port    int64 `mapstructure:"port"`
	Logging Log   `mapstructure:"logging"`
	Store   Store `mapstructure:"store"`
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

func viperDefaults() {
	viper.SetDefault("store.path", "kronos.db")
	viper.SetDefault("port", 9175)
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
