package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const (
	MethodGetUpdates  = "getUpdates"
	MethodSendMessage = "sendMessage"
)

func (t *Telegram) Updates(ctx context.Context) ([]Update, error) {
	v := url.Values{}
	v.Add("offset", strconv.Itoa(t.offset))

	data, err := t.doRequest(ctx, t.url(MethodGetUpdates, v), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %s", err)
	}

	rawUpd, err := t.parseResponse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %s", err)
	}

	u := make([]Update, 0)
	if err := json.Unmarshal(rawUpd, &u); err != nil {
		return nil, fmt.Errorf("failed to parse updates: %s", err)
	}

	if u != nil && len(u) > 0 {
		t.offset = u[len(u)-1].ID + 1
	}
	return u, nil
}

func (t *Telegram) SendMessage(ctx context.Context, opts SendMessageOptions) error {
	b, err := json.Marshal(opts)
	if err != nil {
		return fmt.Errorf("failed to marshal sendMessage options: %s", err)
	}
	resp, err := t.doRequest(ctx, t.url(MethodSendMessage, url.Values{}), bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to send message: %s", err)
	}
	_, err = t.parseResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to send message: %s", err)
	}
	return nil
}
