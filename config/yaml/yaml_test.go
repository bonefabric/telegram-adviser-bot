package yaml

import (
	"os"
	"testing"
)

func TestLoadSuccess(t *testing.T) {
	validYaml := `
app:
  telegram-token: "6084774046"
  store-driver: "sqlite3"
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
	if conf.TelegramToken() != "6084774046" {
		t.Errorf("incorrect tgToken: %s", conf.TelegramToken())
	}
	if conf.StoreDriver() != "sqlite3" {
		t.Errorf("incorrect storeDriver: %s", conf.StoreDriver())
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
