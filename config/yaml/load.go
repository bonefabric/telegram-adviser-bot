package yaml

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type config struct {
	App appConfig `yaml:"app"`
}

type appConfig struct {
	TelegramToken string `yaml:"telegram-token"`
	StoreDriver   string `yaml:"store-driver"`
}

func Load(path string) (conf *Config, err error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %s", err)
	}
	defer func(file *os.File) {
		if e := file.Close(); e != nil && err == nil {
			err = e
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	var c config
	if err = yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}
	return &Config{
		tgToken:     c.App.TelegramToken,
		storeDriver: c.App.StoreDriver,
	}, nil
}
