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
	MenuMain  = "main"
	MenuTeams = "teams"
)

var Navigation = navigationProtocol{}

type navigationProtocol struct{}

const (
	navPrefix = "change_menu"
	navKey    = "menu_name"
)

// Payload генерирует строку для КНОПКИ
// Пример использования: protocol.Navigation.Payload(protocol.MenuTeams)
func (navigationProtocol) Payload(menuID string) string {
	return callbacks.OutcomingData(navPrefix).
		With(navKey, menuID).
		String()
}

// On генерирует условие для РОУТЕРА (Tree)
// Пример использования: protocol.Navigation.On(protocol.MenuTeams)
func (navigationProtocol) On(menuID string) core.Condition {
	return composition.All(
		updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
		conditions.PrefixIs(navPrefix),
		conditions.ValueIs(navKey, menuID),
	)
}

const (
	MenuTeam        = "team"
	ParamTeamID     = "team_id"
	MenuTeamMembers = "team_members"
)

func (navigationProtocol) ToTeam(teamID int64) string {
	return callbacks.OutcomingData(navPrefix).
		With(navKey, MenuTeam).
		With(ParamTeamID, strconv.FormatInt(teamID, 10)).
		String()
}

func (navigationProtocol) ParseTeamID(callbackData string) (int64, error) {
	val, exists := callbacks.IncomingData(callbackData).Value(ParamTeamID)
	if !exists {
		return 0, fmt.Errorf("parameter %s not found", ParamTeamID)
	}
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s: %w", ParamTeamID, err)
	}
	return id, nil
}

func (navigationProtocol) ToTeamMembers(teamID int64) string {
	return callbacks.OutcomingData(navPrefix).
		With(navKey, MenuTeamMembers).
		With(ParamTeamID, strconv.FormatInt(teamID, 10)).
		String()
}
