package telegram

type SendMessageOptions struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}
