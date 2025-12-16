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

type assignTaskToSelfAction struct {
	community domain.Community
}

func (a assignTaskToSelfAction) Perform(ctx context.Context, update telegram.Update) error {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found in update")
	}
	user, exists, err := a.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	if !exists {
		return fmt.Errorf("user with telegram ID %d does not exist", chatID)
	}
	callbackData, found := attributes.CallbackData().Value(update)
	if !found {
		return fmt.Errorf("callback data not found in update")
	}
	taskID, err := protocol.ParseTaskID(callbackData)
	if err != nil {
		return fmt.Errorf("parsing task ID from callback data %q: %w", callbackData, err)
	}
	task, exists, err := a.community.Task(ctx, taskID)
	if err != nil {
		return fmt.Errorf("getting task %d: %w", taskID, err)
	}
	if !exists {
		return fmt.Errorf("task %d does not exist", taskID)
	}
	team, err := task.Team(ctx)
	if err != nil {
		return fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return fmt.Errorf("getting members of team %d: %w", team.ID(), err)
	}
	member, exists, err := members.MemberWithUsername(ctx, user.Username())
	if err != nil {
		return fmt.Errorf("getting member with username %q in team %d: %w", user.Username(), team.ID(), err)
	}
	if !exists {
		return fmt.Errorf("member with username %q does not exist in team %d", user.Username(), team.ID())
	}
	err = member.AssignTask(ctx, task.ID(), member.ID())
	if err != nil {
		return fmt.Errorf("assigning task %d to member %d: %w", task.ID(), member.ID(), err)
	}
	return nil
}

func AssignTaskToSelf(community domain.Community) core.Action {
	return assignTaskToSelfAction{
		community: community,
	}
}
