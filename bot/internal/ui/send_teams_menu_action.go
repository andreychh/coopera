package ui

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/menu"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendTeamsMenuAction struct {
	bot       tg.Bot
	community domain.Community
}

func (s sendTeamsMenuAction) Perform(ctx context.Context, update telegram.Update) error {
	chatID, err := attributes.ChatID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting chat ID: %w", s, err)
	}
	messageID, err := attributes.MessageID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting message ID: %w", s, err)
	}
	details, err := s.details(ctx, chatID)
	if err != nil {
		return err
	}
	err = s.bot.Chat(chatID).Message(int64(messageID)).Edit(ctx, menu.TeamsMenu(details))
	if err != nil {
		return err
	}
	return nil
}

func (s sendTeamsMenuAction) details(ctx context.Context, id int64) ([]domain.TeamDetails, error) {
	teams, err := s.community.UserWithTelegramID(id).CreatedTeams(ctx)
	if err != nil {
		return nil, err
	}
	var details []domain.TeamDetails
	for _, team := range teams {
		detail, err := team.Details(ctx)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	return details, nil
}

func SendTeamsMenu(bot tg.Bot, community domain.Community) core.Action {
	return sendTeamsMenuAction{
		bot:       bot,
		community: community,
	}
}
