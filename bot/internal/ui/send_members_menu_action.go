package ui

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/menu"
	"github.com/andreychh/coopera-bot/pkg/botlib/base/bot"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendMembersMenuAction struct {
	bot  bot.Bot
	team domain.Team
}

func (s sendMembersMenuAction) Perform(ctx context.Context, update telegram.Update) error {
	id, err := attributes.ChatID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting chat ID: %w", s, err)
	}
	details, err := s.details(ctx)
	if err != nil {
		return fmt.Errorf("(%T) getting team members details: %w", s, err)
	}
	err = s.bot.Chat(id).Send(ctx, menu.MembersMenu(details))
	if err != nil {
		return fmt.Errorf("(%T) sending members menu to chat #%d: %w", s, id, err)
	}
	return nil
}

func (s sendMembersMenuAction) details(ctx context.Context) ([]domain.MemberDetails, error) {
	members, err := s.team.Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting team members: %w", err)
	}
	var details []domain.MemberDetails
	for _, member := range members {
		detail, err := member.Details(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting member details: %w", err)
		}
		details = append(details, detail)
	}
	return details, nil
}
