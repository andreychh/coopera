package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var ErrRecordAlreadyExists = errors.New("record already exists")

type Client interface {
	Get(ctx context.Context, path string, resp any) error
	Post(ctx context.Context, path string, req any, resp any) error
	Patch(ctx context.Context, path string, req any, resp any) error
}
type httpClient struct {
	baseURL string
	client  *http.Client
}

func (c *httpClient) Get(ctx context.Context, path string, resp any) error {
	return c.do(ctx, http.MethodGet, path, nil, resp)
}

func (c *httpClient) Post(ctx context.Context, path string, req any, resp any) error {
	return c.do(ctx, http.MethodPost, path, req, resp)
}

func (c *httpClient) Patch(ctx context.Context, path string, req any, resp any) error {
	return c.do(ctx, http.MethodPatch, path, req, resp)
}

func (c *httpClient) do(ctx context.Context, method, path string, reqBody any, respTarget any) (err error) {
	fullURL, err := c.buildURL(path)
	if err != nil {
		return fmt.Errorf("building url: %w", err)
	}
	var bodyReader io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("creating %s request for %s: %w", method, path, err)
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing %s request for %s: %w", method, path, err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return APIError{
			StatusCode: resp.StatusCode,
			URL:        path,
			Body:       body,
		}
	}
	if respTarget != nil {
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if err := json.NewDecoder(resp.Body).Decode(respTarget); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("decoding response from %s: %w", path, err)
		}
	}
	return nil
}

func (c *httpClient) buildURL(path string) (string, error) {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}
	rel, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(rel).String(), nil
}

func HTTPClient(baseURL string, client *http.Client) Client {
	return &httpClient{
		baseURL: baseURL,
		client:  client,
	}
}
