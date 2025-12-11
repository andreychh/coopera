package features

import (
	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	domainconditions "github.com/andreychh/coopera-bot/internal/domain/conditions"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms/actions"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

func CreateTeamFormTeamNameSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecCreateTeamFormTeamName,
		hsm.CoreBehavior(
			base.SendContent(
				bot,
				sources.Static(content.Text("Please provide the name of your team.")),
			),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp("^[A-Za-zА-Яа-я0-9_ -]{3,50}$")),
						base.SendContent(bot,
							sources.Static(
								content.Text("Please provide the name of your team using 3 to 50 characters: letters, numbers, spaces, hyphens, or underscores."),
							),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						domainconditions.TeamExists(c),
						base.SendContent(bot,
							sources.Static(
								content.Text("Team with this name already exists. Please choose a different name."),
							),
						),
						hsm.Stay(),
					),
					hsm.Try(
						routing.Terminal(
							composition.Sequential(
								actions.SaveTextToField(f, "team_name"),
								domainactions.CreateTeam(c, f),
							),
						),
						hsm.Transit(SpecTeamsMenu),
					),
				),
			),
			composition.Nothing(),
		),
	)
}

func CreateTeamFormSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Node(
		SpecCreateTeamForm,
		hsm.CoreBehavior(
			base.EditOrSendContent(
				bot,
				sources.Static(content.Text("Fill out the form below or use /cancel to exit the form.")),
			),
			hsm.Greedy(
				hsm.JustIf(conditions.CommandIs("cancel"), hsm.Transit(SpecTeamsMenu)),
			),
			composition.Nothing(),
		),
		hsm.Group(CreateTeamFormTeamNameSpec(bot, c, f)),
	)
}
