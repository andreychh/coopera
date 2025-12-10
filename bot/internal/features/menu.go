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
		protocol.MenuMain,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.MainMenuView()),
			hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeams), hsm.Transit(protocol.MenuTeams)),
			composition.Nothing(),
		),
	)
}

func TeamMenuSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		protocol.MenuTeam,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.TeamMenu(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMembers), hsm.Transit(protocol.MenuMembers)),
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeams), hsm.Transit(protocol.MenuTeams)),
			),
			composition.Nothing(),
		),
	)
}

func MembersMenuSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		protocol.MenuMembers,
		hsm.CoreBehavior(
			base.EditOrSendContent(bot, views.MembersMenu(c)),
			hsm.FirstHandled(
				hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit(protocol.MenuTeam)),
			),
			composition.Nothing(),
		),
	)
}
