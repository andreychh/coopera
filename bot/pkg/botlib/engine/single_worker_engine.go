package engine

import (
	"context"
	"log"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type singleWorkerEngine struct {
	token  string
	action core.Action
}

func (e singleWorkerEngine) Start(ctx context.Context) {
	updates, err := e.updates()
	if err != nil {
		log.Fatalf("Failed to initialize updates: %v", err)
		return
	}
	for {
		select {
		case update, ok := <-updates:
			if !ok {
				return
			}
			e.action.Perform(ctx, update)
		case <-ctx.Done():
			return
		}
	}
}

func (e singleWorkerEngine) updates() (telegram.UpdatesChannel, error) {
	bot, err := telegram.NewBotAPI(e.token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	u := telegram.NewUpdate(0)
	u.Timeout = 60
	return bot.GetUpdatesChan(u), nil
}

func SingleWorkerEngine(token string, action core.Action) Engine {
	return singleWorkerEngine{token: token, action: action}
}
