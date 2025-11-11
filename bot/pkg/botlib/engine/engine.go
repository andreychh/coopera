package engine

import (
	"context"
)

type Engine interface {
	Start(ctx context.Context)
}
