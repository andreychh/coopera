package views

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/formatting"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type membersMatrixView struct {
	community domain.Community
	forms     forms.Forms
}

func (m membersMatrixView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return nil, fmt.Errorf("chat ID not found in update")
	}
	teamIDStr, err := m.forms.Form(chatID).Field("team_id").Value(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting team_id field for user %d: %w", chatID, err)
	}
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing team_id %q to int64: %w", teamIDStr, err)
	}
	team, exists, err := m.community.Team(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("getting team with ID %d: %w", teamID, err)
	}
	if !exists {
		return nil, fmt.Errorf("team with ID %d does not exist", teamID)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members of team %d: %w", teamID, err)
	}
	slice, err := members.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members slice for team %d: %w", teamID, err)
	}
	matrix := buttons.Matrix(buttons.Row(buttons.TextButton("(Без исполнителя)")))
	for _, member := range slice {
		matrix = matrix.WithRow(buttons.Row(buttons.TextButton(fmt.Sprintf("@%s", member.Username()))))
	}
	return keyboards.Resized(keyboards.Reply(
		formatting.Formatted(
			content.Text("<b>Шаг 4 из 4: Исполнитель</b>\n\nВыберите участника из списка, чтобы назначить задачу сразу.\nИли выберите пункт '(Без исполнителя)', чтобы задача попала в общую очередь."),
			formatting.ParseModeHTML,
		),
		matrix,
	)), nil
}

func MembersMatrixView(c domain.Community, f forms.Forms) sources.Source[content.Content] {
	return membersMatrixView{
		community: c,
		forms:     f,
	}
}
