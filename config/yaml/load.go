package yaml

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	App       appConfig `yaml:"app"`
	Profiling profiling `yaml:"profiling"`
}

type appConfig struct {
	StoreDriver  string       `yaml:"store-driver"`
	StoreOptions storeOptions `yaml:"store-options"`
	Units        units        `yaml:"units"`
}

type storeOptions struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type profiling struct {
	Enabled bool   `yaml:"enabled"`
	File    string `yaml:"file"`
}

type units struct {
	Telegram struct {
		Enabled bool   `yaml:"enabled"`
		Token   string `yaml:"token"`
	} `yaml:"telegram"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	var c config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}

	return &Config{
		storeDriver:   c.App.StoreDriver,
		storeHost:     c.App.StoreOptions.Host,
		storePort:     c.App.StoreOptions.Port,
		storeUser:     c.App.StoreOptions.User,
		storePassword: c.App.StoreOptions.Password,
		storeName:     c.App.StoreOptions.Name,
		profiler:      c.Profiling.Enabled,
		profileFile:   c.Profiling.File,
		tgToken:       c.App.Units.Telegram.Token,
		tgEnabled:     c.App.Units.Telegram.Enabled,
	}, nil
}
