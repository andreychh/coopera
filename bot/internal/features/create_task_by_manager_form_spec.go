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
	"github.com/andreychh/coopera-bot/pkg/botlib/content/formatting"
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
				sources.Static(formatting.Formatted(
					content.Text("<b>Шаг 1 из 4: Название</b>\n\nВведите название задачи (от 3 до 64 символов)."),
					formatting.ParseModeHTML,
				)),
			),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^[^\n]{3,64}$`)),
						base.SendContent(bot, sources.Static(formatting.Formatted(
							content.Text("<b>Ошибка:</b> Название должно быть от 3 до 64 символов и состоять из одной строки."),
							formatting.ParseModeHTML,
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
					formatting.Formatted(
						content.Text("<b>Шаг 2 из 4: Описание</b>\n\nПодробно опишите, что нужно сделать (до 1000 символов).\nЕсли описание не требуется, нажмите кнопку ниже."),
						formatting.ParseModeHTML,
					),
					buttons.Matrix(buttons.Row(buttons.TextButton("(Без описания)"))),
				)),
			)),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`(?s)^.{1,1000}$`)),
						base.SendContent(bot, sources.Static(formatting.Formatted(
							content.Text("<b>Ошибка:</b> Описание слишком длинное (максимум 1000 символов)."),
							formatting.ParseModeHTML,
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
				keyboards.Empty(formatting.Formatted(
					content.Text("<b>Шаг 3 из 4: Стоимость</b>\n\nОцените стоимость задачи в баллах (число от 1 до 99)."),
					formatting.ParseModeHTML,
				)),
			)),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^[1-9][0-9]?$`)),
						base.SendContent(bot, sources.Static(formatting.Formatted(
							content.Text("<b>Ошибка:</b> Введите целое число от 1 до 99."),
							formatting.ParseModeHTML,
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
						conditions.TextIs("(Без исполнителя)"),
						composition.Sequential(
							domainactions.CreateUnassigned(c, f),
							base.SendContent(bot, sources.Static[content.Content](
								keyboards.Empty(formatting.Formatted(
									content.Text("✅ <b>Задача успешно создана!</b>\nОна добавлена на <b>Доску задач</b>."),
									formatting.ParseModeHTML,
								)),
							)),
						),
						hsm.Transit(SpecTeamsMenu),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@`)),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Ошибка:</b> Юзернейм должен начинаться с символа '@' (или выберите из меню)."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(domainconditions.UserInTeam(c, f)),
						base.SendContent(bot, sources.Static(formatting.Formatted(
							content.Text("<b>Ошибка:</b> Пользователь не найден в этой команде."),
							formatting.ParseModeHTML,
						))),
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
									keyboards.Empty(formatting.Formatted(
										content.Text("✅ <b>Задача создана и назначена!</b>"),
										formatting.ParseModeHTML,
									)),
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
					sources.Static(formatting.Formatted(
						content.Text("<b>Создание задачи</b>\n\nЗаполните форму ниже. Для отмены используйте /cancel."),
						formatting.ParseModeHTML,
					)),
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
								keyboards.Empty(content.Text("Создание задачи отменено.")),
							)),
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
