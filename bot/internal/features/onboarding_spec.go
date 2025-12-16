package features

import (
	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/actions"
	"github.com/andreychh/coopera-bot/internal/ui/views"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
)

func OnboardingSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(
		SpecOnboarding,
		hsm.CoreBehavior(
			composition.Sequential(
				actions.CreateUser(domain.IdempotencyCommunity(c)),
				base.SendContent(bot, views.WelcomeMessage()),
			),
			hsm.Just(hsm.Transit(SpecMainMenu)),
			composition.Nothing(),
		),
	)
}
