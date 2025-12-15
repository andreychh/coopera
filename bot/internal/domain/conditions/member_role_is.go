package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type memberRoleIsCondition struct {
	community domain.Community
	expected  domain.MemberRole
}

func (m memberRoleIsCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return false, fmt.Errorf("chat ID not found in update")
	}
	user, exists, err := m.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return false, fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	if !exists {
		return false, fmt.Errorf("user with telegram ID %d does not exist", chatID)
	}
	callback, found := attributes.CallbackData().Value(update)
	if !found {
		return false, fmt.Errorf("callback data not found in update")
	}
	teamID, err := protocol.ParseTeamID(callback)
	if err != nil {
		return false, fmt.Errorf("parsing team ID from callback data %q: %w", callback, err)
	}
	team, exists, err := m.community.Team(ctx, teamID)
	if err != nil {
		return false, fmt.Errorf("getting team %d: %w", teamID, err)
	}
	if !exists {
		return false, fmt.Errorf("team %d does not exist", teamID)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return false, fmt.Errorf("getting members of team %d: %w", teamID, err)
	}
	member, exists, err := members.MemberWithUsername(ctx, user.Username())
	if err != nil {
		return false, fmt.Errorf("getting member with username %q in team %d: %w", user.Username(), teamID, err)
	}
	if !exists {
		return false, fmt.Errorf("member with username %q does not exist in team %d", user.Username(), teamID)
	}
	return member.Role() == m.expected, nil
}

func MemberRoleIs(community domain.Community, expected domain.MemberRole) core.Condition {
	return memberRoleIsCondition{
		community: community,
		expected:  expected,
	}
}
