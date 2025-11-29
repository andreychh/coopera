package protocol

import (
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks/conditions"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	updatesconditions "github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

const (
	MenuMain        = "main"
	MenuTeams       = "teams"
	MenuTeam        = "team"
	MenuTeamMembers = "team_members"
)

const (
	navPrefix = "change_menu"
	navKey    = "menu_name"
	teamIDKey = "team_id"
)

func ToMainMenu() string {
	return callbacks.OutcomingData(navPrefix).
		With(navKey, MenuMain).
		String()
}

func ToTeamsMenu() string {
	return callbacks.OutcomingData(navPrefix).
		With(navKey, MenuTeams).
		String()
}

func ToTeamMenu(teamID int64) string {
	return callbacks.OutcomingData(navPrefix).
		With(navKey, MenuTeam).
		With(teamIDKey, strconv.FormatInt(teamID, 10)).
		String()
}

func ToMembersMenu(teamID int64) string {
	return callbacks.OutcomingData(navPrefix).
		With(navKey, MenuTeamMembers).
		With(teamIDKey, strconv.FormatInt(teamID, 10)).
		String()
}

func OnChangeMenu(id string) core.Condition {
	return composition.All(
		updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
		conditions.PrefixIs(navPrefix),
		conditions.ValueIs(navKey, id),
	)
}

func ParseTeamID(callbackData string) (int64, error) {
	val, exists := callbacks.IncomingData(callbackData).Value(teamIDKey)
	if !exists {
		return 0, fmt.Errorf("parameter %s not found", teamIDKey)
	}
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s: %w", teamIDKey, err)
	}
	return id, nil
}
