package features

import (
	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	domainconditions "github.com/andreychh/coopera-bot/internal/domain/conditions"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/internal/ui/views"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
)

func MainMenuSpec(bot tg.Bot) hsm.Spec {
	return hsm.Leaf(
		SpecMainMenu,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.MainMenuView()),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeams), hsm.Transit(SpecTeamsMenu)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuUserTasks), hsm.Transit(SpecUserTasks)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuUserStats), hsm.Transit(SpecUserStats)),
			),
			composition.Nothing(),
		),
	)
}

func TeamsMenuSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecTeamsMenu,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.TeamsMenu(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit(SpecTeamMenu)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMain), hsm.Transit(SpecMainMenu)),
				hsm.JustIf(protocol.OnStartForm(protocol.FormCreateTeam), hsm.Transit(SpecCreateTeamForm)),
			),
			composition.Nothing(),
		),
	)
}

func TeamMenuSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecTeamMenu,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.TeamMenu(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMembers), hsm.Transit(SpecMembersMenu)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeams), hsm.Transit(SpecTeamsMenu)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeamTasks), hsm.Transit(SpecTeamTasks)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMemberTasks), hsm.Transit(SpecMemberTasks)),
				hsm.JustIf(
					composition.All(
						protocol.OnStartForm(protocol.FormCreateTask),
						domainconditions.MemberRoleIs(c, domain.RoleManager),
					),
					hsm.Transit(SpecCreateTaskByManagerForm),
				),
				hsm.JustIf(
					composition.All(
						protocol.OnStartForm(protocol.FormCreateTask),
						domainconditions.MemberRoleIs(c, domain.RoleMember),
					),
					hsm.Transit(SpecCreateTaskByMemberForm),
				),
			),
			composition.Nothing(),
		),
	)
}

func MembersMenuSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecMembersMenu,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.MembersMenu(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit(SpecTeamMenu)),
				hsm.JustIf(protocol.OnStartForm(protocol.FormAddMember), hsm.Transit(SpecAddMemberForm)),
			),
			composition.Nothing(),
		),
	)
}

func TasksAssignedToUserSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecUserTasks,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.TasksAssignedToUserView(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMain), hsm.Transit(SpecMainMenu)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuUserTask), hsm.Transit(SpecUserTask)),
			),
			composition.Nothing(),
		),
	)
}

func TeamTasksSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecTeamTasks,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.TeamTasks(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit(SpecTeamMenu)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeamTask), hsm.Transit(SpecTeamTask)),
			),
			composition.Nothing(),
		),
	)
}

func MemberTasksSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecMemberTasks,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.MemberTasksMenuView(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit(SpecTeamMenu)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMemberTask), hsm.Transit(SpecMemberTask)),
			),
			composition.Nothing(),
		),
	)
}

func UserTaskSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecUserTask,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.UserTaskMenuView(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuUserTasks), hsm.Transit(SpecUserTasks)),
				hsm.TryAction(
					protocol.OnChangeMenu(protocol.MenuUserTask),
					domainactions.SubmitTaskForReview(c),
					hsm.Transit(SpecUserTask),
				),
			),
			composition.Nothing(),
		),
	)
}

func TeamTaskSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecTeamTask,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.TeamTaskMenuView(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeamTasks), hsm.Transit(SpecTeamTasks)),
				hsm.JustIf(
					protocol.OnStartForm(protocol.FormEstimateTask),
					hsm.Transit(SpecEstimateTaskForm),
				),
				hsm.TryAction(
					composition.All(
						protocol.OnChangeMenu(protocol.MenuTeamTask),
						protocol.OnAction(protocol.ActionAssignTaskToSelf),
					),
					domainactions.AssignTaskToSelf(c),
					hsm.Transit(SpecTeamTask),
				),
				hsm.TryAction(
					composition.All(
						protocol.OnChangeMenu(protocol.MenuTeamTask),
						protocol.OnAction(protocol.ActionSubmitTaskForReview),
					),
					domainactions.SubmitTaskForReview(c),
					hsm.Transit(SpecTeamTask),
				),
				hsm.TryAction(
					composition.All(
						protocol.OnChangeMenu(protocol.MenuTeamTask),
						protocol.OnAction(protocol.ActionApproveTask),
					),
					domainactions.ApproveTask(c),
					hsm.Transit(SpecTeamTask),
				),
			),
			composition.Nothing(),
		),
	)
}

func MemberTaskSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecMemberTask,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.MemberTaskMenuView(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMemberTasks), hsm.Transit(SpecMemberTasks)),
				hsm.TryAction(
					protocol.OnChangeMenu(protocol.MenuMemberTask),
					domainactions.SubmitTaskForReview(c),
					hsm.Transit(SpecMemberTask),
				),
			),
			composition.Nothing(),
		),
	)
}

func UserStatsMenuSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecUserStats,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.UserStatsMenu(c)),
			hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMain), hsm.Transit(SpecMainMenu)),
			composition.Nothing(),
		),
	)
}
