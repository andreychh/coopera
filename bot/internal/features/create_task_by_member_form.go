package features

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
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
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

func CreateTaskByMemberTitleSpec(bot tg.Bot, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecCreateTaskByMemberFormTitle,
		hsm.CoreBehavior(
			base.SendContent(
				bot,
				sources.Static(formatting.Formatted(
					content.Text("<b>Шаг 1 из 2: Название</b>\n\nВведите название задачи (от 3 до 64 символов)."),
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
					formatting.Formatted(
						content.Text("<b>Шаг 2 из 2: Описание</b>\n\nПодробно опишите задачу (до 1000 символов).\nЕсли описание не требуется, нажмите кнопку ниже."),
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
						routing.Terminal(
							composition.Sequential(
								actions.SaveTextToField(f, "task_description"),
								domainactions.CreateDraft(c, f),
								base.SendContent(bot, sources.Static[content.Content](
									keyboards.Empty(formatting.Formatted(
										content.Text("✅ <b>Задача успешно создана!</b>\nОна добавлена на <b>Доску задач</b>."),
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

func CreateTaskByMemberSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Node(
		SpecCreateTaskByMemberForm,
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
			CreateTaskByMemberTitleSpec(bot, f),
			CreateTaskByMemberDescriptionSpec(bot, c, f),
		),
	)
}
