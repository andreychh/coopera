package features

import (
	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	domainconditions "github.com/andreychh/coopera-bot/internal/domain/conditions"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/formatting"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms/actions"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

func CreateTeamFormNameSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecCreateTeamFormName,
		hsm.CoreBehavior(
			base.SendContent(
				bot,
				sources.Static(formatting.Formatted(
					content.Text("<b>Название команды</b>\n\nПридумайте название для вашей команды (от 3 до 50 символов)."),
					formatting.ParseModeHTML,
				)),
			),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp("^[A-Za-zА-Яа-я0-9_ -]{3,50}$")),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Ошибка формата:</b> Используйте от 3 до 50 символов: буквы, цифры, пробелы, дефис (-) или подчеркивание (_)."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						domainconditions.TeamExists(c),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Ошибка:</b> Команда с таким названием уже существует. Пожалуйста, придумайте другое."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.Try(
						routing.Terminal(
							composition.Sequential(
								actions.SaveTextToField(f, "team_name"),
								domainactions.CreateTeam(c, f),
								base.SendContent(bot,
									sources.Static(formatting.Formatted(
										content.Text("<b>Успешно:</b> Команда создана! Теперь вы можете добавить в неё участников."),
										formatting.ParseModeHTML,
									)),
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

func CreateTeamFormSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Node(
		SpecCreateTeamForm,
		hsm.CoreBehavior(
			base.EditOrSendContent(
				bot,
				sources.Static(formatting.Formatted(
					content.Text("<b>Создание новой команды</b>\n\nЗаполните форму ниже. Для отмены используйте /cancel."),
					formatting.ParseModeHTML,
				)),
			),
			hsm.Greedy(
				hsm.If(
					conditions.CommandIs("cancel"),
					hsm.Try(
						routing.Terminal(
							base.SendContent(bot, sources.Static[content.Content](
								keyboards.Empty(content.Text("Создание команды отменено.")),
							)),
						),
						hsm.Transit(SpecTeamsMenu),
					),
				),
			),
			composition.Nothing(),
		),
		hsm.Group(CreateTeamFormNameSpec(bot, c, f)),
	)
}
