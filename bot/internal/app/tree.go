package app

import (
	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/internal/ui/views"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms/actions"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	updcond "github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

func OnboardingBehavior(bot tg.Bot, c domain.Community) hsm.Behavior {
	return hsm.CoreBehavior(
		composition.Sequential(
			domainactions.CreateUser(domain.IdempotencyCommunity(c)),
			base.SendContent(bot, views.WelcomeMessage()),
		),
		hsm.Just(hsm.Transit("main_menu")),
		composition.Nothing(),
	)
}

func MainMenuBehavior(bot tg.Bot) hsm.Behavior {
	return hsm.CoreBehavior(
		base.EditOrSendContent(bot, views.MainMenu()),
		hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeams), hsm.Transit("teams_menu")),
		composition.Nothing(),
	)
}

func TeamsMenuBehavior(bot tg.Bot, c domain.Community) hsm.Behavior {
	return hsm.CoreBehavior(
		base.EditOrSendContent(bot, views.TeamsMenu(c)),
		hsm.FirstHandled(
			hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit("team_menu")),
			hsm.JustIf(protocol.OnStartForm(protocol.FormCreateTeam), hsm.Transit("create_team_form")),
			hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMain), hsm.Transit("main_menu")),
		),
		composition.Nothing(),
	)
}

func TeamMenuBehavior(bot tg.Bot, c domain.Community) hsm.Behavior {
	return hsm.CoreBehavior(
		base.EditOrSendContent(bot, views.TeamMenu(c)),
		hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit("team_menu")),
		composition.Nothing(),
	)
}

func CreateTeamFormBehavior() hsm.Behavior {
	return hsm.CoreBehavior(
		composition.Nothing(),
		hsm.Greedy(
			hsm.JustIf(updcond.CommandIs("cancel"), hsm.Transit("teams_menu")),
		),
		composition.Nothing(),
	)
}

func CreateTeamFormTeamNameBehavior(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Behavior {
	return hsm.CoreBehavior(
		base.EditOrSendContent(bot, content.StaticView(content.Text("Please provide the name of your team."))),
		hsm.If(
			composition.Not(updcond.CommandIs("cancel")),
			hsm.FirstHandled(
				hsm.TryAction(
					composition.Not(updcond.TextMatchesRegexp("^[A-Za-zА-Яа-я0-9_ -]{3,50}$")),
					base.SendContent(bot,
						content.StaticView(
							content.Text("Please provide the name of your team using 3 to 50 characters: letters, numbers, spaces, hyphens, or underscores."),
						),
					),
					hsm.Stay(),
				),
				hsm.Try(routing.Terminal(
					composition.Sequential(
						actions.SaveTextToField(f, "team_name"),
						domainactions.CreateTeam(f, c),
					)),
					hsm.Transit("teams_menu"),
				),
			),
		),
		composition.Nothing(),
	)
}

func Tree(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Node(
		"root",
		hsm.CoreBehavior(
			composition.Nothing(),
			hsm.JustIf(
				composition.All(
					updcond.UpdateTypeIs(updates.UpdateTypeMessage),
					updcond.CommandIs("start"),
				),
				hsm.Transit("onboarding")),
			composition.Nothing(),
		),
		hsm.Group(
			hsm.Leaf("onboarding", OnboardingBehavior(bot, c)),
			hsm.Leaf("main_menu", MainMenuBehavior(bot)),
			hsm.Leaf("teams_menu", TeamsMenuBehavior(bot, c)),
			hsm.Leaf("team_menu", TeamMenuBehavior(bot, c)),
			hsm.Node(
				"create_team_form",
				CreateTeamFormBehavior(),
				hsm.Group(
					hsm.Leaf("create_team_form:team_name", CreateTeamFormTeamNameBehavior(bot, c, f)),
				),
			),
		),
	)
}
