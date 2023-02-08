package telegram

import "encoding/json"

type Response struct {
	Ok          bool             `json:"ok"`
	Description *string          `json:"description"`
	Result      *json.RawMessage `json:"result"`
}

type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"message"`
}

type Message struct {
	ID   int     `json:"message_id"`
	From *User   `json:"from"`
	Text *string `json:"text"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
}
