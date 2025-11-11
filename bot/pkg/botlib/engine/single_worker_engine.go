package engine

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type singleWorkerEngine struct {
	clause  core.Clause
	updates telegram.UpdatesChannel
}

func (e singleWorkerEngine) Start(ctx context.Context) {
	for {
		select {
		case update, ok := <-e.updates:
			if !ok {
				return
			}
			_, _ = e.clause.TryExecute(ctx, update)
		case <-ctx.Done():
			return
		}
	}
}

func SingleWorkerEngine(clause core.Clause, updates telegram.UpdatesChannel) Engine {
	return singleWorkerEngine{clause: clause, updates: updates}
}
