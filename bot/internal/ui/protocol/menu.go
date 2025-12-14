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
	MenuMain                = "main_menu"
	MenuTeams               = "teams_menu"
	MenuTeam                = "team_menu"
	MenuMembers             = "members_menu"
	MenuMember              = "member_menu"
	MenuTasksAssignedToUser = "tasks_assigned_to_user_menu"
	MenuAllTeamTasks        = "all_team_tasks_menu"
	MenuMemberTasks         = "member_tasks_menu"
	MenuUserTask            = "user_task_menu"

	prefixChangeMenu = "change_menu"

	keyMenuName = "menu_name"
	keyTeamID   = "team_id"
	keyMemberID = "member_id"
	keyTaskID   = "task_id"
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

func ToMemberMenu(memberID int64) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuMember).
		With(keyMemberID, strconv.FormatInt(memberID, 10)).
		String()
}

func ToTasksAssignedToUserMenu() string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuTasksAssignedToUser).
		String()
}

func ToAllTeamTasksMenu(teamID int64) string {
	return callbacks.OutcomingData(prefixChangeMenu).
		With(keyMenuName, MenuAllTeamTasks).
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

func OnChangeMenu(id string) core.Condition {
	return composition.All(
		updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
		conditions.PrefixIs(prefixChangeMenu),
		conditions.ValueIs(keyMenuName, id),
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
