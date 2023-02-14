package config

type Config interface {
	StoreDriver() string
	Profiling() bool
	ProfileFile() string
	StoreHost() string
	StorePort() int
	StoreUser() string
	StorePassword() string
	StoreName() string
	TelegramEnabled() bool
	TelegramToken() string
}
