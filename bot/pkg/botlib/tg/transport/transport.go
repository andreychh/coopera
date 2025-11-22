package transport

import "context"

type Client interface {
	Execute(ctx context.Context, method string, payload []byte) ([]byte, error)
}
