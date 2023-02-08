package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const ApiHost = "api.telegram.org"

type Telegram struct {
	token      string
	httpClient http.Client
	basePath   string
	offset     int
}

func New(token string) Telegram {
	return Telegram{
		token:      token,
		httpClient: http.Client{},
		basePath:   basePath(token),
	}
}

func (t *Telegram) Updates(ctx context.Context) ([]Update, error) {
	var v url.Values
	v.Add("offset", strconv.Itoa(t.offset))

	data, err := t.doRequest(ctx, t.url("getUpdates", v), nil)
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

func (t *Telegram) parseResponse(data []byte) ([]byte, error) {
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	if !r.Ok {
		msg := "response ok: false"
		if r.Description != nil {
			msg += "; description: " + *r.Description
		}
		return nil, errors.New(msg)
	}

	if r.Result == nil {
		return nil, nil
	}
	return *r.Result, nil
}

func (t *Telegram) url(method string, v url.Values) string {
	u := url.URL{
		Scheme:   "https",
		Host:     ApiHost,
		Path:     path.Join(t.basePath, method),
		RawQuery: v.Encode(),
	}
	return u.String()
}

func (t *Telegram) doRequest(ctx context.Context, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %s", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return io.ReadAll(resp.Body)
}

func basePath(token string) string {
	return "bot" + token
}
