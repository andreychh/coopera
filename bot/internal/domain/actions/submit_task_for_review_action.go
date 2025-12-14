package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type submitTaskForReviewAction struct {
	community domain.Community
}

func (m submitTaskForReviewAction) Perform(ctx context.Context, update telegram.Update) error {
	callback, found := attributes.CallbackData().Value(update)
	if !found {
		return fmt.Errorf("callback data not found in update")
	}
	taskID, err := protocol.ParseTaskID(callback)
	if err != nil {
		return fmt.Errorf("parsing task ID from callback data %q: %w", callback, err)
	}
	task, exists, err := m.community.Task(ctx, taskID)
	if err != nil {
		return fmt.Errorf("getting task %d: %w", taskID, err)
	}
	if !exists {
		return fmt.Errorf("task %d does not exist", taskID)
	}
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found in update")
	}
	user, exists, err := m.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	if !exists {
		return fmt.Errorf("user with telegram ID %d does not exist", chatID)
	}
	team, err := task.Team(ctx)
	if err != nil {
		return fmt.Errorf("getting team for task %d: %w", taskID, err)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return fmt.Errorf("getting members for team %d: %w", team.ID(), err)
	}
	member, exists, err := members.MemberWithUsername(ctx, user.Username())
	if err != nil {
		return fmt.Errorf("getting member with username %q in team %d: %w", user.Username(), team.ID(), err)
	}
	if !exists {
		return fmt.Errorf("member with username %q does not exist in team %d", user.Username(), team.ID())
	}
	err = member.SubmitTaskForReview(ctx, task.ID())
	if err != nil {
		return fmt.Errorf("submitting task %d for review by member %d: %w", task.ID(), member.ID(), err)
	}
	return nil
}

func SubmitTaskForReview(community domain.Community) core.Action {
	return submitTaskForReviewAction{
		community: community,
	}
}
