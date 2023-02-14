package yaml

import (
	"os"
	"testing"
)

func TestLoadSuccess(t *testing.T) {
	validYaml := `
app:
  store-driver: "sqlite3"
  store-options:
    host: "store-host"
    port: 3306
    user: "store-user"
    password: "store-pass"
    name: "store-name"
  units:
    telegram:
      enabled: true
      token: "some-token"
profiling:
  enabled: false
  file: "cpu.prof"
`
	tmpFile, err := os.CreateTemp(".", "config.*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile.Name())

	_, err = tmpFile.Write([]byte(validYaml))
	if err != nil {
		t.Fatalf("failed to write data to temp file: %s", err)
	}

	conf, err := Load(tmpFile.Name())
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if conf == nil {
		t.Error("nil config")
	}
	if conf.StoreDriver() != "sqlite3" {
		t.Errorf("incorrect storeDriver: %s", conf.StoreDriver())
	}
	if conf.StoreHost() != "store-host" {
		t.Errorf("incorrect storeHost: %s", conf.StoreHost())
	}
	if conf.StorePort() != 3306 {
		t.Errorf("incorrect storePort: %d", conf.StorePort())
	}
	if conf.StoreUser() != "store-user" {
		t.Errorf("incorrect storeUser: %s", conf.StoreUser())
	}
	if conf.StorePassword() != "store-pass" {
		t.Errorf("incorrect storePassword: %s", conf.StorePassword())
	}
	if conf.StoreName() != "store-name" {
		t.Errorf("incorrect storeName: %s", conf.StoreName())
	}
	if !conf.TelegramEnabled() {
		t.Errorf("incorrect telegramEnabled: %t", conf.TelegramEnabled())
	}
	if conf.TelegramToken() != "some-token" {
		t.Errorf("incorrect telegramToken: %s", conf.TelegramToken())
	}
	if conf.Profiling() {
		t.Errorf("incorrect profiling: %t", conf.Profiling())
	}
	if conf.ProfileFile() != "cpu.prof" {
		t.Errorf("incorrect profiler file: %s", conf.ProfileFile())
	}
}

func TestLoadFail(t *testing.T) {
	invalidYaml := `
telegram-token = 6084774046
`
	tmpFile, err := os.CreateTemp(".", "config.*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile.Name())

	_, err = tmpFile.Write([]byte(invalidYaml))
	if err != nil {
		t.Fatalf("failed to write data to temp file: %s", err)
	}

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Error("error nil on invalid file")
	}
}
