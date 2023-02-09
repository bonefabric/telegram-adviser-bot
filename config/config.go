package config

import (
	"flag"
	"log"

	"bonefabric/adviser/store"
)

type Config struct {
	tgToken     string
	storeDriver store.Driver
}

func Load() Config {
	if flag.Parsed() {
		log.Fatal("flag is already parsed")
	}

	var c Config

	var stdrvr string

	flag.StringVar(&c.tgToken, "tt", "", "telegram bot token")
	flag.StringVar(&stdrvr, "sd", "sqlite3", "store driver (sqlite3)")
	flag.Parse()

	if c.tgToken == "" {
		log.Fatal("tt flag required")
	}

	switch stdrvr {
	case "sqlite3":
		c.storeDriver = store.DriverSqlite3
	default:
		log.Fatal("invalid store driver")
	}

	return c
}

func (c *Config) TgToken() string {
	return c.tgToken
}

func (c *Config) StoreDriver() store.Driver {
	return c.storeDriver
}
