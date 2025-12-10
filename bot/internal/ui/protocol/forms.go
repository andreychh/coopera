package protocol

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks/conditions"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	updatesconditions "github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

const (
	FormCreateTeam = "create_team"

	prefixStartForm = "start_form"

	keyFormName = "form_name"
)

func OnStartForm(name string) core.Condition {
	return composition.All(
		updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
		conditions.PrefixIs(prefixStartForm),
		conditions.ValueIs(keyFormName, name),
	)
}

func StartForm(name string) string {
	return callbacks.OutcomingData(prefixStartForm).
		With(keyFormName, name).
		String()
}
