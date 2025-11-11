package keyvalue

import "context"

type Store interface {
	Read(ctx context.Context, key string) (string, error)
	Write(ctx context.Context, key string, value string) error
	Exists(ctx context.Context, key string) (bool, error)
}
