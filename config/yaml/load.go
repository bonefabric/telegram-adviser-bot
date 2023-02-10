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
	TelegramToken string `yaml:"telegram-token"`
	StoreDriver   string `yaml:"store-driver"`
}

type profiling struct {
	Enabled bool   `yaml:"enabled"`
	File    string `yaml:"file"`
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
		tgToken:     c.App.TelegramToken,
		storeDriver: c.App.StoreDriver,
		profiler:    c.Profiling.Enabled,
		profileFile: c.Profiling.File,
	}, nil
}
