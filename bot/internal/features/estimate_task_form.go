package features

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms/actions"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func EstimateTaskPointsSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecEstimateTaskFormPoints,
		hsm.CoreBehavior(
			composition.Nothing(),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^[1-9][0-9]?$`)),
						base.SendContent(bot, sources.Static(content.Text(
							"Invalid points. Please enter a number between 1 and 99.",
						))),
						hsm.Stay(),
					),
					hsm.Try(
						routing.Terminal(
							composition.Sequential(
								actions.SaveTextToField(f, "task_points"),
								domainactions.EstimateTask(c, f),
								base.SendContent(bot,
									sources.Static(
										content.Text("Task Updated successfully!"),
									),
								),
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

func EstimateTaskSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Node(
		SpecEstimateTaskForm,
		hsm.CoreBehavior(
			composition.Sequential(
				base.EditOrSendContent(
					bot,
					sources.Static(content.Text("Please provide a points of the task or use /cancel to exit the form.")),
				),
				sources.Apply(
					forms.CurrentField(f, "task_id"),
					TaskIDFromCallbackData(),
					func(ctx context.Context, field forms.Field, value string) error {
						return field.ChangeValue(ctx, value)
					},
				),
			),
			hsm.Greedy(
				hsm.JustIf(conditions.CommandIs("cancel"), hsm.Transit(SpecTeamsMenu)),
			),
			composition.Nothing(),
		),
		hsm.Group(EstimateTaskPointsSpec(bot, c, f)),
	)
}

type taskIDFromCallbackDataSource struct{}

func (t taskIDFromCallbackDataSource) Value(ctx context.Context, update telegram.Update) (string, error) {
	callbackData, exists := attributes.CallbackData().Value(update)
	if !exists {
		return "", fmt.Errorf("getting callback data from update: callback data not found")
	}
	id, err := protocol.ParseTaskID(callbackData)
	if err != nil {
		return "", fmt.Errorf("parsing team ID from callback data %q: %w", callbackData, err)
	}
	return strconv.FormatInt(id, 10), nil
}

func TaskIDFromCallbackData() sources.Source[string] {
	return taskIDFromCallbackDataSource{}
}
