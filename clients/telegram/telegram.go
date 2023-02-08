package telegram

import "net/http"

type Telegram struct {
	token      string
	httpClient http.Client
}

func New(token string) Telegram {
	return Telegram{
		token:      token,
		httpClient: http.Client{},
	}
}
