package app

import (
	nethttp "net/http"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/http"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	dialogueskeyvalue "github.com/andreychh/coopera-bot/pkg/botlib/dialogues/keyvalue"
	"github.com/andreychh/coopera-bot/pkg/botlib/engine"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	formskeyvalue "github.com/andreychh/coopera-bot/pkg/botlib/forms/keyvalue"
	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	tgtransport "github.com/andreychh/coopera-bot/pkg/botlib/tg/transport"
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

func Client(token string) tgtransport.Client {
	return tgtransport.HTTPClient(token)
}

func Bot(client tgtransport.Client) tg.Bot {
	return tg.NewBot(client)
}

func Community(s string) domain.Community {
	return http.Community(transport.HTTPClient(s, &nethttp.Client{}))
}

func Engine(token string, clause core.Clause) engine.Engine {
	return engine.ShutdownEngine(
		engine.SingleWorkerEngine(
			token, clause,
		),
	)
}
