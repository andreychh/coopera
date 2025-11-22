package app

import (
	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/http"
	t "github.com/andreychh/coopera-bot/internal/transport"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	dialogueskeyvalue "github.com/andreychh/coopera-bot/pkg/botlib/dialogues/keyvalue"
	"github.com/andreychh/coopera-bot/pkg/botlib/engine"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	formskeyvalue "github.com/andreychh/coopera-bot/pkg/botlib/forms/keyvalue"
	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/transport"
)

func Store() keyvalue.Store {
	return keyvalue.MemoryStore()
}

func Dialogues(store keyvalue.Store) dialogues.Dialogues {
	return dialogueskeyvalue.KeyValueDialogues(store)
}

func Forms(store keyvalue.Store) forms.Forms {
	return formskeyvalue.KeyValueForms(store)
}

func Client(token string) transport.Client {
	return transport.HTTPClient(token)
}

func Bot(client transport.Client) tg.Bot {
	return tg.NewBot(client)
}

func Community(client t.Client) domain.Community {
	return http.Community(client)
}

func Engine(token string, clause core.Clause) engine.Engine {
	return engine.ShutdownEngine(
		engine.SingleWorkerEngine(
			token, clause,
		),
	)
}
