package features

import (
	"github.com/andreychh/coopera-bot/internal/domain"
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
			hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeams), hsm.Transit(SpecTeamsMenu)),
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
