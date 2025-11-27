package sessions

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
)

type Sessions interface {
	Session(id int64) Session
}

type Session interface {
	StateID(ctx context.Context) (core.ID, error)
	ChangeStateID(ctx context.Context, id core.ID) error
	Exists(ctx context.Context) (bool, error)
}
