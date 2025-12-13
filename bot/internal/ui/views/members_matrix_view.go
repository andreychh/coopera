package views

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
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
	teamID, err := m.forms.Form(chatID).Field("team_id").Value(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting team_id field for user %d: %w", chatID, err)
	}
	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing team_id %q to int64: %w", teamID, err)
	}
	team, err := m.community.Team(ctx, intTeamID)
	if err != nil {
		return nil, fmt.Errorf("getting team with ID %d: %w", intTeamID, err)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members of team %d: %w", intTeamID, err)
	}
	slice, err := members.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members slice for team %d: %w", intTeamID, err)
	}
	matrix := buttons.Matrix(buttons.Row(buttons.TextButton("(unassigned)")))
	for _, member := range slice {
		matrix = matrix.WithRow(buttons.Row(buttons.TextButton(fmt.Sprintf("@%s", member.Name()))))
	}
	return keyboards.Resized(keyboards.Reply(
		content.Text("Select a member to assign the task to:"),
		matrix,
	)), nil
}

func MembersMatrixView(c domain.Community, f forms.Forms) sources.Source[content.Content] {
	return membersMatrixView{
		community: c,
		forms:     f,
	}
}
