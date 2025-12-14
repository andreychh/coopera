package domain

import (
	"context"
)

type idempotencyCommunity struct {
	origin Community
}

func (i idempotencyCommunity) CreateUser(ctx context.Context, tgID int64, tgUsername string) (User, error) {
	user, exists, err := i.origin.UserWithTelegramID(ctx, tgID)
	if err != nil {
		return nil, err
	}
	if exists {
		return user, nil
	}
	return i.origin.CreateUser(ctx, tgID, tgUsername)
}

func (i idempotencyCommunity) UserWithID(ctx context.Context, id int64) (User, bool, error) {
	return i.origin.UserWithID(ctx, id)
}

func (i idempotencyCommunity) UserWithTelegramID(ctx context.Context, tgID int64) (User, bool, error) {
	return i.origin.UserWithTelegramID(ctx, tgID)
}

func (i idempotencyCommunity) UserWithUsername(ctx context.Context, tgUsername string) (User, bool, error) {
	return i.origin.UserWithUsername(ctx, tgUsername)
}

func (i idempotencyCommunity) Team(ctx context.Context, id int64) (Team, bool, error) {
	return i.origin.Team(ctx, id)
}

func (i idempotencyCommunity) Task(ctx context.Context, id int64) (Task, bool, error) {
	return i.origin.Task(ctx, id)
}

func IdempotencyCommunity(origin Community) Community {
	return idempotencyCommunity{
		origin: origin,
	}
}
