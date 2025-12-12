package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

var ErrRecordAlreadyExists = errors.New("record already exists")

type Client interface {
	Get(ctx context.Context, path string) ([]byte, error)
	Post(ctx context.Context, path string, payload []byte) ([]byte, error)
}

type httpClient struct {
	baseURL string
	client  *http.Client
}

func (c httpClient) Get(ctx context.Context, path string) ([]byte, error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing base url: %w", err)
	}
	relativeURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parsing relative path: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, baseURL.ResolveReference(relativeURL).String(), nil,
	)
	if err != nil {
		return nil, fmt.Errorf("creating GET request for %s: %w", path, err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing GET request for %s: %w", path, err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, APIError{
			StatusCode: resp.StatusCode,
			URL:        path,
			Body:       body,
		}
	}
	return io.ReadAll(resp.Body)
}

func (c httpClient) Post(ctx context.Context, path string, payload []byte) ([]byte, error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing base url: %w", err)
	}
	relativeURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parsing relative path: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, baseURL.ResolveReference(relativeURL).String(), bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("creating POST request for %s: %w", path, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing POST request for %s: %w", path, err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()
	if resp.StatusCode == http.StatusConflict {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading body: %w", err)
		}
		errorResult := gjson.GetBytes(body, "error")
		if !errorResult.Exists() {
			return nil, fmt.Errorf("error message does not exists in body")
		}
		msg := errorResult.String()
		if strings.Contains(msg, "record already exists") {
			return nil, ErrRecordAlreadyExists
		}
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d for POST %s. Body: %s",
			resp.StatusCode,
			path,
			errorBody,
		)
	}
	return io.ReadAll(resp.Body)
}

func HTTPClient(baseURL string, client *http.Client) Client {
	return httpClient{
		baseURL: baseURL,
		client:  client,
	}
}
