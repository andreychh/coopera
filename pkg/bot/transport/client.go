// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/andreychh/coopera/pkg/bot/api"
)

type telegramClient struct {
	token  string
	client HTTPClient
}

func (c telegramClient) SendRequest(ctx context.Context, method api.Method, requestBody, responseBody any) (err error) {
	marshaled, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshaling request body: %w", err)
	}
	request, err := c.createRequest(ctx, c.url(method), bytes.NewReader(marshaled))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	response, err := c.client.Do(request)
	if err != nil {
		return MappedError(fmt.Errorf("sending request: %w", err), c.token, "REDACTED_TOKEN")
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(response.Body)
	var envelope api.Envelope
	err = json.NewDecoder(response.Body).Decode(&envelope)
	if err != nil {
		return fmt.Errorf("decoding envelope: %w", err)
	}
	if !envelope.Ok {
		return api.NewError(envelope)
	}
	err = json.Unmarshal(envelope.Result, responseBody)
	if err != nil {
		return fmt.Errorf("unmarshaling result: %w", err)
	}
	return nil
}

func (c telegramClient) createRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func (c telegramClient) url(method api.Method) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", c.token, method)
}

func NewClient(token string, client HTTPClient) TelegramClient {
	return telegramClient{
		token:  token,
		client: client,
	}
}
