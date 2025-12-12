package protocol

import (
	"strconv"

	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks/conditions"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	updatesconditions "github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

const (
	FormCreateTeam = "create_team_from"
	FormAddMember  = "add_member_form"

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

func StartCreateTeamForm() string {
	return callbacks.OutcomingData(prefixStartForm).
		With(keyFormName, FormCreateTeam).
		String()
}

func StartAddMemberForm(teamID int64) string {
	return callbacks.OutcomingData(prefixStartForm).
		With(keyFormName, FormAddMember).
		With(keyTeamID, strconv.FormatInt(teamID, 10)).
		String()
}
