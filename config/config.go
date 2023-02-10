package config

type Config interface {
	TelegramToken() string
	StoreDriver() string
	Profiling() bool
	ProfileFile() string
}
