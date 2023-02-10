package yaml

type Config struct {
	tgToken     string
	storeDriver string
	profiler    bool
	profileFile string
}

func (c *Config) TelegramToken() string {
	return c.tgToken
}

func (c *Config) StoreDriver() string {
	return c.storeDriver
}

func (c *Config) Profiling() bool {
	return c.profiler
}

func (c *Config) ProfileFile() string {
	return c.profileFile
}
