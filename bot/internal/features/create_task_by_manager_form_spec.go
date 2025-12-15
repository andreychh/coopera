package features

import (
	"context"
	"strings"

	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	domainconditions "github.com/andreychh/coopera-bot/internal/domain/conditions"
	"github.com/andreychh/coopera-bot/internal/ui/views"
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
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

func CreateTaskByManagerTitleSpec(bot tg.Bot, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecCreateTaskByManagerFormTitle,
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
						hsm.Transit(SpecCreateTaskByManagerFormDescription),
					),
				),
			),
			composition.Nothing(),
		),
	)
}

func CreateTaskByManagerDescriptionSpec(bot tg.Bot, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecCreateTaskByManagerFormDescription,
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
						routing.Terminal(actions.SaveTextToField(f, "task_description")),
						hsm.Transit(SpecCreateTaskByManagerFormPoints),
					),
				),
			),
			composition.Nothing(),
		),
	)
}

func CreateTaskByManagerPointsSpec(bot tg.Bot, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecCreateTaskByManagerFormPoints,
		hsm.CoreBehavior(
			base.SendContent(bot, sources.Static[content.Content](
				keyboards.Empty(content.Text(
					"Please provide the number of points for the task (1-99).",
				)),
			)),
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
						routing.Terminal(actions.SaveTextToField(f, "task_points")),
						hsm.Transit(SpecCreateTaskByManagerFormAssignTo),
					),
				),
			),
			composition.Nothing(),
		),
	)
}

func CreateTaskByManagerAssignToSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecCreateTaskByManagerFormAssignTo,
		hsm.CoreBehavior(
			base.SendContent(bot, views.MembersMatrixView(c, f)),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						conditions.TextIs("(unassigned)"),
						composition.Sequential(
							domainactions.CreateUnassigned(c, f),
							base.SendContent(bot, sources.Static[content.Content](
								keyboards.Empty(content.Text("Task created successfully!")),
							)),
						),
						hsm.Transit(SpecTeamsMenu),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@`)),
						base.SendContent(bot,
							sources.Static(content.Text("Invalid format. Username must start with '@' symbol.")),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(domainconditions.UserInTeam(c, f)),
						base.SendContent(bot, sources.Static(content.Text("User is not a member of this team."))),
						hsm.Stay(),
					),
					hsm.Try(
						routing.Terminal(
							composition.Sequential(
								sources.Apply(
									sources.Required(attributes.Text()),
									forms.CurrentField(f, "task_assignee"),
									func(ctx context.Context, text string, field forms.Field) error {
										username := strings.TrimPrefix(text, "@")
										return field.ChangeValue(ctx, username)
									},
								),
								domainactions.CreateAssigned(c, f),
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

func CreateTaskByManagerSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Node(
		SpecCreateTaskByManagerForm,
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
			CreateTaskByManagerTitleSpec(bot, f),
			CreateTaskByManagerDescriptionSpec(bot, f),
			CreateTaskByManagerPointsSpec(bot, f),
			CreateTaskByManagerAssignToSpec(bot, c, f),
		),
	)
}
