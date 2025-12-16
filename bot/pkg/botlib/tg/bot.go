package tg

import (
	"context"
	"encoding/json"

	"github.com/andreychh/coopera-bot/pkg/botlib/tg/transport"
)

type bot struct {
	dataSource transport.Client
}

func (b bot) Chat(id int64) Chat {
	return NewChat(id, b.dataSource)
}

type req struct {
	CallbackQueryId string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
}

func (b bot) AnswerCallbackQuery(ctx context.Context, id string) error {
	req := req{
		CallbackQueryId: id,
		Text:            "",
		ShowAlert:       false,
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = b.dataSource.Execute(ctx, "answerCallbackQuery", payload)
	if err != nil {
		return err
	}
	return nil
}

func NewBot(dataSource transport.Client) Bot {
	return bot{
		dataSource: dataSource,
	}
}
