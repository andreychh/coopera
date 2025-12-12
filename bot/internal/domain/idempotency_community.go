package domain

import (
	"context"
	"errors"

	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type idempotencyCommunity struct {
	origin Community
}

func (i idempotencyCommunity) CreateUser(ctx context.Context, tgID int64, tgUsername string) (User, error) {
	_, err := i.origin.CreateUser(ctx, tgID, tgUsername)
	if errors.Is(err, transport.ErrRecordAlreadyExists) {
		user, err := i.origin.UserWithTelegramID(ctx, tgID)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, err
}

func (i idempotencyCommunity) UserWithTelegramID(ctx context.Context, tgID int64) (User, error) {
	return i.origin.UserWithTelegramID(ctx, tgID)
}

func (i idempotencyCommunity) UserWithUsername(ctx context.Context, tgUsername string) (User, error) {
	return i.origin.UserWithUsername(ctx, tgUsername)
}

func (i idempotencyCommunity) UserWithTelegramIDExists(ctx context.Context, tgID int64) (bool, error) {
	return i.origin.UserWithTelegramIDExists(ctx, tgID)
}

func (i idempotencyCommunity) UserWithUsernameExists(ctx context.Context, tgUsername string) (bool, error) {
	return i.origin.UserWithUsernameExists(ctx, tgUsername)
}

func (i idempotencyCommunity) Team(ctx context.Context, id int64) (Team, error) {
	return i.origin.Team(ctx, id)
}

func IdempotencyCommunity(origin Community) Community {
	return idempotencyCommunity{
		origin: origin,
	}
}
