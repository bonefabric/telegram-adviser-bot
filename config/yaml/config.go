package yaml

type Config struct {
	tgToken     string
	storeDriver string
}

func (c Config) TelegramToken() string {
	return c.tgToken
}

func (c Config) StoreDriver() string {
	return c.storeDriver
}
