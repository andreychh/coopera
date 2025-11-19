package ui

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/menu"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendTeamMenuAction struct {
	bot       tg.Bot
	community domain.Community
}

func (s sendTeamMenuAction) Perform(ctx context.Context, update telegram.Update) error {
	callbackData, err := attributes.CallbackData(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting callback data: %w", s, err)
	}
	teamID, err := callbacks.PrefixedData(callbackData).Value("team_id")
	if err != nil {
		return fmt.Errorf("(%T) getting team ID from callback data: %w", s, err)
	}
	chatID, err := attributes.ChatID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting chat ID: %w", s, err)
	}
	messageID, err := attributes.MessageID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting message ID: %w", s, err)
	}
	tID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		return fmt.Errorf("(%T) parsing team ID %q: %w", s, teamID, err)
	}
	details, err := s.community.Team(tID).Details(ctx)
	if err != nil {
		return fmt.Errorf("(%T) getting team details: %w", s, err)
	}
	err = s.bot.Chat(chatID).Message(int64(messageID)).Edit(ctx, menu.TeamMenu(details))
	if err != nil {
		return fmt.Errorf("(%T) sending team menu to chat #%d: %w", s, chatID, err)
	}
	return nil
}

func SendTeamMenu(bot tg.Bot, c domain.Community) core.Action {
	return sendTeamMenuAction{
		bot:       bot,
		community: c,
	}
}
