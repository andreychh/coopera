package features

import (
	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	updcond "github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

func RootSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Node(
		SpecRoot,
		hsm.CoreBehavior(
			composition.Nothing(),
			hsm.JustIf(updcond.CommandIs("start"), hsm.Transit(SpecMainMenu)),
			composition.Nothing(),
		),
		hsm.Group(
			OnboardingSpec(bot, c),

			MainMenuSpec(bot),
			TeamsMenuSpec(bot, c),
			TeamMenuSpec(bot, c),
			MembersMenuSpec(bot, c),
			TasksAssignedToUserSpec(bot, c),
			AllTeamTasksSpec(bot, c),
			MemberTasksSpec(bot, c),

			CreateTeamFormSpec(bot, c, f),
			AddMemberSpec(bot, c, f),
			CreateTaskSpec(bot, c, f),
		),
	)
}
