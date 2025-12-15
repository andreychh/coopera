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
				sources.Static(content.Text("Please provide the username of the member to add.")),
			),
			hsm.If(
				composition.Not(conditions.CommandIs("cancel")),
				hsm.FirstHandled(
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@`)),
						base.SendContent(bot,
							sources.Static(content.Text("Invalid format. Username must start with '@' symbol.")),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@.{5,32}$`)),
						base.SendContent(bot,
							sources.Static(content.Text("Incorrect length. Username must be between 5 and 32 characters.")),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@[a-zA-Z0-9_]+$`)),
						base.SendContent(bot,
							sources.Static(content.Text("Invalid characters. Only Latin letters (a-z), numbers (0-9), and underscores (_) are allowed.")),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`^@[a-zA-Z]`)),
						base.SendContent(bot,
							sources.Static(content.Text("Invalid start. Username must start with a letter.")),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(conditions.TextMatchesRegexp(`[a-zA-Z0-9]$`)),
						base.SendContent(bot,
							sources.Static(content.Text("Invalid ending. Username cannot end with an underscore.")),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						conditions.TextMatchesRegexp(`__`),
						base.SendContent(bot,
							sources.Static(content.Text("Invalid format. Consecutive underscores (__) are not allowed.")),
						),
						hsm.Stay(),
					),
					hsm.TryAction(
						composition.Not(domainconditions.UserExists(c)),
						base.SendContent(bot, sources.Static(content.Text("User not found in bot database. Ask them to /start the bot first."))),
						hsm.Stay(),
					),
					hsm.TryAction(
						domainconditions.UserInTeam(c, f),
						base.SendContent(bot, sources.Static(content.Text("User is already a member of this team."))),
						hsm.Stay(),
					),
					hsm.Try(routing.Terminal(
						composition.Sequential(
							actions.SaveTextToField(f, "member_username"),
							domainactions.AddMember(c, f),
							base.SendContent(bot,
								sources.Static(
									content.Text("Member added successfully!"),
								),
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
