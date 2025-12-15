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
	MenuMembers     = "members"
	MenuUserTasks   = "user_ts"
	MenuTeamTasks   = "team_ts"
	MenuMemberTasks = "member_ts"
	MenuUserTask    = "user_t"
	MenuMemberTask  = "member_t"
	MenuTeamTask    = "team_t"
	MenuUserStats   = "user_st"

	ActionSubmitTaskForReview = "submit"
	ActionApproveTask         = "approve"
	ActionAssignTaskToSelf    = "assign"

	prefixChangeMenu = "cm"

	keyMenuName = "mn"
	keyTeamID   = "tmid"
	keyTaskID   = "tsid"
	keyAction   = "act"
)

func ToMainMenu() string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuMain).
		String()
}

func ToTeamsMenu() string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuTeams).
		String()
}

func ToTeamMenu(teamID int64) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuTeam).
		With(keyTeamID, strconv.FormatInt(teamID, 10)).
		String()
}

func ToMembersMenu(teamID int64) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuMembers).
		With(keyTeamID, strconv.FormatInt(teamID, 10)).
		String()
}

func ToUserTasksMenu() string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuUserTasks).
		String()
}

func ToTeamTasksMenu(teamID int64) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuTeamTasks).
		With(keyTeamID, strconv.FormatInt(teamID, 10)).
		String()
}

func ToMemberTasksMenu(teamID int64) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuMemberTasks).
		With(keyTeamID, strconv.FormatInt(teamID, 10)).
		String()
}

func ToUserTaskMenu(taskID int64) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuUserTask).
		With(keyTaskID, strconv.FormatInt(taskID, 10)).
		String()
}

func ToMemberTaskMenu(taskID int64) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuMemberTask).
		With(keyTaskID, strconv.FormatInt(taskID, 10)).
		String()
}

func ToTeamTaskMenu(taskID int64) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuTeamTask).
		With(keyTaskID, strconv.FormatInt(taskID, 10)).
		String()
}

func ToTeamTaskMenuWithAction(taskID int64, action string) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuTeamTask).
		With(keyTaskID, strconv.FormatInt(taskID, 10)).
		With(keyAction, action).
		String()
}

func ToUserStatsMenu() string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuUserStats).
		String()
}

func OnChangeMenu(id string) core.Condition {
	return composition.All(
		updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
		conditions.PrefixIs(prefixChangeMenu),
		conditions.ValueIs(keyMenuName, id),
	)
}

func OnAction(action string) core.Condition {
	return composition.All(
		updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
		conditions.ValueIs(keyAction, action),
	)
}

func ParseTeamID(callbackData string) (int64, error) {
	val, exists := callbacks.IncomingData(callbackData).Value(keyTeamID)
	if !exists {
		return 0, fmt.Errorf("parameter %q not found", keyTeamID)
	}
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %q: %w", keyTeamID, err)
	}
	return id, nil
}

func ParseTaskID(callbackData string) (int64, error) {
	val, exists := callbacks.IncomingData(callbackData).Value(keyTaskID)
	if !exists {
		return 0, fmt.Errorf("parameter %q not found", keyTaskID)
	}
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %q: %w", keyTaskID, err)
	}
	return id, nil
}
