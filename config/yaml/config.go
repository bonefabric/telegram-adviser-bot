package yaml

type Config struct {
	storeDriver   string
	storeHost     string
	storePort     int
	storeUser     string
	storePassword string
	storeName     string
	profiler      bool
	profileFile   string
	tgToken       string
	tgEnabled     bool
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

func (c *Config) StoreHost() string {
	return c.storeHost
}

func (c *Config) StorePort() int {
	return c.storePort
}

func (c *Config) StoreUser() string {
	return c.storeUser
}

func (c *Config) StorePassword() string {
	return c.storePassword
}

func (c *Config) StoreName() string {
	return c.storeName
}

func (c *Config) TelegramToken() string {
	return c.tgToken
}

func (c *Config) TelegramEnabled() bool {
	return c.tgEnabled
}
