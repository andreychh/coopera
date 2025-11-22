package protocol

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks/conditions"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	updatesconditions "github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

type Form string

const (
	FormCreateTeam = "create_team"
)

var Forms = formsProtocol{}

type formsProtocol struct{}

const (
	formPrefix = "start_form"
	formKey    = "form_name"
)

// StartPayload — строка для кнопки "Начать форму"
func (formsProtocol) StartPayload(formID string) string {
	return callbacks.OutcomingData(formPrefix).
		With(formKey, formID).
		String()
}

// OnStart — условие роутера "Пользователь нажал начать форму"
func (formsProtocol) OnStart(formID string) core.Condition {
	return composition.All(
		updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
		conditions.PrefixIs(formPrefix),
		conditions.ValueIs(formKey, formID),
	)
}
