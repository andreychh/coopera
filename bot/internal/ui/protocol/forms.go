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
)

const (
	formPrefix = "start_form"
	formKey    = "form_name"
)

func OnStartForm(formID string) core.Condition {
	return composition.All(
		updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
		conditions.PrefixIs(formPrefix),
		conditions.ValueIs(formKey, formID),
	)
}

func StartFromPayload(formID string) string {
	return callbacks.OutcomingData(formPrefix).
		With(formKey, formID).
		String()
}
