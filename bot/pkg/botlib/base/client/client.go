package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client interface {
	Execute(ctx context.Context, method string, payload []byte) ([]byte, error)
}

type httpClient struct {
	baseURL   string
	transport *http.Client
}

func (h httpClient) Execute(ctx context.Context, method string, payload []byte) ([]byte, error) {
	request, err := h.request(ctx, fmt.Sprintf("%s/%s", h.baseURL, method), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	resp, err := h.transport.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("telegram API error %d: %s", resp.StatusCode, body)
	}
	return body, nil
}

func (h httpClient) request(ctx context.Context, url string, reader io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func HTTPClient(token string) Client {
	return httpClient{
		baseURL:   fmt.Sprintf("https://api.telegram.org/bot%s", token),
		transport: &http.Client{Timeout: 30 * time.Second},
	}
}
