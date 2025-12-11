package hsm

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type decisionHandler struct {
	decision Decision
}

func (h decisionHandler) Handle(context.Context, telegram.Update) (Decision, error) {
	return h.decision, nil
}

func Just(decision Decision) Handler {
	return decisionHandler{decision: decision}
}
