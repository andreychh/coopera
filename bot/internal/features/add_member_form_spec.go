package features

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	domainconditions "github.com/andreychh/coopera-bot/internal/domain/conditions"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
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
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddMemberUsernameSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Leaf(
		SpecAddMemberFormUsername,
		hsm.CoreBehavior(
			base.SendContent(
				bot,
				sources.Static(
					formatting.Formatted(
						content.Text("Введите Telegram-юзернейм пользователя (в формате <b>@username</b>):"),
						formatting.ParseModeHTML,
					),
				),
			),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@`)),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Ошибка формата:</b> Юзернейм должен начинаться с символа '@'."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@.{5,32}$`)),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Ошибка длины:</b> Юзернейм должен содержать от 5 до 32 символов."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@[a-zA-Z0-9_]+$`)),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Недопустимые символы:</b> Используйте только латинские буквы (a-z), цифры (0-9) и нижнее подчеркивание (_)."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@[a-zA-Z]`)),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Ошибка формата:</b> Юзернейм должен начинаться с буквы."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`[a-zA-Z0-9]$`)),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Ошибка формата:</b> Юзернейм не может заканчиваться подчеркиванием."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						conditions.TextMatchesRegexp(`__`),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Ошибка формата:</b> Двойное подчеркивание (__) запрещено."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(domainconditions.UserExists(c)),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Пользователь не найден.</b>\nПользователь должен сначала запустить этого бота (/start), чтобы вы могли добавить его в команду."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						domainconditions.UserInTeam(c, f),
						base.SendContent(bot,
							sources.Static(formatting.Formatted(
								content.Text("<b>Ошибка:</b> Этот пользователь уже состоит в данной команде."),
								formatting.ParseModeHTML,
							)),
						),
						hsm.Stay(),
					),
					hsm.Try(routing.Terminal(
						composition.Sequential(
							actions.SaveTextToField(f, "member_username"),
							domainactions.AddMember(c, f),
							base.SendContent(bot,
								sources.Static(formatting.Formatted(
									content.Text("<b>Успешно:</b> Пользователь добавлен в команду."),
									formatting.ParseModeHTML,
								)),
							),
						)),
						hsm.Transit(SpecTeamsMenu),
					),
				),
			),
			composition.Nothing(),
		),
	)
}

func AddMemberSpec(bot tg.Bot, c domain.Community, f forms.Forms) hsm.Spec {
	return hsm.Node(
		SpecAddMemberForm,
		hsm.CoreBehavior(
			composition.Sequential(
				base.EditOrSendContent(
					bot,
					sources.Static(formatting.Formatted(
						content.Text("<b>Добавление участника</b>\n\nЗаполните форму ниже. Для отмены используйте команду /cancel."),
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
								keyboards.Empty(content.Text("Действие отменено.")),
							)),
						),
						hsm.Transit(SpecTeamsMenu),
					),
				),
			),
			composition.Nothing(),
		),
		hsm.Group(AddMemberUsernameSpec(bot, c, f)),
	)
}

type teamIDFromCallbackDataSource struct{}

func (t teamIDFromCallbackDataSource) Value(ctx context.Context, update telegram.Update) (string, error) {
	callbackData, exists := attributes.CallbackData().Value(update)
	if !exists {
		return "", fmt.Errorf("getting callback data from update: callback data not found")
	}
	id, err := protocol.ParseTeamID(callbackData)
	if err != nil {
		return "", fmt.Errorf("parsing team ID from callback data %q: %w", callbackData, err)
	}
	return strconv.FormatInt(id, 10), nil
}

func TeamIDFromCallbackData() sources.Source[string] {
	return teamIDFromCallbackDataSource{}
}
