package config

type Config interface {
	TelegramToken() string
	StoreDriver() string
	Profiling() bool
	ProfileFile() string
	StoreHost() string
	StorePort() int
	StoreUser() string
	StorePassword() string
	StoreName() string
}
