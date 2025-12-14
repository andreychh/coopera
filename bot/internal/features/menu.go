package features

import (
	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
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
				hsm.JustIf(
					protocol.OnChangeMenu(protocol.MenuTasksAssignedToUser),
					hsm.Transit(SpecTasksAssignedToUser),
				),
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
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuAllTeamTasks), hsm.Transit(SpecAllTeamTasks)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMemberTasks), hsm.Transit(SpecMemberTasks)),
				hsm.JustIf(protocol.OnStartForm(protocol.FormCreateTask), hsm.Transit(SpecCreateTaskForm)),
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
		SpecTasksAssignedToUser,
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

func AllTeamTasksSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecAllTeamTasks,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.AllTeamTasks(c)),
			hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit(SpecTeamMenu)),
			composition.Nothing(),
		),
	)
}

func MemberTasksSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecMemberTasks,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.MemberTasksMenuView(c)),
			hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit(SpecTeamMenu)),
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
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTasksAssignedToUser), hsm.Transit(SpecTasksAssignedToUser)),
				hsm.TryAction(
					protocol.OnChangeMenu(protocol.MenuUserTask),
					domainactions.MarkTaskAsCompleted(c),
					hsm.Transit(SpecUserTask),
				),
			),
			composition.Nothing(),
		),
	)
}
