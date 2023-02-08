package config

import (
	"flag"
	"log"
)

type Config struct {
	tgToken string
}

func Load() Config {
	if flag.Parsed() {
		log.Fatal("flag is already parsed")
	}

	var c Config

	flag.StringVar(&c.tgToken, "tt", "", "telegram bot token")
	flag.Parse()

	if c.tgToken == "" {
		log.Fatal("tt flag required")
	}

	return c
}

func (c *Config) TgToken() string {
	return c.tgToken
}
