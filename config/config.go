package config

type Config interface {
	TelegramToken() string
	StoreDriver() string
}
