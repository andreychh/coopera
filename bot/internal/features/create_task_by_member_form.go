package features

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms/actions"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

func CreateTaskByMemberTitleSpec(bot tg.Bot, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecCreateTaskByMemberFormTitle,
		hsm.CoreBehavior(
			base.SendContent(
				bot,
				sources.Static(content.Text("Please provide the title of the task.")),
			),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^[^\n]{3,64}$`)),
						base.SendContent(bot, sources.Static(content.Text(
							"Invalid title. Please use 3 to 64 characters without new lines.",
						))),
						hsm.Stay(),
					),
					hsm.Try(
						routing.Terminal(actions.SaveTextToField(f, "task_title")),
						hsm.Transit(SpecCreateTaskByMemberFormDescription),
					),
				),
			),
			composition.Nothing(),
		),
	)
}

func CreateTaskByMemberDescriptionSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecCreateTaskByMemberFormDescription,
		hsm.CoreBehavior(
			base.SendContent(bot, sources.Static[content.Content](
				keyboards.Resized(keyboards.Reply(
					content.Text("Please provide the description of the task."),
					buttons.Matrix(buttons.Row(buttons.TextButton("(Без описания)"))),
				)),
			)),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`(?s)^.{1,1000}$`)),
						base.SendContent(bot, sources.Static(content.Text(
							"Description is too long. Please keep it under 1000 characters.",
						))),
						hsm.Stay(),
					),
					hsm.Try(
						routing.Terminal(
							composition.Sequential(
								actions.SaveTextToField(f, "task_description"),
								domainactions.CreateDraft(c, f),
								base.SendContent(bot, sources.Static[content.Content](
									keyboards.Empty(content.Text("Task created successfully!")),
								)),
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

func CreateTaskByMemberSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Node(
		SpecCreateTaskByMemberForm,
		hsm.CoreBehavior(
			composition.Sequential(
				base.EditOrSendContent(
					bot,
					sources.Static(content.Text("Please fill out the form below or use /cancel to exit the form.")),
				),
				sources.Apply(
					forms.CurrentField(f, "team_id"),
					TeamIDFromCallbackData(),
					func(ctx context.Context, field forms.Field, value string) error {
						return field.ChangeValue(ctx, value)
					},
				),
			),
			hsm.Greedy(
				hsm.If(
					conditions.CommandIs("cancel"),
					hsm.Try(
						routing.Terminal(
							base.SendContent(bot, sources.Static[content.Content](
								keyboards.Empty(content.Text("Form canceled."))),
							),
						),
						hsm.Transit(SpecTeamsMenu),
					),
				),
			),
			composition.Nothing(),
		),
		hsm.Group(
			CreateTaskByMemberTitleSpec(bot, f),
			CreateTaskByMemberDescriptionSpec(bot, c, f),
		),
	)
}
